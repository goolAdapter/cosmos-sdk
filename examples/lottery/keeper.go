package lottery

import (
	"bytes"
	"encoding/binary"
	"math"

	"github.com/cosmos/cosmos-sdk/examples/lottery/voter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"golang.org/x/crypto/ripemd160"
)

var (
	LotteryPrefixKey = []byte{0x01}
	sequenceKey      = []byte("sequence")
	statusKey        = []byte("status")
	voteKey          = []byte("vote")
	resultKey        = []byte("result")
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

func (lk LotteryKeeper) checkForStartRound(ctx sdk.Context, msg MsgStartLotteryRound) (status int64, seq int64, err sdk.Error) {
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

func (lk LotteryKeeper) SetResult(ctx sdk.Context, address sdk.AccAddress, seq, result int64) error {
	store := ctx.KVStore(lk.key)
	key := GetInfoResultKey(address, seq, lk.cdc)

	bz := marshalBinaryPanic(lk.cdc, result)
	store.Set(key, bz)
	return nil
}

func (lk LotteryKeeper) StartLotteryRound(ctx sdk.Context, msg MsgStartLotteryRound) sdk.Error {
	_, seq, err := lk.checkForStartRound(ctx, msg)
	if err != nil {
		return err
	}

	// iterate to get the voters
	voters := []voter.Voter{}
	appendVoter := func(ac voter.Voter) (stop bool) {
		voters = append(voters, ac)
		return false
	}
	lk.voterKeeper.IterateVoters(ctx, appendVoter)

	var vl VoteItemList
	for _, v := range voters {
		vl = append(vl, VoteItem{v.Address, []byte(""), 0})
	}
	vl = vl.Sort()
	store := ctx.KVStore(lk.key)
	key := GetInfoVoteKey(msg.Address, lk.cdc)
	store.Set(key, lk.cdc.MustMarshalBinary(vl))

	lk.SetSequence(ctx, msg.Address, seq+1)
	lk.SetStatus(ctx, msg.Address, int64(WaitforPreVotePhase))

	return nil
}

func (lk LotteryKeeper) checkFoPreVoteRound(ctx sdk.Context, msg MsgPreVote) (status int64, seq int64, err sdk.Error) {
	status = lk.GetStatus(ctx, msg.Target)
	if status != WaitforPreVotePhase && status != WaitforVotePhase {
		return 0, 0, ErrStatusNotMatch(DefaultCodespace, status)
	}

	seq = lk.GetSequence(ctx, msg.Target)
	if seq != msg.Sequence {
		return 0, 0, ErrSequenceNotMatch(DefaultCodespace, seq, msg.Sequence)
	}

	return status, seq, nil
}

func (lk LotteryKeeper) HandlePreVoteMsg(ctx sdk.Context, msg MsgPreVote) sdk.Error {
	status, _, err := lk.checkFoPreVoteRound(ctx, msg)
	if err != nil {
		return err
	}

	store := ctx.KVStore(lk.key)
	key := GetInfoVoteKey(msg.Target, lk.cdc)
	bz := store.Get(key)
	if bz == nil {
		return err
	}

	var hashedCount = 0
	var votes VoteItemList
	unmarshalBinaryPanic(lk.cdc, bz, &votes)
	for i := 0; i < len(votes); i++ {
		if len(votes[i].Hash) == 0 {
			if bytes.Equal(votes[i].Address, msg.Address) {
				votes[i].Hash = msg.Hash
				hashedCount++
			}
		} else {
			hashedCount++
		}
	}

	store.Set(key, lk.cdc.MustMarshalBinary(votes))

	if status == WaitforPreVotePhase {
		if hashedCount*5 > len(votes)*4 {
			lk.SetStatus(ctx, msg.Target, int64(WaitforVotePhase))
		}
	}
	return nil
}

func (lk LotteryKeeper) checkForVoteRound(ctx sdk.Context, msg MsgVote) (status int64, seq int64, err sdk.Error) {
	status = lk.GetStatus(ctx, msg.Target)
	if status != WaitforVotePhase {
		return 0, 0, ErrStatusNotMatch(DefaultCodespace, status)
	}

	seq = lk.GetSequence(ctx, msg.Target)
	if seq != msg.Sequence {
		return 0, 0, ErrSequenceNotMatch(DefaultCodespace, seq, msg.Sequence)
	}

	return status, seq, nil
}

func (lk LotteryKeeper) voteHash(v int64) []byte {
	// Doesn't write Name, since merkle.SimpleHashFromMap() will
	// include them via the keys.
	bz, _ := lk.cdc.MarshalBinary(v) // Does not error
	hasher := ripemd160.New()
	_, err := hasher.Write(bz)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}
	return hasher.Sum(nil)
}

func (lk LotteryKeeper) HandleVoteMsg(ctx sdk.Context, msg MsgVote) sdk.Error {
	status, _, err := lk.checkForVoteRound(ctx, msg)
	if err != nil {
		return err
	}

	if msg.Value > int64(math.Pow(2, 32)) {
		return ErrInvalidValue(DefaultCodespace)
	}

	store := ctx.KVStore(lk.key)
	key := GetInfoVoteKey(msg.Target, lk.cdc)
	bz := store.Get(key)
	if bz == nil {
		return err
	}

	var resultValue int64
	var hashedCount = 0
	var votes VoteItemList
	unmarshalBinaryPanic(lk.cdc, bz, &votes)
	for i := 0; i < len(votes); i++ {
		if votes[i].Value == 0 {
			if bytes.Equal(votes[i].Address, msg.Address) && bytes.Equal(lk.voteHash(msg.Value), votes[i].Hash) {
				votes[i].Value = msg.Value
				resultValue += votes[i].Value
				hashedCount++
			}
		} else {
			resultValue += votes[i].Value
			hashedCount++
		}
	}

	store.Set(key, lk.cdc.MustMarshalBinary(votes))
	if status == WaitforVotePhase {
		if hashedCount*4 > len(votes)*3 {
			lk.SetResult(ctx, msg.Target, msg.Sequence, int64(resultValue))
			lk.SetStatus(ctx, msg.Target, int64(GenerateResultPhase))
		}
	}
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

func GetInfoVoteKey(address sdk.AccAddress, cdc *wire.Codec) []byte {
	return append(GetInfoPrefix(address, cdc), voteKey...)
}

func GetInfoResultKey(address sdk.AccAddress, seq int64, cdc *wire.Codec) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(seq))
	return append(append(GetInfoPrefix(address, cdc), resultKey...), b...)
}
