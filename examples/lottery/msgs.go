package lottery

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MsgTypeName string = "lottery"

//MsgSetupLottery reset a new lottery status
var _ sdk.Msg = MsgSetupLottery{}

type MsgSetupLottery struct {
	Address          sdk.AccAddress `json:"address"`
	ParticipateNames []string       `json:"participate_names"`
}

func NewMsgSetupLottery(addr sdk.AccAddress, names []string) MsgSetupLottery {
	return MsgSetupLottery{
		Address:          addr,
		ParticipateNames: names,
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
	Sequence int64          `json:"seq"`
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

	if msg.Number <= 0 {
		return sdk.ErrInvalidPubKey("MsgStartLotteryRound.Number must not be zero")
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
