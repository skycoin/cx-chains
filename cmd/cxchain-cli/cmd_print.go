package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/skycoin/cx-chains/src/cx/cxspec"
	"github.com/skycoin/cx-chains/src/cx/cxutil"
	"os"
)

type printFlags struct {
	cmd *flag.FlagSet

	in    string
	field string
}

func processPrintFlags(args []string) printFlags {
	// Specify default flag values.
	f := printFlags{
		cmd:      flag.NewFlagSet("cxchain-cli print", flag.ExitOnError),
		in: cxspec.DefaultSpecFilepath,
		field:    "seckey",
	}

	f.cmd.Usage = func() {
		usage := cxutil.DefaultUsageFormat("flags")
		usage(f.cmd, nil)
	}

	f.cmd.StringVar(&f.in, "in", f.in, "`FILENAME` of file to read in")
	f.cmd.StringVar(&f.field, "field", f.field, "`NAME` of field to print")

	if err := f.cmd.Parse(args); err != nil {
		os.Exit(1)
	}

	return f
}

func cmdPrint(args []string) {
	flags := processPrintFlags(args)

	f, err := os.Open(flags.in)
	if err != nil {
		log.WithError(err).
			Fatal("Failed to read in file.")
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.WithError(err).
				Fatal("Failed to close file.")
		}
	}()

	var kSpec cxspec.KeySpec
	if err := json.NewDecoder(f).Decode(&kSpec); err != nil {
		log.WithError(err).
			Fatal("Failed to decode file.")
	}

	var out string

	switch flags.field {
	case "spec_era":
		out = kSpec.SpecEra
	case "key_type":
		out = kSpec.KeyType
	case "pubkey":
		out = kSpec.PubKey
	case "seckey":
		out = kSpec.SecKey
	case "address":
		out = kSpec.Address
	default:
		log.WithField("field", flags.field).
			Fatal("Invalid field input.")
	}

	fmt.Println(out)
}
