package auth

import (
	"fmt"
	"testing"

	crypto "github.com/tepleton/go-crypto"
	cmn "github.com/tepleton/tmlibs/common"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/stack"
	"github.com/tepleton/basecoin/state"
)

func makeSignTx() basecoin.Tx {
	key := crypto.GenPrivKeyEd25519().Wrap()
	payload := cmn.RandBytes(32)
	tx := NewSig(stack.NewRawTx(payload))
	Sign(tx, key)
	return tx.Wrap()
}

func makeMultiSignTx(cnt int) basecoin.Tx {
	payload := cmn.RandBytes(32)
	tx := NewMulti(stack.NewRawTx(payload))
	for i := 0; i < cnt; i++ {
		key := crypto.GenPrivKeyEd25519().Wrap()
		Sign(tx, key)
	}
	return tx.Wrap()
}

func makeHandler() basecoin.Handler {
	return stack.New(Signatures{}).Use(stack.OKHandler{})
}

func BenchmarkCheckOneSig(b *testing.B) {
	tx := makeSignTx()
	h := makeHandler()
	store := state.NewMemKVStore()
	for i := 1; i <= b.N; i++ {
		ctx := stack.NewContext("foo", log.NewNopLogger())
		_, err := h.DeliverTx(ctx, store, tx)
		// never should error
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkCheckMultiSig(b *testing.B) {
	sigs := []int{1, 3, 8, 20}
	for _, cnt := range sigs {
		label := fmt.Sprintf("%dsigs", cnt)
		b.Run(label, func(sub *testing.B) {
			benchmarkCheckMultiSig(sub, cnt)
		})
	}
}

func benchmarkCheckMultiSig(b *testing.B, cnt int) {
	tx := makeMultiSignTx(cnt)
	h := makeHandler()
	store := state.NewMemKVStore()
	for i := 1; i <= b.N; i++ {
		ctx := stack.NewContext("foo", log.NewNopLogger())
		_, err := h.DeliverTx(ctx, store, tx)
		// never should error
		if err != nil {
			panic(err)
		}
	}
}
