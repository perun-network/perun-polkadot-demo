// Copyright 2021 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package demo

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/pkg/errors"

	sr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"perun.network/go-perun/channel/persistence/keyvalue"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
	"perun.network/go-perun/pkg/sortedkv/leveldb"
	wirenet "perun.network/go-perun/wire/net"
	"perun.network/go-perun/wire/net/simple"
)

var backend *node

// Setup initializes the node, can not be done in init() since it needs the
// configuration from viper.
func Setup() {
	SetConfig(flags.cfgFile, flags.cfgNetFile)

	var err error
	if backend, err = newNode(); err != nil {
		log.WithError(err).Fatalln("Could not initialize node.")
	}
}

func newNode() (*node, error) {
	wallet, acc, err := setupWallet(config.Sk)
	if err != nil {
		return nil, errors.WithMessage(err, "importing mnemonic")
	}
	dot, err := newDotSetup(acc, config.Chain)
	if err != nil {
		return nil, errors.WithMessage(err, "creating dot setup")
	}
	dialer := simple.NewTCPDialer(config.Node.DialTimeout)

	n := &node{
		log:         log.Get(),
		onChain:     acc,
		wallet:      wallet,
		api:         dot.Api,
		adjudicator: dot.Adjudicator,
		funder:      dot.Funder,
		dialer:      dialer,
		peers:       make(map[string]*peer),
	}
	return n, n.setup()
}

func (n *node) setup() error {
	var err error

	sk, err := sr25519.NewSKFromRng(rand.Reader)
	if err != nil {
		return errors.WithMessage(err, "generating off-chain account")
	}
	n.offChain = n.wallet.ImportSK(sk)
	n.log.WithField("off-chain", n.offChain.Address()).Info("Generated account")
	n.bus = wirenet.NewBus(n.onChain, n.dialer)

	if n.client, err = client.New(n.onChain.Address(), n.bus, n.funder, n.adjudicator, n.wallet); err != nil {
		return errors.WithMessage(err, "creating client")
	}

	host := config.Node.IP + ":" + strconv.Itoa(int(config.Node.Port))
	n.log.WithField("host", host).Trace("Listening for connections")
	listener, err := simple.NewTCPListener(host)
	if err != nil {
		return errors.WithMessage(err, "could not start tcp listener")
	}

	n.client.OnNewChannel(n.setupChannel)
	if err := n.setupPersistence(); err != nil {
		return errors.WithMessage(err, "setting up persistence")
	}
	go n.client.Handle(n, n)
	go n.bus.Listen(listener)
	return n.PrintConfig()
}

func (n *node) setupPersistence() error {
	if config.Node.PersistenceEnabled {
		n.log.Info("Starting persistence")
		db, err := leveldb.LoadDatabase(config.Node.PersistencePath)
		if err != nil {
			return errors.WithMessage(err, "creating/loading database")
		}
		persister := keyvalue.NewPersistRestorer(db)
		n.client.EnablePersistence(persister)

		ctx, cancel := context.WithTimeout(context.Background(), config.Node.ReconnectTimeout)
		defer cancel()
		if err := n.client.Restore(ctx); err != nil {
			n.log.WithError(err).Warn("Could not restore client")
			// return the error.
		}
	} else {
		n.log.Info("Persistence disabled")
	}
	return nil
}

func setupWallet(hexSk string) (*dotwallet.Wallet, *dotwallet.Account, error) {
	wallet := dotwallet.NewWallet()
	sk, err := sr25519.NewSKFromHex(hexSk)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "creating hdwallet")
	}
	return wallet, wallet.ImportSK(sk), nil
}

func (n *node) PrintConfig() error {
	fmt.Printf(
		"Alias: %s\n"+
			"Listening: %s:%d\n"+
			"Node RPC URL: %s\n"+
			"Perun ID: %s\n"+
			"OffChain: %s\n"+
			"", config.Alias, config.Node.IP, config.Node.Port, config.Chain.NodeUrl, n.onChain.Address().String(), n.offChain.Address().String())

	fmt.Println("Known peers:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	for alias, peer := range config.Peers {
		fmt.Fprintf(w, "%s\t%v\t%s:%d\n", alias, peer.PerunID, peer.Hostname, peer.Port)
	}
	return w.Flush()
}
