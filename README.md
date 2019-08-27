# Tepleton
<img src="docs/tepleton_logo.png" width="250" height="250">

The Tepleton is a framework for building blockchain applications in Golang.


**Note**: Requires [Go 1.12+](https://golang.org/dl/)


# Basecoin


Basecoin is an [WRSP application](https://github.com/tepleton/wrsp) designed to be used with the [Tendermint consensus engine](https://tepleton.com/) to form a Proof-of-Stake cryptocurrency.
It also provides a general purpose framework for extending the feature-set of the cryptocurrency
by implementing plugins.

Basecoin serves as a reference implementation for how we build WRSP applications in Go,
and is the framework in which we implement the [Cosmos Hub](https://cosmos.network).
It's easy to use, and doesn't require any forking - just implement your plugin, import the basecoin libraries,
and away you go with a full-stack blockchain and command line tool for transacting.

WARNING: Currently uses plain-text private keys for transactions and is otherwise not production ready.

## Installation

On a good day, basecoin can be installed like a normal Go program:

```
go get -u github.com/tepleton/basecoin/cmd/basecoin
```

In some cases, if that fails, or if another branch is required,
we use `glide` for dependency management.

The guaranteed correct way of compiling from source, assuming you've already 
run `go get` or otherwise cloned the repo, is:

```
cd $GOPATH/src/github.com/tepleton/basecoin
git checkout develop # (until we release tepleton v0.9)
make get_vendor_deps
make install
```

This will create the `basecoin` binary in `$GOPATH/bin`.


## Command Line Interface

The basecoin CLI can be used to start a stand-alone basecoin instance (`basecoin start`),
or to start basecoin with Tendermint in the same process (`basecoin start --in-proc`).
It can also be used to send transactions, eg. `basecoin tx send --to 0x4793A333846E5104C46DD9AB9A00E31821B2F301 --amount 100btc,10gold`
See `basecoin --help` and `basecoin [cmd] --help` for more details`.

