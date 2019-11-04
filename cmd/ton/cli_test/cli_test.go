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
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/stake"
)

func TestGaiaCLISend(t *testing.T) {

	tests.ExecuteT(t, "tond unsafe_reset_all", 1)
	pass := "1234567890"
	executeWrite(t, "toncli keys delete foo", pass)
	executeWrite(t, "toncli keys delete bar", pass)
	keys, chainID := executeInit(t, "tond init -o --accounts=foo-100000fermion-true", "foo")
	require.Equal(t, 1, len(keys))

	// get a free port, also setup some common flags
	servAddr := server.FreeTCPAddr(t)
	flags := fmt.Sprintf("--node=%v --chain-id=%v", servAddr, chainID)

	// start tond server
	cmd, _, _ := tests.GoExecuteT(t, fmt.Sprintf("tond start --rpc.laddr=%v", servAddr))
	defer cmd.Process.Kill()

	executeWrite(t, "toncli keys add foo --recover", pass, keys[0])
	executeWrite(t, "toncli keys add bar", pass)

	fooAddr, _ := executeGetAddrPK(t, "toncli keys show foo --output=json")
	barAddr, _ := executeGetAddrPK(t, "toncli keys show bar --output=json")

	fooAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(100000), fooAcc.GetCoins().AmountOf("fermion"))

	executeWrite(t, fmt.Sprintf("toncli send %v --amount=10fermion --to=%v --name=foo", flags, barAddr), pass)
	time.Sleep(time.Second * 3) // waiting for some blocks to pass

	barAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", barAddr, flags))
	assert.Equal(t, int64(10), barAcc.GetCoins().AmountOf("fermion"))
	fooAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(99990), fooAcc.GetCoins().AmountOf("fermion"))
}

func TestGaiaCLIDeclareCandidacy(t *testing.T) {

	tests.ExecuteT(t, "tond unsafe_reset_all", 1)
	pass := "1234567890"
	executeWrite(t, "toncli keys delete foo", pass)
	keys, chainID := executeInit(t, "tond init -o --accounts=bar-100000fermion-true;foo-100000fermion-true", "bar", "foo")
	require.Equal(t, 2, len(keys))

	// get a free port, also setup some common flags
	servAddr := server.FreeTCPAddr(t)
	flags := fmt.Sprintf("--node=%v --chain-id=%v", servAddr, chainID)

	// start tond server
	cmd, _, _ := tests.GoExecuteT(t, fmt.Sprintf("tond start --rpc.laddr=%v", servAddr))
	defer cmd.Process.Kill()

	executeWrite(t, "toncli keys add foo --recover", pass, keys[1])
	fooAddr, fooPubKey := executeGetAddrPK(t, "toncli keys show foo --output=json")
	fooAcc := executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(100000), fooAcc.GetCoins().AmountOf("fermion"))

	// declare candidacy
	declStr := fmt.Sprintf("toncli declare-candidacy %v", flags)
	declStr += fmt.Sprintf(" --name=%v", "foo")
	declStr += fmt.Sprintf(" --address-candidate=%v", fooAddr)
	declStr += fmt.Sprintf(" --pubkey=%v", fooPubKey)
	declStr += fmt.Sprintf(" --amount=%v", "3fermion")
	declStr += fmt.Sprintf(" --moniker=%v", "foo-vally")
	fmt.Printf("debug declStr: %v\n", declStr)
	executeWrite(t, declStr, pass)
	time.Sleep(time.Second * 3) // waiting for some blocks to pass
	fooAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	assert.Equal(t, int64(99997), fooAcc.GetCoins().AmountOf("fermion"))
	candidate := executeGetCandidate(t, fmt.Sprintf("toncli candidate %v --address-candidate=%v", flags, fooAddr))
	assert.Equal(t, candidate.Address.String(), fooAddr)
	assert.Equal(t, int64(3), candidate.Assets.Evaluate())

	// TODO timeout issues if not connected to the internet
	// unbond a single share
	//unbondStr := fmt.Sprintf("toncli unbond %v", flags)
	//unbondStr += fmt.Sprintf(" --name=%v", "foo")
	//unbondStr += fmt.Sprintf(" --address-candidate=%v", fooAddr)
	//unbondStr += fmt.Sprintf(" --address-delegator=%v", fooAddr)
	//unbondStr += fmt.Sprintf(" --shares=%v", "1")
	//unbondStr += fmt.Sprintf(" --sequence=%v", "1")
	//fmt.Printf("debug unbondStr: %v\n", unbondStr)
	//executeWrite(t, unbondStr, pass)
	//time.Sleep(time.Second * 3) // waiting for some blocks to pass
	//fooAcc = executeGetAccount(t, fmt.Sprintf("toncli account %v %v", fooAddr, flags))
	//assert.Equal(t, int64(99998), fooAcc.GetCoins().AmountOf("fermion"))
	//candidate = executeGetCandidate(t, fmt.Sprintf("toncli candidate %v --address-candidate=%v", flags, fooAddr))
	//assert.Equal(t, int64(2), candidate.Assets.Evaluate())
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) {
	cmd, wc, _ := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := wc.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	cmd.Wait()
}

func executeWritePrint(t *testing.T, cmdStr string, writes ...string) {
	cmd, wc, rc := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := wc.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	cmd.Wait()

	bz := make([]byte, 100000)
	rc.Read(bz)
	fmt.Printf("debug read: %v\n", string(bz))
}

func executeInit(t *testing.T, cmdStr string, names ...string) (keys []string, chainID string) {
	out := tests.ExecuteT(t, cmdStr, 1)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	var appMessageRes map[string]json.RawMessage
	err = json.Unmarshal(initRes["app_message"], &appMessageRes)
	require.NoError(t, err)

	for _, name := range names {
		var key string
		err = json.Unmarshal(appMessageRes["secret-"+name], &key)
		require.NoError(t, err)
		keys = append(keys, key)
	}

	return
}

func executeGetAddrPK(t *testing.T, cmdStr string) (addr, pubKey string) {
	out := tests.ExecuteT(t, cmdStr, 2)
	var ko keys.KeyOutput
	keys.UnmarshalJSON([]byte(out), &ko)
	return ko.Address, ko.PubKey
}

func executeGetAccount(t *testing.T, cmdStr string) auth.BaseAccount {
	out := tests.ExecuteT(t, cmdStr, 2)
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)
	value := initRes["value"]
	var acc auth.BaseAccount
	_ = json.Unmarshal(value, &acc) //XXX pubkey can't be decoded go amino issue
	require.NoError(t, err, "value %v, err %v", string(value), err)
	return acc
}

func executeGetCandidate(t *testing.T, cmdStr string) stake.Candidate {
	out := tests.ExecuteT(t, cmdStr, 2)
	var candidate stake.Candidate
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &candidate)
	require.NoError(t, err, "out %v, err %v", out, err)
	return candidate
}
