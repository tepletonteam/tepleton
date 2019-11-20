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
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/stake"
	crypto "github.com/tepleton/go-crypto"
)

func TestGaiaCLISend(t *testing.T) {
	fmt.Println("wackydebugoutput TestGaiaCLISend 0")

	tests.ExecuteT(t, "tond unsafe_reset_all")
	pass := "1234567890"
	executeWrite(t, "toncli keys delete foo", pass)
	executeWrite(t, "toncli keys delete bar", pass)
	chainID := executeInit(t, "tond init -o --name=foo")
	executeWrite(t, "toncli keys add bar", pass)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--node=%v --chain-id=%v", servAddr, chainID)

	// start tond server
	proc := tests.GoExecuteT(t, fmt.Sprintf("tond start --rpc.laddr=%v", servAddr))
	defer proc.Stop(false)
	tests.WaitForStart(port)

	fooAddr, _ := executeGetAddrPK(t, "toncli keys show foo --output=json")
	fooCech, err := sdk.Bech32TepletonifyAcc(fooAddr)
	require.NoError(t, err)
	barAddr, _ := executeGetAddrPK(t, "toncli keys show bar --output=json")
	barCech, err := sdk.Bech32TepletonifyAcc(barAddr)
	require.NoError(t, err)

	fooAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooCech, flags))
	assert.Equal(t, int64(50), fooAcc.GetCoins().AmountOf("steak"))

	executeWrite(t, fmt.Sprintf("toncli send %v --amount=10steak --to=%v --name=foo", flags, barCech), pass)
	fmt.Println("wackydebugoutput TestGaiaCLISend 1")
	fmt.Println("wackydebugoutput TestGaiaCLISend 2")
	time.Sleep(time.Second * 2) // waiting for some blocks to pass

	barAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barCech, flags))
	assert.Equal(t, int64(10), barAcc.GetCoins().AmountOf("steak"))
	fooAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooCech, flags))
	assert.Equal(t, int64(40), fooAcc.GetCoins().AmountOf("steak"))

	// test autosequencing
	executeWrite(t, fmt.Sprintf("toncli send %v --amount=10steak --to=%v --name=foo", flags, barCech), pass)
	time.Sleep(time.Second * 2) // waiting for some blocks to pass

	barAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barCech, flags))
	assert.Equal(t, int64(20), barAcc.GetCoins().AmountOf("steak"))
	fooAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooCech, flags))
	assert.Equal(t, int64(30), fooAcc.GetCoins().AmountOf("steak"))
}

func TestGaiaCLICreateValidator(t *testing.T) {

	tests.ExecuteT(t, "tond unsafe_reset_all")
	pass := "1234567890"
	executeWrite(t, "toncli keys delete foo", pass)
	executeWrite(t, "toncli keys delete bar", pass)
	chainID := executeInit(t, "tond init -o --name=foo")
	executeWrite(t, "toncli keys add bar", pass)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--node=%v --chain-id=%v", servAddr, chainID)

	// start tond server
	proc := tests.GoExecuteT(t, fmt.Sprintf("tond start --rpc.laddr=%v", servAddr))
	defer proc.Stop(false)
	tests.WaitForStart(port)

	fooAddr, _ := executeGetAddrPK(t, "toncli keys show foo --output=json")
	fooCech, err := sdk.Bech32TepletonifyAcc(fooAddr)
	require.NoError(t, err)
	barAddr, barPubKey := executeGetAddrPK(t, "toncli keys show bar --output=json")
	barCech, err := sdk.Bech32TepletonifyAcc(barAddr)
	require.NoError(t, err)
	barCeshPubKey, err := sdk.Bech32TepletonifyValPub(barPubKey)
	require.NoError(t, err)

	executeWrite(t, fmt.Sprintf("toncli send %v --amount=10steak --to=%v --name=foo", flags, barCech), pass)
	time.Sleep(time.Second * 2) // waiting for some blocks to pass

	barAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barCech, flags))
	assert.Equal(t, int64(10), barAcc.GetCoins().AmountOf("steak"))
	fooAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooCech, flags))
	assert.Equal(t, int64(40), fooAcc.GetCoins().AmountOf("steak"))

	// create validator
	cvStr := fmt.Sprintf("toncli create-validator %v", flags)
	cvStr += fmt.Sprintf(" --name=%v", "bar")
	cvStr += fmt.Sprintf(" --address-validator=%v", barCech)
	cvStr += fmt.Sprintf(" --pubkey=%v", barCeshPubKey)
	cvStr += fmt.Sprintf(" --amount=%v", "2steak")
	cvStr += fmt.Sprintf(" --moniker=%v", "bar-vally")

	executeWrite(t, cvStr, pass)
	time.Sleep(time.Second * 3) // waiting for some blocks to pass

	barAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barCech, flags))
	require.Equal(t, int64(8), barAcc.GetCoins().AmountOf("steak"), "%v", barAcc)

	validator := executeGetValidator(t, fmt.Sprintf("toncli validator %v --output=json %v", barCech, flags))
	assert.Equal(t, validator.Owner, barAddr)
	assert.Equal(t, "2/1", validator.PoolShares.Amount.String())

	// unbond a single share
	unbondStr := fmt.Sprintf("toncli unbond %v", flags)
	unbondStr += fmt.Sprintf(" --name=%v", "bar")
	unbondStr += fmt.Sprintf(" --address-validator=%v", barCech)
	unbondStr += fmt.Sprintf(" --address-delegator=%v", barCech)
	unbondStr += fmt.Sprintf(" --shares=%v", "1")
	unbondStr += fmt.Sprintf(" --sequence=%v", "1")
	t.Log(fmt.Sprintf("debug unbondStr: %v\n", unbondStr))

	executeWrite(t, unbondStr, pass)
	time.Sleep(time.Second * 3) // waiting for some blocks to pass

	barAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barCech, flags))
	require.Equal(t, int64(9), barAcc.GetCoins().AmountOf("steak"), "%v", barAcc)
	validator = executeGetValidator(t, fmt.Sprintf("toncli validator %v --output=json %v", barCech, flags))
	assert.Equal(t, "1/1", validator.PoolShares.Amount.String())
}

//___________________________________________________________________________________
// executors

func executeWrite(t *testing.T, cmdStr string, writes ...string) {
	proc := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	proc.Wait()
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

func executeGetAddrPK(t *testing.T, cmdStr string) (sdk.Address, crypto.PubKey) {
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

func executeGetValidator(t *testing.T, cmdStr string) stake.Validator {
	out := tests.ExecuteT(t, cmdStr)
	var validator stake.Validator
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &validator)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return validator
}
