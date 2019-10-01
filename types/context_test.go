package types_test

import (
	"testing"

	"github.com/tepleton/tepleton-sdk/types"
	wrsp "github.com/tepleton/wrsp/types"
)

func TestContextGetOpShouldNeverPanic(t *testing.T) {
	var ms types.MultiStore
	ctx := types.NewContext(ms, wrsp.Header{}, false, nil)
	indices := []int64{
		-10, 1, 0, 10, 20,
	}

	for _, index := range indices {
		_, _ = ctx.GetOp(index)
	}
}
