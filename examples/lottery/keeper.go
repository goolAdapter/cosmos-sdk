package lottery

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var (
	collectedFeesKey = []byte("collectedFees")
)

// LotteryKeeper - handlers sets/gets of custom variables for your module
type LotteryKeeper struct {
	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewFeeKeeper returns a new FeeKeeper
func NewLotteryKeeper(cdc *wire.Codec, key sdk.StoreKey) LotteryKeeper {
	return LotteryKeeper{
		key: key,
		cdc: cdc,
	}
}

func (lk LotteryKeeper) SetupLottery(ctx sdk.Context, msg MsgSetupLottery) sdk.Error {
	if len(msg.ParticipateNames) == 0 {
		return ErrEmptyActor(DefaultCodespace)
	}

	lk.clearInfo(ctx, msg.Address)
	store := ctx.KVStore(lk.key)

	key := GetInfoSetupKey(msg.Address, lk.cdc)
	bz := lk.cdc.MustMarshalBinary(msg)
	store.Set(key, bz)
	return nil
}

func (lk LotteryKeeper) clearInfo(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(lk.key)

	prefix := GetInfoPrefix(address, lk.cdc)

	iter := sdk.KVStorePrefixIterator(store, prefix)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
	iter.Close()
}

func GetInfoPrefix(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return append([]byte{0x00}, address...)
}

func GetInfoSetupKey(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return GetInfoPrefix(address, cdc)
}
