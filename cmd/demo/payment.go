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
	"fmt"
	"math/big"

	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/pkg/errors"
	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	"perun.network/go-perun/log"
)

type (
	paymentChannel struct {
		*client.Channel

		log     log.Logger
		handler chan bool
	}
)

func newPaymentChannel(ch *client.Channel) *paymentChannel {
	return &paymentChannel{
		Channel: ch,
		log:     log.WithField("channel", ch.ID()),
		handler: make(chan bool, 1),
	}
}
func (ch *paymentChannel) sendMoney(amount *big.Int) error {
	return ch.sendUpdate(
		func(state *channel.State) error {
			transferBal(stateBals(state), ch.Idx(), amount)
			return nil
		}, "sendMoney")
}

func (ch *paymentChannel) sendFinal() error {
	ch.log.Debugf("Sending final state")
	return ch.sendUpdate(func(state *channel.State) error {
		state.IsFinal = true
		return nil
	}, "final")
}

func (ch *paymentChannel) sendUpdate(update func(*channel.State) error, desc string) error {
	ch.log.Debugf("Sending update: %s", desc)
	ctx, cancel := context.WithTimeout(context.Background(), config.Channel.Timeout)
	defer cancel()

	stateBefore := ch.State()
	err := ch.UpdateBy(ctx, update)
	ch.log.Debugf("Sent update: %s, err: %v", desc, err)

	state := ch.State()
	balChanged := stateBefore.Balances[0][0].Cmp(state.Balances[0][0]) != 0
	if balChanged {
		bals := dot.NewDotsFromPlanks(state.Allocation.Balances[0]...)
		fmt.Printf("💰 Sent payment. New balance: [My: %v, Peer: %v]\n", bals[ch.Idx()], bals[1-ch.Idx()]) // assumes two-party channel
	}

	return err
}

func transferBal(bals []channel.Bal, ourIdx channel.Index, amount *big.Int) {
	a := new(big.Int).Set(amount) // local copy because we mutate it
	otherIdx := ourIdx ^ 1
	ourBal := bals[ourIdx]
	otherBal := bals[otherIdx]
	otherBal.Add(otherBal, a)
	ourBal.Sub(ourBal, a)
}

func stateBals(state *channel.State) []channel.Bal {
	return state.Balances[0]
}

func (ch *paymentChannel) Handle(old *channel.State, update client.ChannelUpdate, res *client.UpdateResponder) {
	oldBal := stateBals(old)
	balChanged := oldBal[0].Cmp(update.State.Balances[0][0]) != 0
	ctx, cancel := context.WithTimeout(context.Background(), config.Channel.Timeout)
	defer cancel()
	if err := assertValidTransition(old, update.State, update.ActorIdx); err != nil {
		if err := res.Reject(ctx, "invalid transition"); err != nil {
			ch.log.WithError(err).Error("Could not reject channel proposal")
		} else {
			return
		}
	} else if err := res.Accept(ctx); err != nil {
		ch.log.Error(errors.WithMessage(err, "handling payment update"))
	}

	if balChanged {
		bals := dot.NewDotsFromPlanks(update.State.Allocation.Balances[0]...)
		PrintfAsync("💰 Received payment. New balance: [My: %v, Peer: %v]\n", bals[ch.Idx()], bals[1-ch.Idx()])
	}
}

// assertValidTransition checks that money flows only from the actor to the
// other participants.
func assertValidTransition(from, to *channel.State, actor channel.Index) error {
	if !channel.IsNoData(to.Data) {
		return errors.New("channel must not have app data")
	}
	for i, asset := range from.Balances {
		for j, bal := range asset {
			if int(actor) == j && bal.Cmp(to.Balances[i][j]) == -1 {
				return errors.Errorf("payer[%d] steals asset %d, so %d < %d", j, i, bal, to.Balances[i][j])
			} else if int(actor) != j && bal.Cmp(to.Balances[i][j]) == 1 {
				return errors.Errorf("payer[%d] reduces participant[%d]'s asset %d", actor, j, i)
			}
		}
	}
	return nil
}

func (ch *paymentChannel) GetBalances() (our, other *big.Int) {
	bals := stateBals(ch.State())
	if len(bals) != 2 {
		return new(big.Int), new(big.Int)
	}
	return bals[ch.Idx()], bals[1-ch.Idx()]
}
