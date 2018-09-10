package app

import (
	"encoding/json"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/examples/lottery/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "LotteryApp"
)

// LotteryApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore keys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type LotteryApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the multistore
	keyMain    *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey

	// manage getting and setting accounts
	accountMapper auth.AccountMapper
}

// NewLotteryApp returns a reference to a new LotteryApp given a logger and
// database. Internally, a codec is created along with all the necessary keys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewLotteryApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *LotteryApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()

	// create your application type
	var app = &LotteryApp{
		cdc:        cdc,
		BaseApp:    bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...),
		keyMain:    sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("acc"),
	}

	// define and attach the mappers and keepers
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyAccount, // target store
		func() auth.Account {
			return &types.AppAccount{}
		},
	)
	// register message routes
	app.Router().
		AddRoute("lottery", types.NewHandler())

	// perform initialization logic
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// mount the multistore and load the latest state
	app.MountStoresIAVL(app.keyMain, app.keyAccount)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	app.Seal()

	return app
}

// MakeCodec creates a new wire codec and registers all the necessary types
// with the codec.
func MakeCodec() *wire.Codec {
	cdc := wire.NewCodec()

	wire.RegisterCrypto(cdc)
	sdk.RegisterWire(cdc)
	auth.RegisterWire(cdc)

	// register custom type
	cdc.RegisterConcrete(&types.AppAccount{}, "lottery/Account", nil)

	cdc.Seal()

	return cdc
}

// BeginBlocker reflects logic to run before any TXs application are processed
// by the application.
func (app *LotteryApp) BeginBlocker(_ sdk.Context, _ abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

// EndBlocker reflects logic to run after all TXs are processed by the
// application.
func (app *LotteryApp) EndBlocker(_ sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

// initChainer implements the custom application logic that the BaseApp will
// invoke upon initialization. In this case, it will take the application's
// state provided by 'req' and attempt to deserialize said state. The state
// should contain all the genesis accounts. These accounts will be added to the
// application's account mapper.
func (app *LotteryApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
		panic(err)
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			// TODO: https://github.com/cosmos/cosmos-sdk/issues/468
			panic(err)
		}

		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	return abci.ResponseInitChain{}
}

// ExportAppStateAndValidators implements custom application logic that exposes
// various parts of the application's state and set of validators. An error is
// returned if any step getting the state or set of validators fails.
func (app *LotteryApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*types.GenesisAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountMapper.IterateAccounts(ctx, appendAccountsFn)

	genState := types.GenesisState{Accounts: accounts}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}
