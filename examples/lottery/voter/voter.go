package voter

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

//voter represent who want to vote
type Voter struct {
	Address sdk.AccAddress `json:"address"`
	Memo    string         `json:"memo"`
}

func NewVoter(addr sdk.AccAddress, memo string) Voter {
	return Voter{
		Address: addr,
		Memo:    memo,
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

type voterValue struct {
	Memo string
}

// return the redelegation without fields contained within the key for the store
func MustMarshalVoter(cdc *wire.Codec, v Voter) []byte {
	val := voterValue{
		Memo: v.Memo,
	}
	return cdc.MustMarshalBinary(val)
}

// unmarshal a redelegation from a store key and value
func MustUnmarshalVoter(cdc *wire.Codec, ownerAddr, value []byte) Voter {
	v, err := UnmarshalVoter(cdc, ownerAddr, value)
	if err != nil {
		panic(err)
	}

	return v
}

// unmarshal a redelegation from a store key and value
func UnmarshalVoter(cdc *wire.Codec, ownerAddr, value []byte) (v Voter, err error) {
	var storeValue Voter
	err = cdc.UnmarshalBinary(value, &storeValue)
	if err != nil {
		return
	}

	if len(ownerAddr) != 20 {
		err = errors.New("unexpected address length")
		return
	}

	return Voter{
		Address: ownerAddr,
		Memo:    storeValue.Memo,
	}, nil
}

// HumanReadableString returns a human readable string representation of a
// validator. An error is returned if the owner or the owner's public key
// cannot be converted to Bech32 format.
func (v Voter) HumanReadableString() (string, error) {
	resp := "Voter \n"
	resp += fmt.Sprintf("Address: %s\n", v.Address)
	resp += fmt.Sprintf("Memo: %s\n", v.Memo)

	return resp, nil
}
