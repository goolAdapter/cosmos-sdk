package lottery

import (
	"github.com/cosmos/cosmos-sdk/examples/lottery/voter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var (
	LotteryPrefixKey = []byte{0x01}
	sequenceKey      = []byte("sequence")
	statusKey        = []byte("status")
	prevoteKey       = []byte("prevote")
	voteKey          = []byte("vote")
)

const (
	WaitforNewRoundPhase int64 = iota
	WaitforPreVotePhase
	WaitforVotePhase
	GenerateResultPhase
)

// LotteryKeeper - handlers sets/gets of custom variables for your module
type LotteryKeeper struct {
	// The (unexposed) key used to access the fee store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec

	voterKeeper voter.VoterKeeper
}

// NewLotteryKeeper returns a new LotteryKeeper
func NewLotteryKeeper(cdc *wire.Codec, key sdk.StoreKey, vc voter.VoterKeeper) LotteryKeeper {
	return LotteryKeeper{
		key:         key,
		cdc:         cdc,
		voterKeeper: vc,
	}
}

func (lk LotteryKeeper) SetupLottery(ctx sdk.Context, msg MsgSetupLottery) sdk.Error {
	seq := lk.GetSequence(ctx, msg.Address)
	if seq != msg.Sequence {
		return ErrSequenceNotMatch(DefaultCodespace, seq, msg.Sequence)
	}

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

func (lk LotteryKeeper) CheckForStartRound(ctx sdk.Context, msg MsgStartLotteryRound) (status int64, seq int64, err sdk.Error) {
	status = lk.GetStatus(ctx, msg.Address)
	if status != WaitforNewRoundPhase {
		return 0, 0, ErrStatusNotMatch(DefaultCodespace, status)
	}

	seq = lk.GetSequence(ctx, msg.Address)
	if seq+1 != msg.Sequence {
		return 0, 0, ErrSequenceNotMatch(DefaultCodespace, seq, msg.Sequence)
	}

	return status, seq, nil
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

func (lk LotteryKeeper) SetSequence(ctx sdk.Context, address sdk.AccAddress, seq int64) error {
	store := ctx.KVStore(lk.key)
	key := GetInfoSequenceKey(address, lk.cdc)

	bz := marshalBinaryPanic(lk.cdc, seq)
	store.Set(key, bz)
	return nil
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

func (lk LotteryKeeper) SetStatus(ctx sdk.Context, address sdk.AccAddress, status int64) error {
	store := ctx.KVStore(lk.key)
	key := GetInfoStatusKey(address, lk.cdc)

	bz := marshalBinaryPanic(lk.cdc, status)
	store.Set(key, bz)
	return nil
}

func (lk LotteryKeeper) StartLotteryRound(ctx sdk.Context, msg MsgStartLotteryRound) sdk.Error {
	_, seq, err := lk.CheckForStartRound(ctx, msg)
	if err != nil {
		return err
	}

	lk.SetSequence(ctx, msg.Address, seq+1)
	lk.SetStatus(ctx, msg.Address, int64(WaitforPreVotePhase))

	// iterate to get the voters
	voters := []voter.Voter{}
	appendVoter := func(ac voter.Voter) (stop bool) {
		voters = append(voters, ac)
		return false
	}
	lk.voterKeeper.IterateVoters(ctx, appendVoter)

	var preVotes PreVoteItemList
	for _, v := range voters {
		preVotes = append(preVotes, PreVoteItem{v.Address, []byte("")})
	}
	preVotes = preVotes.Sort()
	store := ctx.KVStore(lk.key)
	key := GetInfoPrevoteKey(msg.Address, lk.cdc)
	store.Set(key, lk.cdc.MustMarshalBinary(preVotes))

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
	return append(LotteryPrefixKey, address...)
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

func GetInfoPrevoteKey(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return append(GetInfoPrefix(address, cdc), prevoteKey...)
}

func GetInfoVoteKey(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return append(GetInfoPrefix(address, cdc), voteKey...)
}
