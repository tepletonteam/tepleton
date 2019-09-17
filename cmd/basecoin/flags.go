package main

import (
	"github.com/urfave/cli"
)

// start flags
var (
	addrFlag = cli.StringFlag{
		Name:  "address",
		Value: "tcp://0.0.0.0:46658",
		Usage: "Listen address",
	}

	eyesFlag = cli.StringFlag{
		Name:  "eyes",
		Value: "local",
		Usage: "MerkleEyes address, or 'local' for embedded",
	}

	// TODO: move to config file
	// eyesCacheSizePtr := flag.Int("eyes-cache-size", 10000, "MerkleEyes db cache size, for embedded")

	dirFlag = cli.StringFlag{
		Name:  "dir",
		Value: ".",
		Usage: "Root directory",
	}

	inProcTMFlag = cli.BoolFlag{
		Name:  "in-proc",
		Usage: "Run Tendermint in-process with the App",
	}

	abiPluginFlag = cli.BoolFlag{
		Name:  "abi-plugin",
		Usage: "Enable the abi plugin",
	}

	counterPluginFlag = cli.BoolFlag{
		Name:  "counter-plugin",
		Usage: "Enable the counter plugin",
	}
)

// tx flags

var (
	nodeFlag = cli.StringFlag{
		Name:  "node",
		Value: "tcp://localhost:46657",
		Usage: "Tendermint RPC address",
	}

	toFlag = cli.StringFlag{
		Name:  "to",
		Value: "",
		Usage: "Destination address for the transaction",
	}

	amountFlag = cli.IntFlag{
		Name:  "amount",
		Value: 0,
		Usage: "Amount of coins to send in the transaction",
	}

	fromFlag = cli.StringFlag{
		Name:  "from",
		Value: "priv_validator.json",
		Usage: "Path to a private key to sign the transaction",
	}

	seqFlag = cli.IntFlag{
		Name:  "sequence",
		Value: 0,
		Usage: "Sequence number for the account",
	}

	coinFlag = cli.StringFlag{
		Name:  "coin",
		Value: "blank",
		Usage: "Specify a coin denomination",
	}

	gasFlag = cli.IntFlag{
		Name:  "gas",
		Value: 0,
		Usage: "The amount of gas for the transaction",
	}

	feeFlag = cli.IntFlag{
		Name:  "fee",
		Value: 0,
		Usage: "The transaction fee",
	}

	dataFlag = cli.StringFlag{
		Name:  "data",
		Value: "",
		Usage: "Data to send with the transaction",
	}

	nameFlag = cli.StringFlag{
		Name:  "name",
		Value: "",
		Usage: "Plugin to send the transaction to",
	}

	chainIDFlag = cli.StringFlag{
		Name:  "chain_id",
		Value: "test_chain_id",
		Usage: "ID of the chain for replay protection",
	}

	validFlag = cli.BoolFlag{
		Name:  "valid",
		Usage: "Set valid field in CounterTx",
	}
)

// abi flags
var (
	abiChainIDFlag = cli.StringFlag{
		Name:  "chain_id",
		Usage: "ChainID for the new blockchain",
		Value: "",
	}

	abiGenesisFlag = cli.StringFlag{
		Name:  "genesis",
		Usage: "Genesis file for the new blockchain",
		Value: "",
	}

	abiHeaderFlag = cli.StringFlag{
		Name:  "header",
		Usage: "Block header for an abi update",
		Value: "",
	}

	abiCommitFlag = cli.StringFlag{
		Name:  "commit",
		Usage: "Block commit for an abi update",
		Value: "",
	}

	abiFromFlag = cli.StringFlag{
		Name:  "from",
		Usage: "Source ChainID",
		Value: "",
	}

	abiToFlag = cli.StringFlag{
		Name:  "to",
		Usage: "Destination ChainID",
		Value: "",
	}

	abiTypeFlag = cli.StringFlag{
		Name:  "type",
		Usage: "ABI packet type (eg. coin)",
		Value: "",
	}

	abiPayloadFlag = cli.StringFlag{
		Name:  "payload",
		Usage: "ABI packet payload",
		Value: "",
	}

	abiPacketFlag = cli.StringFlag{
		Name:  "packet",
		Usage: "hex-encoded ABI packet",
		Value: "",
	}

	abiProofFlag = cli.StringFlag{
		Name:  "proof",
		Usage: "hex-encoded proof of ABI packet from source chain",
		Value: "",
	}

	abiSequenceFlag = cli.IntFlag{
		Name:  "sequence",
		Usage: "sequence number for ABI packet",
		Value: 0,
	}

	abiHeightFlag = cli.IntFlag{
		Name:  "height",
		Usage: "Height the packet became egress in source chain",
		Value: 0,
	}
)

// proof flags
var (
	proofFlag = cli.StringFlag{
		Name:  "proof",
		Usage: "hex-encoded IAVL proof",
		Value: "",
	}

	keyFlag = cli.StringFlag{
		Name:  "key",
		Usage: "key to the IAVL tree",
		Value: "",
	}

	valueFlag = cli.StringFlag{
		Name:  "value",
		Usage: "value in the IAVL tree",
		Value: "",
	}

	rootFlag = cli.StringFlag{
		Name:  "root",
		Usage: "root hash of the IAVL tree",
		Value: "",
	}
)
