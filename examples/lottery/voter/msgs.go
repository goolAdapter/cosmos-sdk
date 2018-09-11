package voter

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MsgTypeName string = "voter"

//MsgCreatevoter create a voter
var _ sdk.Msg = MsgCreateVoter{}

type MsgCreateVoter struct {
	Address sdk.AccAddress `json:"address"`
	Memo    string         `json:"memo"`
}

func NewMsgCreateVoter(addr sdk.AccAddress, memo string) MsgCreateVoter {
	return MsgCreateVoter{
		Address: addr,
		Memo:    memo,
	}
}

func (msg MsgCreateVoter) Type() string                 { return MsgTypeName }
func (msg MsgCreateVoter) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgCreateVoter) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgCreatevoter.Address must not be empty")
	}

	return nil
}

func (msg MsgCreateVoter) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(bz)
}
