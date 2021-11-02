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
	"log"
	"time"

	"github.com/spf13/viper"
	"perun.network/go-perun/wire"
)

// Config contains all configuration read from config.yaml and network.yaml
type (
	Config struct {
		Alias   string
		Sk      string
		Channel channelConfig
		Node    nodeConfig
		Chain   chainConfig
		// Read from the network.yaml. The key is the alias.
		Peers map[string]*netConfigEntry
	}

	channelConfig struct {
		Timeout              time.Duration
		FundTimeout          time.Duration
		SettleTimeout        time.Duration
		ChallengeDurationSec uint64
	}

	nodeConfig struct {
		IP               string
		Port             uint16
		DialTimeout      time.Duration
		HandleTimeout    time.Duration
		ReconnectTimeout time.Duration

		PersistencePath    string
		PersistenceEnabled bool
	}

	netConfigEntry struct {
		PerunID  string
		perunID  wire.Address
		Hostname string
		Port     uint16
	}
)

var config Config

// GetConfig returns a pointer to the current `Config`.
// This is needed to make viper and cobra work together.
func GetConfig() *Config {
	return &config
}

// SetConfig called by viper when the config file was parsed
func SetConfig(cfgPath, cfgNetPath string) {
	ParseConfig(cfgPath, cfgNetPath, &config)
}

func ParseConfig(cfgPath, cfgNetPath string, cfg *Config) {
	// Load config files
	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	viper.SetConfigFile(cfgNetPath)
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading network config file, %s", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatal(err)
	}

	for _, peer := range config.Peers {
		addr, err := strToAddress(peer.PerunID)
		if err != nil {
			log.Panic(err)
		}
		peer.perunID = addr
	}
}
