package app

import (
	"encoding/json"

	wrsp "github.com/tepleton/wrsp/types"
	oldwire "github.com/tepleton/go-wire"
	cmn "github.com/tepleton/tmlibs/common"
	dbm "github.com/tepleton/tmlibs/db"
	"github.com/tepleton/tmlibs/log"

	bam "github.com/tepleton/tepleton-sdk/baseapp"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/tepleton/tepleton-sdk/x/bank"
	"github.com/tepleton/tepleton-sdk/x/ibc"
	"github.com/tepleton/tepleton-sdk/x/stake"
)

const (
	appName = "GaiaApp"
)

// Extended WRSP application
type GaiaApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	capKeyMainStore    *sdk.KVStoreKey
	capKeyAccountStore *sdk.KVStoreKey
	capKeyIBCStore     *sdk.KVStoreKey
	capKeyStakeStore   *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper sdk.AccountMapper
}

func NewGaiaApp(logger log.Logger, dbs map[string]dbm.DB) *GaiaApp {
	// create your application object
	var app = &GaiaApp{
		BaseApp:            bam.NewBaseApp(appName, logger, dbs["main"]),
		cdc:                MakeCodec(),
		capKeyMainStore:    sdk.NewKVStoreKey("main"),
		capKeyAccountStore: sdk.NewKVStoreKey("acc"),
		capKeyIBCStore:     sdk.NewKVStoreKey("ibc"),
		capKeyStakeStore:   sdk.NewKVStoreKey("stake"),
	}

	// define the accountMapper
	app.accountMapper = auth.NewAccountMapperSealed(
		app.capKeyMainStore, // target store
		&auth.BaseAccount{}, // prototype
	)

	// add handlers
	coinKeeper := bank.NewCoinKeeper(app.accountMapper)
	ibcMapper := ibc.NewIBCMapper(app.cdc, app.capKeyIBCStore)
	stakeKeeper := stake.NewKeeper(app.cdc, app.capKeyStakeStore, coinKeeper)
	app.Router().
		AddRoute("bank", bank.NewHandler(coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(ibcMapper, coinKeeper)).
		AddRoute("stake", stake.NewHandler(stakeKeeper))

	// initialize BaseApp
	app.SetTxDecoder(app.txDecoder)
	app.SetInitChainer(app.initChainer)
	app.SetEndBlocker(stake.NewEndBlocker(stakeKeeper))
	app.MountStoreWithDB(app.capKeyMainStore, sdk.StoreTypeIAVL, dbs["main"])
	app.MountStoreWithDB(app.capKeyAccountStore, sdk.StoreTypeIAVL, dbs["acc"])
	app.MountStoreWithDB(app.capKeyIBCStore, sdk.StoreTypeIAVL, dbs["ibc"])
	app.MountStoreWithDB(app.capKeyStakeStore, sdk.StoreTypeIAVL, dbs["stake"])

	// NOTE: Broken until #532 lands
	//app.MountStoresIAVL(app.capKeyMainStore, app.capKeyIBCStore, app.capKeyStakingStore)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper))
	err := app.LoadLatestVersion(app.capKeyMainStore)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// custom tx codec
// TODO: use new go-wire
func MakeCodec() *wire.Codec {
	const (
		msgTypeSend           = 0x1
		msgTypeIssue          = 0x2
		msgTypeIBCTransferMsg = 0x3
		msgTypeIBCReceiveMsg  = 0x4
		msgDeclareCandidacy   = 0x5
		msgEditCandidacy      = 0x6
		msgDelegate           = 0x7
		msgUnbond             = 0x8
	)
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Msg }{},
		oldwire.ConcreteType{bank.SendMsg{}, msgTypeSend},
		oldwire.ConcreteType{bank.IssueMsg{}, msgTypeIssue},
		oldwire.ConcreteType{ibc.IBCTransferMsg{}, msgTypeIBCTransferMsg},
		oldwire.ConcreteType{ibc.IBCReceiveMsg{}, msgTypeIBCReceiveMsg},
		oldwire.ConcreteType{stake.MsgDeclareCandidacy{}, msgDeclareCandidacy},
		oldwire.ConcreteType{stake.MsgEditCandidacy{}, msgEditCandidacy},
		oldwire.ConcreteType{stake.MsgDelegate{}, msgDelegate},
		oldwire.ConcreteType{stake.MsgUnbond{}, msgUnbond},
	)

	const accTypeApp = 0x1
	var _ = oldwire.RegisterInterface(
		struct{ sdk.Account }{},
		oldwire.ConcreteType{&auth.BaseAccount{}, accTypeApp},
	)
	cdc := wire.NewCodec()

	// cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	// bank.RegisterWire(cdc)   // Register bank.[SendMsg,IssueMsg]
	// crypto.RegisterWire(cdc) // Register crypto.[PubKey,PrivKey,Signature]
	// ibc.RegisterWire(cdc) // Register ibc.[IBCTransferMsg, IBCReceiveMsg]
	return cdc
}

// custom logic for transaction decoding
func (app *GaiaApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx{}

	if len(txBytes) == 0 {
		return nil, sdk.ErrTxDecode("txBytes are empty")
	}

	// StdTx.Msg is an interface. The concrete types
	// are registered by MakeTxCodec in bank.RegisterWire.
	err := app.cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("").TraceCause(err, "")
	}
	return tx, nil
}

// custom logic for ton initialization
func (app *GaiaApp) initChainer(ctx sdk.Context, req wrsp.RequestInitChain) wrsp.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := json.Unmarshal(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/tepleton/tepleton-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		app.accountMapper.SetAccount(ctx, &acc)
	}
	return wrsp.ResponseInitChain{}
}
