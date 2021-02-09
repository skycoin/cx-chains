# Example Environment

Create temporary working directory.

```bash
$ mkdir temp
```

Run `cx-tracker` in a terminal.

```bash
$ cx-tracker --db="temp/cx-tracker.db" --addr=":9091"
```

Generate chain spec file.

```bash
$ cxchain-cli new \
    --coin="tempcoin" \
    --chain-keys-output="temp/tempcoin.chain_keys.json" \
    --chain-spec-output="temp/tempcoin.chain_spec.json" \
    --genesis-keys-output="temp/tempcoin.genesis_keys.json" \
    ./cx/examples/counter-bc.cx
```

Run publisher node.

> TODO @evanlinjin: Make cxchain-cli target to extract keys.