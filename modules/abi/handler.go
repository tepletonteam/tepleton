package abi

import (
	"github.com/tepleton/go-wire/data"
	"github.com/tepleton/tmlibs/log"

	"github.com/tepleton/basecoin"
	"github.com/tepleton/basecoin/errors"
	"github.com/tepleton/basecoin/stack"
	"github.com/tepleton/basecoin/state"
)

const (
	// NameABI is the name of this module
	NameABI = "abi"
	// OptionRegistrar is the option name to set the actor
	// to handle abi chain registration
	OptionRegistrar = "registrar"
)

var (
	allowABI = []byte{0x42, 0xbe, 0xef, 0x1}
)

// AllowABI is the special code that an app must set to
// enable sending ABI packets for this app-type
func AllowABI(app string) basecoin.Actor {
	return basecoin.Actor{App: app, Address: allowABI}
}

// Handler allows us to update the chain state or create a packet
type Handler struct {
	Registrar basecoin.Actor
}

var _ basecoin.Handler = &Handler{}

// NewHandler makes a Handler that allows all chains to connect via ABI.
// Set a Registrar via SetOption to restrict it.
func NewHandler() *Handler {
	return new(Handler)
}

// Name - return name space
func (*Handler) Name() string {
	return NameABI
}

// SetOption - sets the registrar for ABI
func (h *Handler) SetOption(l log.Logger, store state.KVStore, module, key, value string) (log string, err error) {
	if module != NameABI {
		return "", errors.ErrUnknownModule(module)
	}
	if key == OptionRegistrar {
		var act basecoin.Actor
		err = data.FromJSON([]byte(value), &act)
		if err != nil {
			return "", err
		}
		h.Registrar = act
		// TODO: save/load from disk!
		return "Success", nil
	}
	return "", errors.ErrUnknownKey(key)
}

// CheckTx verifies the packet is formated correctly, and has the proper sequence
// for a registered chain
func (h *Handler) CheckTx(ctx basecoin.Context, store state.KVStore, tx basecoin.Tx) (res basecoin.Result, err error) {
	err = tx.ValidateBasic()
	if err != nil {
		return res, err
	}

	switch t := tx.Unwrap().(type) {
	case RegisterChainTx:
		return h.initSeed(ctx, store, t)
	case UpdateChainTx:
		return h.updateSeed(ctx, store, t)
	case CreatePacketTx:
		return h.createPacket(ctx, store, t)
	}
	return res, errors.ErrUnknownTxType(tx.Unwrap())
}

// DeliverTx verifies all signatures on the tx and updated the chain state
// apropriately
func (h *Handler) DeliverTx(ctx basecoin.Context, store state.KVStore, tx basecoin.Tx) (res basecoin.Result, err error) {
	err = tx.ValidateBasic()
	if err != nil {
		return res, err
	}

	switch t := tx.Unwrap().(type) {
	case RegisterChainTx:
		return h.initSeed(ctx, store, t)
	case UpdateChainTx:
		return h.updateSeed(ctx, store, t)
	case CreatePacketTx:
		return h.createPacket(ctx, store, t)
	}
	return res, errors.ErrUnknownTxType(tx.Unwrap())
}

// initSeed imports the first seed for this chain and
// accepts it as the root of trust.
//
// only the registrar, if set, is allowed to do this
func (h *Handler) initSeed(ctx basecoin.Context, store state.KVStore,
	t RegisterChainTx) (res basecoin.Result, err error) {

	// check permission to attach
	// nothing set, means anyone can connect
	if !h.Registrar.Empty() && !ctx.HasPermission(h.Registrar) {
		return res, errors.ErrUnauthorized()
	}

	chainID := t.ChainID()
	s := NewChainSet(store)
	err = s.Register(chainID, ctx.BlockHeight(), t.Seed.Height())
	if err != nil {
		return res, err
	}

	space := stack.PrefixedStore(chainID, store)
	provider := newDBProvider(space)
	err = provider.StoreSeed(t.Seed)
	return res, err
}

// updateSeed checks the seed against the existing chain data and rejects it if it
// doesn't fit (or no chain data)
func (h *Handler) updateSeed(ctx basecoin.Context, store state.KVStore,
	t UpdateChainTx) (res basecoin.Result, err error) {

	chainID := t.ChainID()
	if !NewChainSet(store).Exists([]byte(chainID)) {
		return res, ErrNotRegistered(chainID)
	}

	// load the certifier for this chain
	seed := t.Seed
	space := stack.PrefixedStore(chainID, store)
	cert, err := newCertifier(space, chainID, seed.Height())
	if err != nil {
		return res, err
	}

	// this will import the seed if it is valid in the current context
	err = cert.Update(seed.Checkpoint, seed.Validators)
	return res, err
}

// createPacket makes sure all permissions are good and the destination
// chain is registed.  If so, it appends it to the outgoing queue
func (h *Handler) createPacket(ctx basecoin.Context, store state.KVStore,
	t CreatePacketTx) (res basecoin.Result, err error) {

	// make sure the chain is registed
	dest := t.DestChain
	if !NewChainSet(store).Exists([]byte(dest)) {
		return res, ErrNotRegistered(dest)
	}

	// make sure we have the special ABI permission
	mod, err := t.Tx.GetMod()
	if err != nil {
		return res, err
	}
	if !ctx.HasPermission(AllowABI(mod)) {
		return res, ErrNeedsABIPermission()
	}

	// start making the packet to send
	packet := Packet{
		DestChain:   t.DestChain,
		Tx:          t.Tx,
		Permissions: make([]basecoin.Actor, len(t.Permissions)),
	}

	// make sure we have all the permissions we want to send
	for i, p := range t.Permissions {
		if !ctx.HasPermission(p) {
			return res, ErrCannotSetPermission()
		}
		// add the permission with the current ChainID
		packet.Permissions[i] = p.WithChain(ctx.ChainID())
	}

	// now add it to the output queue....
	// TODO: where to store, also set the sequence....
	return res, nil
}