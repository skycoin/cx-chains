package cxspec

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/logging"
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
	// defaultSpecFilepath is the default cx spec filepath.
	// This is for internal use.
	defaultSpecFilepath = "skycoin.chain_spec.json"

	// DefaultSpecLocStr is the default cx spec location string.
	DefaultSpecLocStr = string(FileLoc + ":" + defaultSpecFilepath)

	// DefaultTrackerURL is the default cx tracker URL.
	DefaultTrackerURL = "https://cxt.skycoin.com"
)

// Possible errors when executing 'Locate'.
var (
	ErrEmptySpec        = errors.New("empty chain spec provided")
	ErrEmptyTracker     = errors.New("tracker is not provided")
	ErrInvalidLocPrefix = errors.New("invalid spec location prefix")
)

// Locate locates the chain spec given a 'loc' string.
// The 'loc' string is to be of format '<location-prefix>:<location>'.
// * <location-prefix> is 'tracker' if undefined.
// * <location> either specifies the cx chain's genesis hash (if
// <location-prefix> is 'tracker') or filepath of the spec file (if
// <location-prefix> is 'file').
func Locate(ctx context.Context, log logrus.FieldLogger, tracker *CXTrackerClient, loc string) (ChainSpec, error) {
	// Ensure logger is existent.
	if log == nil {
		log = logging.MustGetLogger("cxspec").WithField("func", "Locate")
	}

	prefix, suffix, err := splitLocString(loc)
	if err != nil {
		return ChainSpec{}, err
	}

	// Check location prefix (LocPrefix).
	switch prefix {
	case FileLoc:
		if suffix == "" {
			suffix = defaultSpecFilepath
		}

		return ReadSpecFile(suffix)

	case TrackerLoc:
		// Check that 'tracker' is not nil.
		if tracker == nil {
			return ChainSpec{}, ErrEmptyTracker
		}

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

// LocateConfig contains flag values for Locate.
type LocateConfig struct {
	CXChain   string // CX Chain spec location string.
	CXTracker string // CX Tracker URL.
}

// FillDefaults fills LocateConfig with default values.
func (c *LocateConfig) FillDefaults() {
	c.CXChain = DefaultSpecLocStr
	c.CXTracker = DefaultTrackerURL
}

// DefaultLocateConfig returns the default LocateConfig set.
func DefaultLocateConfig() LocateConfig {
	var lc LocateConfig
	lc.FillDefaults()
	return lc
}

// Parse parses the OS args for the 'chain' flag.
func (c *LocateConfig) Parse(args []string) {
	c.CXChain = obtainFlagValue(args, "chain")
	c.CXTracker = obtainFlagValue(args, "tracker")
}

// RegisterFlags ensures that the 'help' menu contains the locate flags and that
// the flags are recognized.
func (c *LocateConfig) RegisterFlags(fs *flag.FlagSet) {
	var temp string
	fs.StringVar(&temp, "chain", c.CXChain, fmt.Sprintf("cx chain location. Prepend with '%s:' or '%s:' for spec location type.", FileLoc, TrackerLoc))
	fs.StringVar(&temp, "tracker", c.CXTracker, "CX Tracker `URL`.")
}

func obtainFlagValue(args []string, key string) string {
	var (
		keyPrefix1 = "-" + key
		keyPrefix2 = keyPrefix1 + "="
	)

	for i, a := range args {
		// Standardize flag prefix to single '-'.
		if strings.HasPrefix(a, "--") {
			a = a[1:]
		}

		// If there is no '=', the flag value is the next arg.
		if a == "-"+key && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(a, keyPrefix2) {
			return strings.TrimPrefix(a, keyPrefix2)
		}
	}

	return ""
}
