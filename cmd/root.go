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

package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"perun.network/go-perun/log"

	"github.com/perun-network/perun-polkadot-demo/cmd/demo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:              "perun-polkadot-demo",
	Short:            "Perun State Channels Demo",
	Long:             "Demonstrator for the Perun state channel framework using the Polkadot backend.",
	PersistentPreRun: runRoot,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logConfig.Level, "log-level", "warn", "Logrus level")
	err := viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
	if err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().StringVar(&logConfig.File, "log-file", "", "log file")
	err = viper.BindPFlag("log.file", rootCmd.PersistentFlags().Lookup("log-file"))
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(demo.GetDemoCmd())
}

func runRoot(c *cobra.Command, args []string) {
	setConfig()
}

// Execute called by rootCmd
func Execute() {
	defer func() {
		if err := recover(); err != nil {
			log.Panicf("err=%s, trace=%s\n", err, debug.Stack())
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
