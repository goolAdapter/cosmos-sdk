package lottery

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var maxMemoCharacters = 1024
var maxRoundNumber int64 = 10240

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
			fmt.Sprintf("maximum number of characters is %d but received %d characters",
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
	Number   int64          `json:"number"` //how many prize number to generate on this round.
	Memo     string         `json:"memo"`
}

func NewMsgStartLotteryRound(addr sdk.AccAddress, seq int64, number int64, memo string) MsgStartLotteryRound {
	return MsgStartLotteryRound{
		Address:  addr,
		Sequence: seq,
		Number:   number,
		Memo:     memo,
	}
}

func (msg MsgStartLotteryRound) Type() string                 { return MsgTypeName }
func (msg MsgStartLotteryRound) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgStartLotteryRound) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgStartLotteryRound.Address must not be empty")
	}

	if msg.Number <= 0 || msg.Number > maxRoundNumber {
		return sdk.ErrInvalidPubKey(fmt.Sprintf("MsgStartLotteryRound.Number %d is invalid", msg.Number))
	}

	if len(msg.Memo) > maxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf("maximum number of characters is %d but received %d characters",
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
	Address sdk.AccAddress `json:"address"`
	Hash    []byte         `json:"hash"`
}

func NewMsgPreVote(addr sdk.AccAddress, hash []byte) MsgPreVote {
	return MsgPreVote{
		Address: addr,
		Hash:    hash,
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
	Address sdk.AccAddress `json:"address"`
	Value   int64          `json:"value"`
}

func NewMsgVote(addr sdk.AccAddress, val int64) MsgVote {
	return MsgVote{
		Address: addr,
		Value:   val,
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
