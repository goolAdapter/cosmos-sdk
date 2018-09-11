package voter

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//voter want to vote
type Voter struct {
	Address sdk.AccAddress `json:"address"`
}

func NewVoter(addr sdk.AccAddress) Voter {
	return Voter{
		Address: addr,
	}
}

func (ac Voter) GetAddress() sdk.AccAddress {
	return ac.Address
}

func (ac *Voter) SetAddress(addr sdk.AccAddress) error {
	if len(ac.Address) != 0 {
		return errors.New("cannot override Voter address")
	}
	ac.Address = addr
	return nil
}
