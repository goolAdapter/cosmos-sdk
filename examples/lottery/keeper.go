package lottery

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var (
	sequenceKey = []byte("sequence")
	statusKey   = []byte("status")
)

const (
	waitforNewRoundPhase int64 = iota
	waitforPreVotePhase
	waitforVotePhase
	generateResultPhase
)

// LotteryKeeper - handlers sets/gets of custom variables for your module
type LotteryKeeper struct {
	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewLotteryKeeper returns a new LotteryKeeper
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

func marshalBinaryPanic(cdc *wire.Codec, value interface{}) []byte {
	res, err := cdc.MarshalBinary(value)
	if err != nil {
		panic(err)
	}
	return res
}

func unmarshalBinaryPanic(cdc *wire.Codec, bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinary(bz, ptr)
	if err != nil {
		panic(err)
	}
}

func (lk LotteryKeeper) CheckForStartRound(ctx sdk.Context, msg MsgStartLotteryRound) sdk.Error {
	status := lk.GetStatus(ctx, msg.Address)
	if status != waitforNewRoundPhase {
		return ErrStatusNotMatch(DefaultCodespace, status)
	}

	seq := lk.GetSequence(ctx, msg.Address)
	if seq+1 != msg.Sequence {
		return ErrSequenceNotMatch(DefaultCodespace, seq, msg.Sequence)
	}

	return nil
}

func (lk LotteryKeeper) GetSequence(ctx sdk.Context, address sdk.AccAddress) int64 {
	store := ctx.KVStore(lk.key)
	key := GetInfoSequenceKey(address, lk.cdc)

	bz := store.Get(key)
	if bz == nil {
		return 0
	}

	var res int64
	unmarshalBinaryPanic(lk.cdc, bz, &res)
	return res
}

func (lk LotteryKeeper) GetStatus(ctx sdk.Context, address sdk.AccAddress) int64 {
	store := ctx.KVStore(lk.key)
	key := GetInfoStatusKey(address, lk.cdc)

	bz := store.Get(key)
	if bz == nil {
		return 0
	}

	var res int64
	unmarshalBinaryPanic(lk.cdc, bz, &res)
	return res
}

func (lk LotteryKeeper) StartLotteryRound(ctx sdk.Context, msg MsgStartLotteryRound) sdk.Error {
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

func GetInfoSequenceKey(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return append(GetInfoPrefix(address, cdc), sequenceKey...)
}

func GetInfoStatusKey(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return append(GetInfoPrefix(address, cdc), statusKey...)
}
