package rest

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/stake"
)

func registerQueryRoutes(ctx context.CoreContext, r *mux.Router, cdc *wire.Codec) {
	r.HandleFunc(
		"/stake/{delegator}/bonding_status/{validator}",
		bondingStatusHandlerFn(ctx, "stake", cdc),
	).Methods("GET")
	r.HandleFunc(
		"/stake/validators",
		validatorsHandlerFn(ctx, "stake", cdc),
	).Methods("GET")
}

// http request handler to query delegator bonding status
func bondingStatusHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read parameters
		vars := mux.Vars(r)
		delegator := vars["delegator"]
		validator := vars["validator"]

		bz, err := hex.DecodeString(delegator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		delegatorAddr := sdk.Address(bz)

		bz, err = hex.DecodeString(validator)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		validatorAddr := sdk.Address(bz)

		key := stake.GetDelegationKey(delegatorAddr, validatorAddr, cdc)

		res, err := ctx.Query(key, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't query bond. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there is no data for this bond
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var bond stake.Delegation
		err = cdc.UnmarshalBinary(res, &bond)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode bond. Error: %s", err.Error())))
			return
		}

		output, err := cdc.MarshalJSON(bond)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}

<<<<<<< HEAD
=======
// TODO move exist next to validator struct for maintainability
type StakeValidatorOutput struct {
	Owner   string `json:"owner"`   // in bech32
	PubKey  string `json:"pub_key"` // in bech32
	Revoked bool   `json:"revoked"` // has the validator been revoked from bonded status?

	PoolShares      stake.PoolShares `json:"pool_shares"`      // total shares for tokens held in the pool
	DelegatorShares sdk.Rat          `json:"delegator_shares"` // total shares issued to a validator's delegators

	Description        stake.Description `json:"description"`           // description terms for the validator
	BondHeight         int64             `json:"bond_height"`           // earliest height as a bonded validator
	BondIntraTxCounter int16             `json:"bond_intra_tx_counter"` // block-local tx index of validator change
	ProposerRewardPool sdk.Coins         `json:"proposer_reward_pool"`  // XXX reward pool collected from being the proposer

	Commission            sdk.Rat `json:"commission"`              // XXX the commission rate of fees charged to any delegators
	CommissionMax         sdk.Rat `json:"commission_max"`          // XXX maximum commission rate which this validator can ever charge
	CommissionChangeRate  sdk.Rat `json:"commission_change_rate"`  // XXX maximum daily increase of the validator commission
	CommissionChangeToday sdk.Rat `json:"commission_change_today"` // XXX commission rate change today, reset each day (UTC time)

	// fee related
	PrevBondedShares sdk.Rat `json:"prev_bonded_shares"` // total shares of a global hold pools
}

func bech32StakeValidatorOutput(validator stake.Validator) (StakeValidatorOutput, error) {
	bechOwner, err := sdk.Bech32ifyVal(validator.Owner)
	if err != nil {
		return StakeValidatorOutput{}, err
	}
	bechValPubkey, err := sdk.Bech32ifyValPub(validator.PubKey)
	if err != nil {
		return StakeValidatorOutput{}, err
	}

	return StakeValidatorOutput{
		Owner:   bechOwner,
		PubKey:  bechValPubkey,
		Revoked: validator.Revoked,

		PoolShares:      validator.PoolShares,
		DelegatorShares: validator.DelegatorShares,

		Description:        validator.Description,
		BondHeight:         validator.BondHeight,
		BondIntraTxCounter: validator.BondIntraTxCounter,
		ProposerRewardPool: validator.ProposerRewardPool,

		Commission:            validator.Commission,
		CommissionMax:         validator.CommissionMax,
		CommissionChangeRate:  validator.CommissionChangeRate,
		CommissionChangeToday: validator.CommissionChangeToday,

		PrevBondedShares: validator.PrevBondedShares,
	}, nil
}

// TODO bech32
>>>>>>> dev
// http request handler to query list of validators
func validatorsHandlerFn(ctx context.CoreContext, storeName string, cdc *wire.Codec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kvs, err := ctx.QuerySubspace(cdc, stake.ValidatorsKey, storeName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't query validators. Error: %s", err.Error())))
			return
		}

		// the query will return empty if there are no validators
		if len(kvs) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// parse out the validators
		validators := make([]stake.Validator, len(kvs))
		for i, kv := range kvs {
			var validator stake.Validator
			err = cdc.UnmarshalBinary(kv.Value, &validator)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Couldn't decode validator. Error: %s", err.Error())))
				return
			}
			validators[i] = validator
		}

		output, err := cdc.MarshalJSON(validators)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(output)
	}
}
