package coin

import (
	"fmt"

	"github.com/tepleton/go-wire/data"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/stack"
	"github.com/tepleton/basecoin/types"
)

const (
	NameCoin = "coin"
)

// Handler writes
type Handler struct {
	Accountant
}

var _ basecoin.Handler = Handler{}

func NewHandler() Handler {
	return Handler{
		Accountant: NewAccountant(""),
	}
}

func (_ Handler) Name() string {
	return NameCoin
}

// CheckTx checks if there is enough money in the account
func (h Handler) CheckTx(ctx basecoin.Context, store types.KVStore, tx basecoin.Tx) (res basecoin.Result, err error) {
	send, err := checkTx(ctx, tx)
	if err != nil {
		return res, err
	}

	// now make sure there is money
	for _, in := range send.Inputs {
		_, err = h.CheckCoins(store, in.Address, in.Coins.Negative(), in.Sequence)
		if err != nil {
			return res, err
		}
	}

	// otherwise, we are good
	return res, nil
}

// DeliverTx moves the money
func (h Handler) DeliverTx(ctx basecoin.Context, store types.KVStore, tx basecoin.Tx) (res basecoin.Result, err error) {
	send, err := checkTx(ctx, tx)
	if err != nil {
		return res, err
	}

	// deduct from all input accounts
	for _, in := range send.Inputs {
		_, err = h.ChangeCoins(store, in.Address, in.Coins.Negative(), in.Sequence)
		if err != nil {
			return res, err
		}
	}

	// add to all output accounts
	for _, out := range send.Outputs {
		// note: sequence number is ignored when adding coins, only checked for subtracting
		_, err = h.ChangeCoins(store, out.Address, out.Coins, 0)
		if err != nil {
			return res, err
		}
	}

	// a-ok!
	return basecoin.Result{}, nil
}

func (h Handler) SetOption(l log.Logger, store types.KVStore, module, key, value string) (log string, err error) {
	if module != NameCoin {
		return "", errors.ErrUnknownModule(module)
	}
	if key == "account" {
		var acc GenesisAccount
		err = data.FromJSON([]byte(value), &acc)
		if err != nil {
			return "", err
		}
		acc.Balance.Sort()
		addr, err := acc.GetAddr()
		if err != nil {
			return "", ErrInvalidAddress()
		}
		// this sets the permission for a public key signature, use that app
		actor := stack.SigPerm(addr)
		err = storeAccount(store, h.MakeKey(actor), acc.ToAccount())
		if err != nil {
			return "", err
		}
		return "Success", nil

	} else {
		msg := fmt.Sprintf("Unknown key: %s", key)
		return "", errors.ErrInternal(msg)
	}
}

func checkTx(ctx basecoin.Context, tx basecoin.Tx) (send SendTx, err error) {
	// check if the tx is proper type and valid
	send, ok := tx.Unwrap().(SendTx)
	if !ok {
		return send, errors.ErrInvalidFormat(tx)
	}
	err = send.ValidateBasic()
	if err != nil {
		return send, err
	}

	// check if all inputs have permission
	for _, in := range send.Inputs {
		if !ctx.HasPermission(in.Address) {
			return send, errors.ErrUnauthorized()
		}
	}
	return send, nil
}