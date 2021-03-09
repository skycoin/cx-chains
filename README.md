# CX Chain

CX Chain is a fork of [Skycoin](https://github.com/SkycoinProject/skycoin) with the ability to run application-specific blockchains with smart-contracts written with [CX](https://github.com/skycoin/cx).

Each ***CX Chain*** is identified by a genesis hash, which in turn reference a ***CX Chain Spec***. The CX Chain Spec specifies the properties and genesis block of the specific CX Chain.

## Install

### Dependencies

CX Chain requires [Golang](https://golang.org/) to compile (version `1.14+`). Detailed installation instructions can be found here: https://github.com/SkycoinProject/skycoin/blob/develop/INSTALLATION.md

### Build

To build `cxchain`, the typical Golang binary build process applies. The following command builds `cxchain` and `cxchain-cli` into the target directory specified by the `GOBIN` env.

```bash
$ git clone git@github.com:skycoin/cx-chains.git && cd cx-chains
$ go install ./cmd/...
```

The `go install` command is also available as a `Makefile` target.

```bash
$ make install
```

## Run

### Dependencies

You will need to specify an address of a `cx-tracker` for a `cxchain` instance to function properly. A local `cx-tracker` instance can be installed via [this repository](https://github.com/skycoin/cx-tracker).

### Run a Local CX Chain Environment

*This local environment has two `cxchain` instances and a `cx-tracker`.*

Start `cx-tracker`.
```bash
$ cx-tracker -addr ":9091"
```

Generate new chain spec (assuming that the repository root is your working directory).
```bash
$ cxchain-cli new ./cx/examples/counter-bc.cx
```

Post chain spec to `cx-tracker`.
```bash
$ export CXCHAIN_SK=$(cxchain-cli key -in skycoin.chain_keys.json -field "seckey")
$ cxchain-cli post -t "http://127.0.0.1:9091" -s skycoin.chain_spec.json
```

At this point, you can head to [http://127.0.0.1:9091/api/specs](http://127.0.0.1:9091/api/specs) to see whether the spec is posted to `cx-tracker`.

Run publisher node with generated chain spec.
* Obtain the chain secret key from generated `{coin}.chain_keys.json` file.
```bash
$ export CXCHAIN_SK=$(cxchain-cli key -in skycoin.chain_keys.json -field "seckey")
$ export CXCHAIN_HASH=$(cxchain-cli genesis -in skycoin.chain_spec.json)
$ cxchain -chain "tracker:$CXCHAIN_HASH" -tracker "http://127.0.0.1:9091" -enable-all-api-sets -data-dir ./master_node -port 6001 -web-interface-port 6421
```

Run client node with generated chain spec (use different data dir, and ports to publisher node).
* As no `CXCHAIN_SK` is provided, a random key pair is generated for the node.
```bash
$ export CXCHAIN_HASH=$(cxchain-cli genesis -in skycoin.chain_spec.json)
$ cxchain -chain "tracker:$CXCHAIN_HASH" -tracker "http://127.0.0.1:9091" -client -enable-all-api-sets -data-dir ./client_node -port 6002 -web-interface-port 6422
```

Run transaction against publisher node.
```bash
$ export CXCHAIN_HASH=$(cxchain-cli genesis -in skycoin.chain_spec.json)
$ cxchain-cli run -chain "tracker:$CXCHAIN_HASH" -tracker "http://127.0.0.1:9091" ./cx/examples/counter-tx.cx
```

Run transaction against client node and inject.
```bash
$ export CXCHAIN_GEN_SK=$(cxchain-cli key -in skycoin.genesis_keys.json -field "seckey")
$ export CXCHAIN_HASH=$(cxchain-cli genesis -in skycoin.chain_spec.json)
$ cxchain-cli run -chain "tracker:$CXCHAIN_HASH" -tracker "http://127.0.0.1:9091" -node "http://127.0.0.1:6422" -inject ./cx/examples/counter-tx.cx
```

## Resources

- [CX Chain Technical Overview](./doc/CXCHAIN_OVERVIEW.md)
- [`skycoin/cx` Repository](https://github.com/skycoin/cx)
- [`skycoin/cx-tracker` Repository](https://github.com/skycoin/cx-tracker)