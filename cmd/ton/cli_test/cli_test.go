package clitest

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton-sdk/client/keys"
	"github.com/tepleton/tepleton-sdk/cmd/ton/app"
	"github.com/tepleton/tepleton-sdk/server"
	"github.com/tepleton/tepleton-sdk/tests"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/stake"
)

func TestGaiaCLISend(t *testing.T) {

	tests.ExecuteT(t, "tond unsafe_reset_all")
	pass := "1234567890"
	executeWrite(t, "toncli keys delete foo", pass)
	executeWrite(t, "toncli keys delete bar", pass)
	chainID := executeInit(t, "tond init -o --name=foo")
	executeWrite(t, "toncli keys add bar", pass)

	// get a free port, also setup some common flags
	servAddr := server.FreeTCPAddr(t)
	flags := fmt.Sprintf("--node=%v --chain-id=%v", servAddr, chainID)

	// start tond server
	cmd, _, _ := tests.GoExecuteT(t, fmt.Sprintf("tond start --rpc.laddr=%v", servAddr))
	defer cmd.Process.Kill()

	fooAddr, _ := executeGetAddrPK(t, "toncli keys show foo --output=json")
	barAddr, _ := executeGetAddrPK(t, "toncli keys show bar --output=json")

	fooAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(50), fooAcc.GetCoins().AmountOf("steak"))

	executeWrite(t, fmt.Sprintf("toncli send %v --amount=10steak --to=%v --name=foo", flags, barAddr), pass)
	time.Sleep(time.Second * 3) // waiting for some blocks to pass

	barAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barAddr, flags))
	assert.Equal(t, int64(10), barAcc.GetCoins().AmountOf("steak"))
	fooAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(40), fooAcc.GetCoins().AmountOf("steak"))
}

func TestGaiaCLIDeclareCandidacy(t *testing.T) {

	tests.ExecuteT(t, "tond unsafe_reset_all")
	pass := "1234567890"
	executeWrite(t, "toncli keys delete foo", pass)
	executeWrite(t, "toncli keys delete bar", pass)
	chainID := executeInit(t, "tond init -o --name=foo")
	executeWrite(t, "toncli keys add bar", pass)

	// get a free port, also setup some common flags
	servAddr := server.FreeTCPAddr(t)
	flags := fmt.Sprintf("--node=%v --chain-id=%v", servAddr, chainID)

	// start tond server
	cmd, _, _ := tests.GoExecuteT(t, fmt.Sprintf("tond start --rpc.laddr=%v", servAddr))
	defer cmd.Process.Kill()

	fooAddr, _ := executeGetAddrPK(t, "toncli keys show foo --output=json")
	barAddr, barPubKey := executeGetAddrPK(t, "toncli keys show bar --output=json")

	executeWrite(t, fmt.Sprintf("toncli send %v --amount=10steak --to=%v --name=foo", flags, barAddr), pass)
	time.Sleep(time.Second * 3) // waiting for some blocks to pass

	fooAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(40), fooAcc.GetCoins().AmountOf("steak"))
	barAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barAddr, flags))
	assert.Equal(t, int64(10), barAcc.GetCoins().AmountOf("steak"))

	// declare candidacy
	declStr := fmt.Sprintf("toncli declare-candidacy %v", flags)
	declStr += fmt.Sprintf(" --name=%v", "bar")
	declStr += fmt.Sprintf(" --address-candidate=%v", barAddr)
	declStr += fmt.Sprintf(" --pubkey=%v", barPubKey)
	declStr += fmt.Sprintf(" --amount=%v", "3steak")
	declStr += fmt.Sprintf(" --moniker=%v", "bar-vally")
	fmt.Printf("debug declStr: %v\n", declStr)
	executeWrite(t, declStr, pass)
	time.Sleep(time.Second) // waiting for some blocks to pass
	barAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barAddr, flags))
	assert.Equal(t, int64(7), barAcc.GetCoins().AmountOf("steak"))
	candidate := executeGetCandidate(t, fmt.Sprintf("toncli candidate %v --address-candidate=%v", flags, barAddr))
	assert.Equal(t, candidate.Address.String(), barAddr)
	assert.Equal(t, int64(3), candidate.Assets.Evaluate())

	// TODO timeout issues if not connected to the internet
	// unbond a single share
	//unbondStr := fmt.Sprintf("toncli unbond %v", flags)
	//unbondStr += fmt.Sprintf(" --name=%v", "bar")
	//unbondStr += fmt.Sprintf(" --address-candidate=%v", barAddr)
	//unbondStr += fmt.Sprintf(" --address-delegator=%v", barAddr)
	//unbondStr += fmt.Sprintf(" --shares=%v", "1")
	//unbondStr += fmt.Sprintf(" --sequence=%v", "1")
	//fmt.Printf("debug unbondStr: %v\n", unbondStr)
	//executeWrite(t, unbondStr, pass)
	//time.Sleep(time.Second * 3) // waiting for some blocks to pass
	//barAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barAddr, flags))
	//assert.Equal(t, int64(99998), barAcc.GetCoins().AmountOf("steak"))
	//candidate = executeGetCandidate(t, fmt.Sprintf("toncli candidate %v --address-candidate=%v", flags, barAddr))
	//assert.Equal(t, int64(2), candidate.Assets.Evaluate())
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) {
	cmd, wc, _ := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := wc.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	fmt.Printf("debug waiting cmdStr: %v\n", cmdStr)
	cmd.Wait()
}

func executeWritePrint(t *testing.T, cmdStr string, writes ...string) {
	cmd, wc, rc := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := wc.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	fmt.Printf("debug waiting cmdStr: %v\n", cmdStr)
	cmd.Wait()

	bz := make([]byte, 100000)
	rc.Read(bz)
	fmt.Printf("debug read: %v\n", string(bz))
}

func executeInit(t *testing.T, cmdStr string) (chainID string) {
	out := tests.ExecuteT(t, cmdStr)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	return
}

func executeGetAddrPK(t *testing.T, cmdStr string) (addr, pubKey string) {
	out := tests.ExecuteT(t, cmdStr)
	var ko keys.KeyOutput
	keys.UnmarshalJSON([]byte(out), &ko)
	return ko.Address, ko.PubKey
}

func executeGetAccount(t *testing.T, cmdStr string) auth.BaseAccount {
	out := tests.ExecuteT(t, cmdStr)
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)
	value := initRes["value"]
	var acc auth.BaseAccount
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	err = cdc.UnmarshalJSON(value, &acc)
	require.NoError(t, err, "value %v, err %v", string(value), err)
	return acc
}

func executeGetCandidate(t *testing.T, cmdStr string) stake.Candidate {
	out := tests.ExecuteT(t, cmdStr)
	var candidate stake.Candidate
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &candidate)
	require.NoError(t, err, "out %v, err %v", out, err)
	return candidate
}
