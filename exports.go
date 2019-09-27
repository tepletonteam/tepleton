package sdk

import (
	"github.com/tepleton/tepleton-sdk/store"
	types "github.com/tepleton/tepleton-sdk/types"
)

type (
	// Type aliases for the tepleton-sdk/types module.  We keep all of them in
	// types/* but they are all meant to be imported as
	// "github.com/tepleton/tepleton-sdk".  So, add all of them.
	Handler   = types.Handler
	Context   = types.Context
	Decorator = types.Decorator

	// Type aliases for other modules.
	MultiStore = store.MultiStore
)
