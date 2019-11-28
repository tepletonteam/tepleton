package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tepleton/go-crypto/keys"
	ctypes "github.com/tepleton/tepleton/rpc/core/types"

	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/stake"
)

func registerTxRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/stake/delegations",
		editDelegationsRequestHandlerFn(cdc, kb, ctx),
	).Methods("POST")
}

type editDelegationsBody struct {
<<<<<<< HEAD
	LocalAccountName string              `json:"name"`
	Password         string              `json:"password"`
	ChainID          string              `json:"chain_id"`
	Sequence         int64               `json:"sequence"`
	Delegate         []stake.MsgDelegate `json:"delegate"`
	Unbond           []stake.MsgUnbond   `json:"unbond"`
=======
	LocalAccountName string             `json:"name"`
	Password         string             `json:"password"`
	ChainID          string             `json:"chain_id"`
	AccountNumber    int64              `json:"account_number"`
	Sequence         int64              `json:"sequence"`
	Gas              int64              `json:"gas"`
	Delegate         []msgDelegateInput `json:"delegate"`
	Unbond           []msgUnbondInput   `json:"unbond"`
>>>>>>> dev
}

func editDelegationsRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, ctx context.CoreContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m editDelegationsBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.Unmarshal(body, &m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		// build messages
		messages := make([]sdk.Msg, len(m.Delegate)+len(m.Unbond))
		i := 0
		for _, msg := range m.Delegate {
			if !bytes.Equal(info.Address(), msg.DelegatorAddr) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Must use own delegator address"))
				return
			}
			messages[i] = msg
			i++
		}
		for _, msg := range m.Unbond {
			if !bytes.Equal(info.Address(), msg.DelegatorAddr) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Must use own delegator address"))
				return
			}
			messages[i] = msg
			i++
		}

		// add gas to context
		ctx = ctx.WithGas(m.Gas)

		// sign messages
		signedTxs := make([][]byte, len(messages[:]))
		for i, msg := range messages {
			// increment sequence for each message
			ctx = ctx.WithAccountNumber(m.AccountNumber)
			ctx = ctx.WithSequence(m.Sequence)
			m.Sequence++

			txBytes, err := ctx.SignAndBuild(m.LocalAccountName, m.Password, msg, cdc)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
				return
			}

			signedTxs[i] = txBytes
		}

		// send
		// XXX the operation might not be atomic if a tx fails
		//     should we have a sdk.MultiMsg type to make sending atomic?
		results := make([]*ctypes.ResultBroadcastTxCommit, len(signedTxs[:]))
		for i, txBytes := range signedTxs {
			res, err := ctx.BroadcastTx(txBytes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			results[i] = res
		}

		output, err := json.MarshalIndent(results[:], "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
