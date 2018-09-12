package voter

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var maxMemoCharacters = 1024

const MsgTypeName string = "voter"

//MsgCreateVoter create a voter
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

	if len(msg.Memo) > maxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf("maximum number of characters is %d but received %d characters",
				maxMemoCharacters, len(msg.Memo)))
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

//MsgRevocationVoter revocation a voter
var _ sdk.Msg = MsgRevocationVoter{}

type MsgRevocationVoter struct {
	Address sdk.AccAddress `json:"address"`
}

func NewMsgRevocationVoter(addr sdk.AccAddress) MsgRevocationVoter {
	return MsgRevocationVoter{
		Address: addr,
	}
}

func (msg MsgRevocationVoter) Type() string                 { return MsgTypeName }
func (msg MsgRevocationVoter) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Address} }

func (msg MsgRevocationVoter) ValidateBasic() sdk.Error {
	if msg.Address == nil {
		return sdk.ErrInvalidAddress("MsgRevocationVoter.Address must not be empty")
	}

	return nil
}

func (msg MsgRevocationVoter) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(bz)
}
