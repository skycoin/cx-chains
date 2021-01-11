package cxspec

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

// LocPrefix determines the location type of the location string.
type LocPrefix string

// Locations types.
const (
	FileLoc    = LocPrefix("file")
	TrackerLoc = LocPrefix("tracker")
)

// Constants.
const (
	defaultChainSpecFile = "./skycoin.chain_spec.json"
)

// Possible errors when executing 'Locate'.
var (
	ErrEmptySpec = errors.New("empty chain spec provided")
	ErrInvalidLocPrefix = errors.New("invalid spec location prefix")
)

// Locate locates the chain spec given a 'loc' string.
// The 'loc' string is to be of format '<location-prefix>:<location>'.
// * <location-prefix> is 'tracker' if undefined.
// * <location> either specifies the cx chain's genesis hash (if
// <location-prefix> is 'tracker') or filepath of the spec file (if
// <location-prefix> is 'file').
func Locate(ctx context.Context, tracker *CXTrackerClient, loc string) (ChainSpec, error) {
	prefix, suffix, err := splitLocString(loc)
	if err != nil {
		return ChainSpec{}, err
	}

	// Check location prefix (LocPrefix).
	switch prefix {
	case FileLoc:
		if suffix == "" {
			suffix = defaultChainSpecFile
		}

		return ReadSpecFile(suffix)

	case TrackerLoc:
		// Obtain genesis hash from hex string.
		hash, err := cipher.SHA256FromHex(suffix)
		if err != nil {
			return ChainSpec{}, fmt.Errorf("invalid genesis hash provided '%s': %w", loc, err)
		}

		// Obtain spec from tracker.
		signedChainSpec, err := tracker.SpecByGenesisHash(ctx, hash)
		if err != nil {
			return ChainSpec{}, fmt.Errorf("chain spec not of genesis hash not found in tracker: %w", err)
		}

		// Verify again (no harm in doing it twice).
		if err := signedChainSpec.Verify(); err != nil {
			return ChainSpec{}, err
		}

		return signedChainSpec.Spec, nil

	default:
		return ChainSpec{}, fmt.Errorf("%w '%s'", ErrInvalidLocPrefix, prefix)
	}
}

func splitLocString(loc string) (prefix LocPrefix, suffix string, err error) {
	loc = strings.TrimSpace(loc)
	if loc == "" {
		return "", "", ErrEmptySpec
	}

	locParts := strings.SplitN(loc, ":", 2)

	switch len(locParts) {
	case 1:
		locParts = append([]string{string(TrackerLoc)}, locParts...)
	case 2:
		// continue
	default:
		panic("internal error: Locate() should never return >2 location parts")
	}

	return LocPrefix(locParts[0]), locParts[1], nil
}
