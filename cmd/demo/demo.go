// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of
// perun-polkadot-demo. Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package demo

import (
	prompt "github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Two party payment Demo",
	Long: `Enables two user to send payments between each other in a ledger state channel.
	The channels are funded and settled on an Polkadot blockchain, leaving out the dispute case.

	It illustrates what Perun is capable of.`,
	Run: runDemo,
}

// CommandLineFlags contains the command line flags.
type CommandLineFlags struct {
	cfgFile    string
	cfgNetFile string
}

var flags CommandLineFlags

func init() {
	demoCmd.PersistentFlags().StringVar(&flags.cfgFile, "config", "config.yaml", "General config file")
	demoCmd.PersistentFlags().StringVar(&flags.cfgNetFile, "network", "network.yaml", "Network config file")
	demoCmd.PersistentFlags().BoolVar(&GetConfig().Node.PersistenceEnabled, "persistence", false, "Enables the persistence")
	demoCmd.PersistentFlags().StringVar(&GetConfig().Sk, "secretkey", "", "Hex secret key 0x…")
	viper.BindPFlag("secretkey", demoCmd.PersistentFlags().Lookup("secretkey"))
}

// GetDemoCmd exposes demoCmd so that it can be used as a sub-command by another cobra command instance.
func GetDemoCmd() *cobra.Command {
	return demoCmd
}

// runDemo is executed everytime the program is started with the `demo` sub-command.
func runDemo(c *cobra.Command, args []string) {
	Setup()
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("perun"),
	)
	p.Run()
}

func completer(prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{}
}

// executor wraps the demo executor to print error messages.
func executor(in string) {
	AddInput(in)
}
