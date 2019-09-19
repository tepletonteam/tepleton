package rest

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"

	wrsp "github.com/tepleton/wrsp/types"
	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/client/commands"
	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/modules/base"
	"github.com/tepleton/basecoin/modules/nonce"
	"github.com/tepleton/basecoin/modules/roles"
	"github.com/tepleton/tmlibs/common"
)

type RoleInput struct {
	Role string `json:"role" validate:"required,min=2"`

	MinimumSigners uint32 `json:"min_sigs" validate:"required,min=1"`

	Signers []basecoin.Actor `json:"signers" validate:"required,min=1"`

	Sequence uint32 `json:"seq" validate:"required,min=1"`
}

func parseRole(roleInHex string) ([]byte, error) {
	parsedRole, err := hex.DecodeString(common.StripHex(roleInHex))
	if err != nil {
		err = errors.WithMessage("invalid hex", err, wrsp.CodeType_EncodingError)
		return nil, err
	}
	return parsedRole, nil
}

// mux.Router registrars

// RegisterQueryAccount is a mux.Router handler that exposes GET
// method access on route /query/account/{signature} to query accounts
func RegisterCreateRole(r *mux.Router) error {
	r.HandleFunc("/build/create_role", doCreateRole).Methods("POST")
	return nil
}

func doCreateRole(w http.ResponseWriter, r *http.Request) {
	ri := new(RoleInput)
	if err := common.ParseRequestAndValidateJSON(r, ri); err != nil {
		common.WriteError(w, err)
		return
	}

	parsedRole, err := parseRole(ri.Role)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	// Note the ordering of Tx wrapping matters:
	// 1. NonceTx
	tx := (nonce.Tx{}).Wrap()
	tx = nonce.NewTx(ri.Sequence, ri.Signers, tx)

	// 2. CreateRoleTx
	tx = roles.NewCreateRoleTx(parsedRole, ri.MinimumSigners, ri.Signers)

	// 3. ChainTx
	tx = base.NewChainTx(commands.GetChainID(), 0, tx)

	common.WriteSuccess(w, tx)
}

// End of mux.Router registrars
