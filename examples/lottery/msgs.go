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
