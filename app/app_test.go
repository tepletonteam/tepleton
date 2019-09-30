package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/go-crypto"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"

	"github.com/tepleton/tepleton-sdk/store"
	sdk "github.com/tepleton/tepleton-sdk/types"
)

// A mock transaction to update a validator's voting power.
type testTx struct {
	Addr     []byte
	NewPower int64
}

const txType = "testTx"

func (tx testTx) Type() string                            { return txType }
func (tx testTx) Get(key interface{}) (value interface{}) { return nil }
func (tx testTx) GetSignBytes() []byte                    { return nil }
func (tx testTx) ValidateBasic() error                    { return nil }
func (tx testTx) GetSigners() []crypto.Address            { return nil }
func (tx testTx) GetFeePayer() crypto.Address             { return nil }
func (tx testTx) GetSignatures() []sdk.StdSignature       { return nil }

func TestBasic(t *testing.T) {
	store, storeKeys := newCommitMultiStore()

	// Create app.
	app := NewApp(t.Name(), store)
	app.SetTxDecoder(func(txBytes []byte) (sdk.Tx, error) {
		var ttx testTx
		fromJSON(txBytes, &ttx)
		return ttx, nil
	})

	app.SetDefaultAnteHandler(func(ctx sdk.Context, tx sdk.Tx) (newCtx sdk.Context, res sdk.Result, abort bool) { return })
	app.Router().AddRoute(txType, func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		// TODO
		return sdk.Result{}
	})

	// Load latest state, which should be empty.
	err := app.LoadLatestVersion(storeKeys["main"])
	assert.Nil(t, err)
	assert.Equal(t, app.LastBlockHeight(), int64(0))

	// Create the validators
	var numVals = 3
	var valSet = make([]wrsp.Validator, numVals)
	for i := 0; i < numVals; i++ {
		valSet[i] = makeVal(secret(i))
	}

	// Initialize the chain
	app.InitChain(wrsp.RequestInitChain{
		Validators: valSet,
	})

	// Simulate the start of a block.
	app.BeginBlock(wrsp.RequestBeginBlock{})

	// Add 1 to each validator's voting power.
	for i, val := range valSet {
		tx := testTx{
			Addr:     makePubKey(secret(i)).Address(),
			NewPower: val.Power + 1,
		}
		txBytes := toJSON(tx)
		res := app.DeliverTx(txBytes)
		assert.True(t, res.IsOK(), "%#v", res)
	}

	// Simulate the end of a block.
	// Get the summary of validator updates.
	res := app.EndBlock(wrsp.RequestEndBlock{})
	valUpdates := res.ValidatorUpdates

	// Assert that validator updates are correct.
	for _, val := range valSet {
		// Sanity
		assert.NotEqual(t, len(val.PubKey), 0)

		// Find matching update and splice it out.
		for j := 0; j < len(valUpdates); {
			valUpdate := valUpdates[j]

			// Matched.
			if bytes.Equal(valUpdate.PubKey, val.PubKey) {
				assert.Equal(t, valUpdate.Power, val.Power+1)
				if j < len(valUpdates)-1 {
					// Splice it out.
					valUpdates = append(valUpdates[:j], valUpdates[j+1:]...)
				}
				break
			}

			// Not matched.
			j += 1
		}
	}
	assert.Equal(t, len(valUpdates), 0, "Some validator updates were unexpected")
}

//----------------------------------------

func randPower() int64 {
	return cmn.RandInt64()
}

func makeVal(secret string) wrsp.Validator {
	return wrsp.Validator{
		PubKey: makePubKey(secret).Bytes(),
		Power:  randPower(),
	}
}

func makePubKey(secret string) crypto.PubKey {
	return makePrivKey(secret).PubKey()
}

func makePrivKey(secret string) crypto.PrivKey {
	privKey := crypto.GenPrivKeyEd25519FromSecret([]byte(secret))
	return privKey
}

func secret(index int) string {
	return fmt.Sprintf("secret%d", index)
}

func copyVal(val wrsp.Validator) wrsp.Validator {
	// val2 := *val
	// return &val2
	return val
}

func toJSON(o interface{}) []byte {
	bz, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	// fmt.Println(">> toJSON:", string(bz))
	return bz
}

func fromJSON(bz []byte, ptr interface{}) {
	// fmt.Println(">> fromJSON:", string(bz))
	err := json.Unmarshal(bz, ptr)
	if err != nil {
		panic(err)
	}
}

// Creates a sample CommitMultiStore
func newCommitMultiStore() (sdk.CommitMultiStore, map[string]sdk.SubstoreKey) {
	dbMain := dbm.NewMemDB()
	dbXtra := dbm.NewMemDB()
	keyMain := sdk.NewKVStoreKey("main")
	keyXtra := sdk.NewKVStoreKey("xtra")
	ms := store.NewCommitMultiStore(dbMain) // Also store rootMultiStore metadata here (it shouldn't clash)
	ms.SetSubstoreLoader(keyMain, store.NewIAVLStoreLoader(dbMain, 0, 0))
	ms.SetSubstoreLoader(keyXtra, store.NewIAVLStoreLoader(dbXtra, 0, 0))
	return ms, map[string]sdk.SubstoreKey{
		"main": keyMain,
		"xtra": keyXtra,
	}
}
