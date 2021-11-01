// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of
// perun-polkadot-demo. Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package demo

import (
	"math/big"
	"strconv"

	sr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"github.com/pkg/errors"
)

func valBal(input string) error {
	_, _, err := big.ParseFloat(input, 10, 64, big.ToNearestEven)
	return errors.Wrap(err, "parsing float")
}

func valUInt(input string) error {
	if n, err := strconv.Atoi(input); err != nil {
		return errors.New("Invalid integer")
	} else if n < 0 {
		return errors.New("Value must be > 0")
	}
	return nil
}

func valPeer(arg string) error {
	if !backend.ExistsPeer(arg) {
		return errors.Errorf("Unknown peer, use 'info' to see connected")
	}
	return nil
}

func valAlias(arg string) error {
	for alias := range config.Peers {
		if alias == arg {
			return nil
		}
	}
	return errors.Errorf("Unknown alias, use 'config' to see available")
}

// strToAddress parses a string as dotwallet.Address
func strToAddress(str string) (*dotwallet.Address, error) {
	pk, err := sr25519.NewPKFromHex(str)
	return dotwallet.NewAddressFromPK(pk), err
}

func dotToPlank(dots ...*big.Float) []*big.Int {
	planks := make([]*big.Int, len(dots))
	for idx, d := range dots {
		plankFloat := new(big.Float).Mul(d, new(big.Float).SetFloat64(dot.PlankPerDot))
		// accuracy (second return value) returns "exact" for specified input range, hence ignored.
		planks[idx], _ = plankFloat.Int(nil)
	}
	return planks
}
