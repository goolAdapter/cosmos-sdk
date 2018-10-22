package lottery

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var maxMemoCharacters = 1024
var maxRoundAmount int64 = 10240

const MsgTypeName string = "lottery"

//MsgSetupLottery reset a new lottery status
var _ sdk.Msg = MsgSetupLottery{}

type MsgSetupLottery struct {
	Address          sdk.AccAddress `json:"address"`
	Sequence         int64          `json:"sequence"`
	ParticipateNames []string       `json:"participate_names"`
	Memo             string         `json:"memo"`
}

func NewMsgSetupLottery(addr sdk.AccAddress, seq int64, names []string, memo string) MsgSetupLottery {
	return MsgSetupLottery{
		Address:          addr,
		Sequence:         seq,
		ParticipateNames: names,
		Memo:             memo,
	}
}

func (msg MsgSetupLottery) Type() string                 { return MsgTypeName }
func (msg MsgSetupLottery) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgSetupLottery) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgSetupLottery.Address must not be empty")
	}

	if msg.ParticipateNames == nil {
		return sdk.ErrInvalidPubKey("MsgSetupLottery.ParticipateNames must not be empty")
	}

	if len(msg.Memo) > maxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf("maximum length of characters is %d but received %d characters",
				maxMemoCharacters, len(msg.Memo)))
	}

	return nil
}

func (msg MsgSetupLottery) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(bz)
}

//MsgStartLotteryRound start a lottery round
var _ sdk.Msg = MsgStartLotteryRound{}

type MsgStartLotteryRound struct {
	Address  sdk.AccAddress `json:"address"`
	Sequence int64          `json:"sequence"`
	Amount   int64          `json:"amount"` //how many prize number to generate on this round.
	Memo     string         `json:"memo"`
}

func NewMsgStartLotteryRound(addr sdk.AccAddress, seq int64, amount int64, memo string) MsgStartLotteryRound {
	return MsgStartLotteryRound{
		Address:  addr,
		Sequence: seq,
		Amount:   amount,
		Memo:     memo,
	}
}

func (msg MsgStartLotteryRound) Type() string                 { return MsgTypeName }
func (msg MsgStartLotteryRound) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgStartLotteryRound) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgStartLotteryRound.Address must not be empty")
	}

	if msg.Amount <= 0 || msg.Amount > maxRoundAmount {
		return sdk.ErrInvalidPubKey(fmt.Sprintf("MsgStartLotteryRound.Amount %d is invalid", msg.Amount))
	}

	if len(msg.Memo) > maxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf("maximum length of characters is %d but received %d characters",
				maxMemoCharacters, len(msg.Memo)))
	}

	return nil
}

func (msg MsgStartLotteryRound) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(bz)
}

//MsgPreVote
var _ sdk.Msg = MsgPreVote{}

type MsgPreVote struct {
	Target   sdk.AccAddress `json:"address"`
	Address  sdk.AccAddress `json:"address"`
	Sequence int64          `json:"sequence"`
	Hash     []byte         `json:"hash"`
}

func NewMsgPreVote(target, from sdk.AccAddress, seq int64, hash []byte) MsgPreVote {
	return MsgPreVote{
		Target:   target,
		Address:  from,
		Sequence: seq,
		Hash:     hash,
	}
}

func (msg MsgPreVote) Type() string                 { return MsgTypeName }
func (msg MsgPreVote) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgPreVote) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgPreVote.Address must not be empty")
	}

	if len(msg.Hash) == 0 {
		return ErrInvalidValue(DefaultCodespace)
	}

	return nil
}

func (msg MsgPreVote) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(bz)
}

//MsgVote
var _ sdk.Msg = MsgVote{}

type MsgVote struct {
	Target   sdk.AccAddress `json:"address"`
	Address  sdk.AccAddress `json:"address"`
	Sequence int64          `json:"sequence"`
	Value    int64          `json:"value"`
}

func NewMsgVote(target, from sdk.AccAddress, seq int64, val int64) MsgVote {
	return MsgVote{
		Target:   target,
		Address:  from,
		Sequence: seq,
		Value:    val,
	}
}

func (msg MsgVote) Type() string                 { return MsgTypeName }
func (msg MsgVote) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgVote) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgVote.Address must not be empty")
	}

	return nil
}

func (msg MsgVote) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(bz)
}
