package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/urfave/cli"

	"github.com/tepleton/basecoin/plugins/abi"

	cmn "github.com/tepleton/go-common"
	"github.com/tepleton/go-merkle"
	"github.com/tepleton/go-wire"
	tmtypes "github.com/tepleton/tepleton/types"
)

func cmdABIRegisterTx(c *cli.Context) error {
	chainID := c.String("chain_id")
	genesisFile := c.String("genesis")
	parent := c.Parent()

	genesisBytes, err := ioutil.ReadFile(genesisFile)
	if err != nil {
		return errors.New(cmn.Fmt("Error reading genesis file %v: %v", genesisFile, err))
	}

	abiTx := abi.ABIRegisterChainTx{
		abi.BlockchainGenesis{
			ChainID: chainID,
			Genesis: string(genesisBytes),
		},
	}

	fmt.Println("ABITx:", string(wire.JSONBytes(abiTx)))

	data := []byte(wire.BinaryBytes(struct {
		abi.ABITx `json:"unwrap"`
	}{abiTx}))
	name := "ABI"

	return appTx(parent, name, data)
}

func cmdABIUpdateTx(c *cli.Context) error {
	headerBytes, err := hex.DecodeString(stripHex(c.String("header")))
	if err != nil {
		return errors.New(cmn.Fmt("Header (%v) is invalid hex: %v", c.String("header"), err))
	}
	commitBytes, err := hex.DecodeString(stripHex(c.String("commit")))
	if err != nil {
		return errors.New(cmn.Fmt("Commit (%v) is invalid hex: %v", c.String("commit"), err))
	}

	header := new(tmtypes.Header)
	commit := new(tmtypes.Commit)

	if err := wire.ReadBinaryBytes(headerBytes, &header); err != nil {
		return errors.New(cmn.Fmt("Error unmarshalling header: %v", err))
	}
	if err := wire.ReadBinaryBytes(commitBytes, &commit); err != nil {
		return errors.New(cmn.Fmt("Error unmarshalling commit: %v", err))
	}

	abiTx := abi.ABIUpdateChainTx{
		Header: *header,
		Commit: *commit,
	}

	fmt.Println("ABITx:", string(wire.JSONBytes(abiTx)))

	data := []byte(wire.BinaryBytes(struct {
		abi.ABITx `json:"unwrap"`
	}{abiTx}))
	name := "ABI"

	return appTx(c.Parent(), name, data)
}

func cmdABIPacketCreateTx(c *cli.Context) error {
	fromChain, toChain := c.String("from"), c.String("to")
	packetType := c.String("type")

	payloadBytes, err := hex.DecodeString(stripHex(c.String("payload")))
	if err != nil {
		return errors.New(cmn.Fmt("Payload (%v) is invalid hex: %v", c.String("payload"), err))
	}

	sequence, err := getABISequence(c)
	if err != nil {
		return err
	}

	abiTx := abi.ABIPacketCreateTx{
		Packet: abi.Packet{
			SrcChainID: fromChain,
			DstChainID: toChain,
			Sequence:   sequence,
			Type:       packetType,
			Payload:    payloadBytes,
		},
	}

	fmt.Println("ABITx:", string(wire.JSONBytes(abiTx)))

	data := []byte(wire.BinaryBytes(struct {
		abi.ABITx `json:"unwrap"`
	}{abiTx}))

	return appTx(c.Parent().Parent(), "ABI", data)
}

func cmdABIPacketPostTx(c *cli.Context) error {
	fromChain, fromHeight := c.String("from"), c.Int("height")

	packetBytes, err := hex.DecodeString(stripHex(c.String("packet")))
	if err != nil {
		return errors.New(cmn.Fmt("Packet (%v) is invalid hex: %v", c.String("packet"), err))
	}
	proofBytes, err := hex.DecodeString(stripHex(c.String("proof")))
	if err != nil {
		return errors.New(cmn.Fmt("Proof (%v) is invalid hex: %v", c.String("proof"), err))
	}

	var packet abi.Packet
	proof := new(merkle.IAVLProof)

	if err := wire.ReadBinaryBytes(packetBytes, &packet); err != nil {
		return errors.New(cmn.Fmt("Error unmarshalling packet: %v", err))
	}
	if err := wire.ReadBinaryBytes(proofBytes, &proof); err != nil {
		return errors.New(cmn.Fmt("Error unmarshalling proof: %v", err))
	}

	abiTx := abi.ABIPacketPostTx{
		FromChainID:     fromChain,
		FromChainHeight: uint64(fromHeight),
		Packet:          packet,
		Proof:           proof,
	}

	fmt.Println("ABITx:", string(wire.JSONBytes(abiTx)))

	data := []byte(wire.BinaryBytes(struct {
		abi.ABITx `json:"unwrap"`
	}{abiTx}))

	return appTx(c.Parent().Parent(), "ABI", data)
}

func getABISequence(c *cli.Context) (uint64, error) {
	if c.IsSet("sequence") {
		return uint64(c.Int("sequence")), nil
	}

	// TODO: get sequence
	return 0, nil
}
