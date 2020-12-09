# First Run

> This document should be integrated within `README.md` in the near future.

Start `cx-tracker` with default setting.
* Compile from [github.com/skycoin/cx-tracker](https://github.com/skycoin/cx-tracker).
```bash
$ cx-tracker
```

Generate new chain spec.
```bash
$ cxchain-cli new ./examples/blockchain/counter-bc.cx
```

Run publisher node with generated chain spec.
* Obtain the chain secret key from generated `{coin}.chain_keys.json` file.
```bash
$ CXCHAIN_SK={publisher_secret_key} cxchain -enable-all-api-sets
```

Run client node with generated chain spec (use different data dir, and ports to publisher node).
* As no `CXCHAIN_SK` is provided, a random key pair is generated for the node.
```bash
$ cxchain -enable-all-api-sets -data-dir "$HOME/.cxchain/skycoin_client" -port 6002 -web-interface-port 6422
```

Run transaction against publisher node.
```bash
$ cxchain-cli run ./examples/blockchain/counter-txn.cx
```

Run transaction against client node and inject.
```bash
$ CXCHAIN_GEN_SK={genesis_secret_key} cxchain-cli run -n "http://127.0.0.1:6422" -i ./examples/blockchain/counter-txn.cx
```
