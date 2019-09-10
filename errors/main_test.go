package errors

import (
	stderr "errors"
	"strconv"
	"testing"

	pkerr "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	wrsp "github.com/tepleton/wrsp/types"
)

func TestCreateResult(t *testing.T) {
	assert := assert.New(t)

	cases := []struct {
		err  error
		msg  string
		code wrsp.CodeType
	}{
		{stderr.New("base"), "base", defaultErrCode},
		{pkerr.New("dave"), "dave", defaultErrCode},
		{New("nonce", wrsp.CodeType_BadNonce), "nonce", wrsp.CodeType_BadNonce},
		{Wrap(stderr.New("wrap")), "wrap", defaultErrCode},
		{WithCode(stderr.New("coded"), wrsp.CodeType_BaseInvalidInput), "coded", wrsp.CodeType_BaseInvalidInput},
		{ErrDecoding(), errDecoding.Error(), wrsp.CodeType_EncodingError},
		{ErrUnauthorized(), errUnauthorized.Error(), wrsp.CodeType_Unauthorized},
	}

	for idx, tc := range cases {
		i := strconv.Itoa(idx)

		res := Result(tc.err)
		assert.True(res.IsErr(), i)
		assert.Equal(tc.msg, res.Log, i)
		assert.Equal(tc.code, res.Code, i)
	}
}
