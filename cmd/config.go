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

package cmd // import "github.com/perun-network/perun-polkadot-demo/cmd"

import (
	"os"

	"perun.network/go-perun/log"
	plogrus "perun.network/go-perun/log/logrus"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type logCfg struct {
	Level string
	level logrus.Level
	File  string
}

var logConfig logCfg

func setConfig() {
	lvl, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "parsing log level"))
	}
	logConfig.level = lvl

	// Set the logging output file
	logger := logrus.New()
	if logConfig.File != "" {
		f, err := os.OpenFile(logConfig.File,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(errors.WithMessage(err, "opening logging file"))
		}
		logger.SetOutput(f)
	}
	logger.SetLevel(lvl)
	log.Set(plogrus.FromLogrus(logger))
}
