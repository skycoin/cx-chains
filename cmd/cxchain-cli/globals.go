package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/skycoin/skycoin/src/cipher"

	"github.com/skycoin/cx-chains/src/cx/cxspec"
)

const (
	// genSKEnv is the ENV for genesis secret key.
	genSKEnv = "CXCHAIN_GEN_SK"

	// chainSKEnv is the ENV for chain secret key.
	chainSKEnv = "CXCHAIN_SK"
)

var ErrNoSKProvided = errors.New("no secret key provided")

// parseSKEnv parses secret key from CXCHAIN_SECRET_KEY env.
// The secret key can be null.
func parseSKEnv(envKey string) (cipher.SecKey, error) {
	if skStr, ok := os.LookupEnv(envKey); ok {
		sk, err := cipher.SecKeyFromHex(skStr)
		if err != nil {
			return cipher.SecKey{}, fmt.Errorf("failed to parse secret key defined in ENV '%s': %w", envKey, err)
		}
		return sk, nil
	}
	return cipher.SecKey{}, ErrNoSKProvided
}

// processSpecFlags parses the 'chain' flag which defines where to discover the
// chain spec file.
func processSpecFlags(ctx context.Context, cmd *flag.FlagSet, args []string) cxspec.ChainSpec {
	conf := cxspec.DefaultLocateConfig()
	conf.RegisterFlags(cmd)
	conf.SoftParse(args)

	spec, err := cxspec.LocateWithConfig(ctx, &conf)
	if err != nil {
		log.WithError(err).
			WithField("chain", conf.CXChain).
			Fatal("Failed to find cx spec.")
	}

	return spec
}
