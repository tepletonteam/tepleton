package app

import (
	"os"

	wrsp "github.com/tepleton/wrsp/types"
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

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.toncli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.tond")
)

// Extended WRSP application
type GaiaApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	keyMain    *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyIBC     *sdk.KVStoreKey
	keyStake   *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper sdk.AccountMapper
	coinKeeper    bank.Keeper
	ibcMapper     ibc.Mapper
	stakeKeeper   stake.Keeper
}

func NewGaiaApp(logger log.Logger, db dbm.DB) *GaiaApp {
	cdc := MakeCodec()

	// create your application object
	var app = &GaiaApp{
		BaseApp:    bam.NewBaseApp(appName, cdc, logger, db),
		cdc:        cdc,
		keyMain:    sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("acc"),
		keyIBC:     sdk.NewKVStoreKey("ibc"),
		keyStake:   sdk.NewKVStoreKey("stake"),
	}

	// add accountMapper/handlers
	app.accountMapper = auth.NewAccountMapper(app.cdc, app.keyMain, &auth.BaseAccount{})
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))

	// register message routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		AddRoute("stake", stake.NewHandler(app.stakeKeeper))

	// initialize BaseApp
	app.SetInitChainer(app.initChainer)
	app.SetEndBlocker(stake.NewEndBlocker(app.stakeKeeper))
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keyStake)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, stake.FeeHandler))
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	ibc.RegisterWire(cdc)
	bank.RegisterWire(cdc)
	stake.RegisterWire(cdc)
	auth.RegisterWire(cdc)
	sdk.RegisterWire(cdc)
	wire.RegisterCrypto(cdc)
	return cdc
}

// custom logic for ton initialization
func (app *GaiaApp) initChainer(ctx sdk.Context, req wrsp.RequestInitChain) wrsp.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/tepleton/tepleton-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		app.accountMapper.SetAccount(ctx, acc)
	}

	// load the initial stake information
	stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)

	return wrsp.ResponseInitChain{}
}

// custom logic for export
func (app *GaiaApp) ExportGenesis() GenesisState {
	return GenesisState{
		Accounts:  []GenesisAccount{},
		StakeData: stake.WriteGenesis(app.NewContext(true, wrsp.Header{}), app.stakeKeeper),
	}
}
