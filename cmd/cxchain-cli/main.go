package main

import (
	"flag"
	"os"

	"github.com/skycoin/skycoin/src/util/logging"

	"github.com/skycoin/cx-chains/src/cx/cxutil"
)

var log = logging.MustGetLogger("cxchain-cli")

func main() {
	usageMenu := cxutil.UsageFormat(func(cmd *flag.FlagSet, subcommands []string) {
		// print: Usage
		cxutil.PrintCmdUsage(cmd, "Usage", subcommands, []string{"args"})

		// print: ENVs
		cxutil.CmdPrintf(cmd, "ENVs:")
		cxutil.CmdPrintf(cmd, "  $%s\n  \t%s", chainSKEnv, "chain secret key (hex)")
		cxutil.CmdPrintf(cmd, "  $%s\n  \t%s", genSKEnv, "genesis secret key (hex)")

		// print: Flags
		cxutil.PrintCmdFlags(cmd, "Flags")
	})

	root := cxutil.NewCommandMap(flag.CommandLine, 7, usageMenu).
		AddSubcommand("version", func([]string) { cmdVersion() }).
		AddSubcommand("help", func([]string) { flag.CommandLine.Usage() }).
		AddSubcommand("tokenize", cmdTokenize).
		AddSubcommand("new", cmdNew).
		AddSubcommand("post", cmdPost).
		AddSubcommand("run", cmdRun).
		AddSubcommand("state", cmdState).
		AddSubcommand("peers", cmdPeers).
		AddSubcommand("key", cmdKey).
		AddSubcommand("genesis", cmdGenesis)

	os.Exit(root.ParseAndRun(os.Args[1:]))
}
