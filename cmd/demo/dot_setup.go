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
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"

	"github.com/perun-network/perun-polkadot-backend/channel/pallet"
	dot "github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	pchannel "perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"
)

type (
	chainConfig struct {
		NodeUrl         string        `json:"node_url"`
		NetworkId       dot.NetworkID `json:"network_id"`
		BlockTimeSec    uint32        `json:"block_time_sec"`
		TxTimeoutSec    uint32        `json:"tx_timeout_sec"`
		BlockQueryDepth uint32        `json:"block_query_depth"` // Actually of type types.BlockNumber.
	}

	dotSetup struct {
		Api         *dot.API
		Funder      pchannel.Funder
		Adjudicator pchannel.Adjudicator
	}
)

func newDotSetup(acc pwallet.Account, cfg chainConfig) (*dotSetup, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	api, err := dot.NewAPI(cfg.NodeUrl, cfg.NetworkId)
	if err != nil {
		return nil, err
	}
	perun := pallet.NewPallet(pallet.NewPerunPallet(api), api.Metadata())
	funder := pallet.NewFunder(perun, acc, 3)
	adj := pallet.NewAdjudicator(acc, perun, api, types.BlockNumber(cfg.BlockQueryDepth))
	return &dotSetup{api, funder, adj}, nil
}

// validate checks the config for some obvious errors.
func (c *chainConfig) validate() error {
	switch {
	case c.NodeUrl == "":
		return errors.New("empty node url")
	case c.BlockTimeSec == 0 || c.BlockTimeSec > 60:
		return errors.New("block time out of range")
	case c.BlockQueryDepth < 1 || c.BlockQueryDepth > 1000:
		return errors.New("block query depth out of range")
	case c.TxTimeoutSec < c.BlockTimeSec || c.TxTimeoutSec > 60*c.BlockTimeSec:
		return errors.New("tx timeout out of range")
	default:
		return nil
	}
}
