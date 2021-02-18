package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher"

	"github.com/skycoin/cx-chains/src/api"

	"github.com/skycoin/cx/cxgo/parser"

	"github.com/skycoin/cx-chains/src/cx/cxutil"

	"github.com/skycoin/cx-chains/src/cx/cxflags"
	"github.com/skycoin/cx-chains/src/cx/cxspec"
)

type runFlags struct {
	cmd *flag.FlagSet

	debugLexer   bool
	debugProfile int
	*cxflags.MemoryFlags

	inject   bool   // Whether to inject transaction to cx chain.
	nodeAddr string // CX Chain node address.
}

// printRunENVs prints ENVs in the 'run' help menu
func printRunENVs(cmd *flag.FlagSet) {
	cxutil.CmdPrintf(cmd, "ENVs:")
	cxutil.CmdPrintf(cmd, "  $%s\n  \t%s", genSKEnv, "genesis secret key (hex), required if '-inject' flag is set")
}

// readRunENVs parses ENVs specified for the 'run' subcommand
func readRunENVs(specAddr cipher.Address) cipher.SecKey {
	genSK, err := parseGenesisSKEnv()
	if err != nil {
		log.WithError(err).
			WithField(genSKEnv, genSK.Hex()).
			Fatal("Failed to read secret key from ENV.")
	}

	genAddr, err := cipher.AddressFromSecKey(genSK)
	if err != nil {
		log.WithError(err).
			WithField(genSKEnv, genSK.Hex()).
			Fatal("Failed to extract genesis address.")
	}

	if genAddr != specAddr {
		log.WithField(genSKEnv, genSK.Hex()).
			Fatal("Provided genesis secret key does not match genesis address from chain spec.")
	}

	return genSK
}

func processRunFlags(args []string) (runFlags, cxspec.ChainSpec, cipher.SecKey) {
	cmd := flag.NewFlagSet("cxchain-cli run", flag.ExitOnError)
	spec := processSpecFlags(context.Background(), cmd, args)

	f := runFlags{
		cmd: cmd,

		debugLexer:   false,
		debugProfile: 0,
		MemoryFlags:  cxflags.DefaultMemoryFlags(),

		inject: false,

		// TODO @evanlinjin: Find a way to set this later on.
		// TODO @evanlinjin: This way, we would not need to call '.Locate' so
		// TODO @evanlinjin: early within processSpecFlags()
		nodeAddr: fmt.Sprintf("http://127.0.0.1:%d", spec.Node.WebInterfacePort),
	}

	f.cmd.Usage = func() {
		usage := cxutil.DefaultUsageFormat("flags", "cx source files")
		usage(f.cmd, nil)
		printRunENVs(f.cmd)
	}

	f.cmd.BoolVar(&f.debugLexer, "debug-lexer", f.debugLexer, "enable lexer debugging by printing all scanner tokens")
	f.cmd.IntVar(&f.debugProfile, "debug-profile", f.debugProfile, "enable CPU+MEM profiling and set CPU profiling rate. Visualize .pprof files with 'go get github.com/google/pprof' and 'pprof -http=:8080 file.pprof'")
	f.MemoryFlags.Register(f.cmd)

	f.cmd.BoolVar(&f.inject, "inject", f.inject, "whether to inject this as a transaction on the cx chain")
	f.cmd.BoolVar(&f.inject, "i", f.inject, "shorthand for 'inject'")

	f.cmd.StringVar(&f.nodeAddr, "node", f.nodeAddr, "HTTP API `ADDRESS` of cxchain node")
	f.cmd.StringVar(&f.nodeAddr, "n", f.nodeAddr, "shorthand for 'node'")

	// Parse flags.
	if err := f.cmd.Parse(args); err != nil {
		os.Exit(1)
	}

	// Ensure genesis secret key is provided if 'inject' flag is set.
	var genSK cipher.SecKey
	if f.inject {
		genSK = readRunENVs(cipher.MustDecodeBase58Address(spec.GenesisAddr))
	}

	// Log stuff.
	cxflags.LogMemFlags(log)

	// Return.
	return f, spec, genSK
}

func cmdRun(args []string) {
	flags, spec, genSK := processRunFlags(args)

	// Apply debug flags.
	parser.DebugLexer = flags.debugLexer

	// Parse for cx args for genesis program state.
	log.Info("Parsing for CX args...")
	cxRes, err := cxutil.ExtractCXArgs(flags.cmd, true)
	if err != nil {
		log.WithError(err).Fatal("Failed to extract CX args.")
	}
	cxFilenames := cxutil.ListSourceNames(cxRes.CXSources, true)
	log.WithField("filenames", cxFilenames).Info("Obtained CX sources.")

	// Prepare API Client.
	c := api.NewClient(flags.nodeAddr)

	// Prepare address.
	addr := cipher.MustDecodeBase58Address(spec.GenesisAddr)

	// Parse and run program.
	ux, progB, err := PrepareChainProg(cxFilenames, cxRes.CXSources, c, addr, flags.debugLexer, flags.debugProfile)
	if err != nil {
		log.WithError(err).Fatal("Failed to prepare chain CX program.")
	}

	if flags.inject {
		// Run: inject.
		if err := BroadcastMainExp(c, genSK, ux); err != nil {
			log.WithError(err).Fatal("Failed to broadcast transaction.")
		}
	} else {
		// Run: without injection.
		if err := RunChainProg(cxRes.CXFlags, progB); err != nil {
			log.WithError(err).Fatal("Failed to run chain CX program.")
		}
	}
}
