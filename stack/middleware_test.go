package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tepleton/go-wire/data"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/state"
)

const (
	nameSigner = "signer"
)

func TestPermissionSandbox(t *testing.T) {
	require := require.New(t)

	// generic args
	ctx := NewContext("test-chain", 20, log.NewNopLogger())
	store := state.NewMemKVStore()
	raw := NewRawTx([]byte{1, 2, 3, 4})
	rawBytes, err := data.ToWire(raw)
	require.Nil(err)

	// test cases to make sure permissioning is solid
	grantee := basecoin.Actor{App: NameGrant, Address: []byte{1}}
	grantee2 := basecoin.Actor{App: NameGrant, Address: []byte{2}}
	// abi and grantee are the same, just different chains
	abi := basecoin.Actor{ChainID: "other", App: NameGrant, Address: []byte{1}}
	abi2 := basecoin.Actor{ChainID: "other", App: nameSigner, Address: []byte{21}}
	signer := basecoin.Actor{App: nameSigner, Address: []byte{21}}
	cases := []struct {
		asABI       bool
		grant       basecoin.Actor
		require     basecoin.Actor
		expectedRes data.Bytes
		expected    func(error) bool
	}{
		// grant as normal app middleware
		{false, grantee, grantee, rawBytes, nil},
		{false, grantee, grantee2, nil, errors.IsUnauthorizedErr},
		{false, grantee2, grantee2, rawBytes, nil},
		{false, abi, grantee, nil, errors.IsInternalErr},
		{false, grantee, abi, nil, errors.IsUnauthorizedErr},
		{false, grantee, signer, nil, errors.IsUnauthorizedErr},
		{false, signer, signer, nil, errors.IsInternalErr},

		// grant as abi middleware
		{true, abi, abi, rawBytes, nil},   // abi can set permissions
		{true, abi2, abi2, rawBytes, nil}, // for any app
		// the must match, both app and chain
		{true, abi, abi2, nil, errors.IsUnauthorizedErr},
		{true, abi, grantee, nil, errors.IsUnauthorizedErr},
		// cannot set local apps from abi middleware
		{true, grantee, grantee, nil, errors.IsInternalErr},
	}

	for i, tc := range cases {
		app := New(Recovery{})
		if tc.asABI {
			app = app.ABI(GrantMiddleware{Auth: tc.grant})
		} else {
			app = app.Apps(GrantMiddleware{Auth: tc.grant})
		}
		app = app.
			Apps(CheckMiddleware{Required: tc.require}).
			Use(EchoHandler{})

		res, err := app.CheckTx(ctx, store, raw)
		checkPerm(t, i, tc.expectedRes, tc.expected, res, err)

		res, err = app.DeliverTx(ctx, store, raw)
		checkPerm(t, i, tc.expectedRes, tc.expected, res, err)
	}
}

func checkPerm(t *testing.T, idx int, data []byte, check func(error) bool, res basecoin.Result, err error) {
	assert := assert.New(t)

	if len(data) > 0 {
		assert.Nil(err, "%d: %+v", idx, err)
		assert.EqualValues(data, res.Data)
	} else {
		assert.NotNil(err, "%d", idx)
		// check error code!
		assert.True(check(err), "%d: %+v", idx, err)
	}
}
