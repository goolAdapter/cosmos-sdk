package voter

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

var (
	VoterPrefixKey = []byte{0x02}
)

type VoterKeeper struct {

	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of voters.
	cdc *wire.Codec
}

func NewVoterKeeper(cdc *wire.Codec, key sdk.StoreKey) VoterKeeper {
	return VoterKeeper{
		key: key,
		cdc: cdc,
	}
}

// Turn an address to key used to get it from the account store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(VoterPrefixKey, addr.Bytes()...)
}

func (ak VoterKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) (Voter, error) {
	store := ctx.KVStore(ak.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return Voter{}, fmt.Errorf("%v 's Voter not exist", addr)
	}

	ac := MustUnmarshalVoter(ak.cdc, addr, bz)
	return ac, nil
}

func (ak VoterKeeper) SetAccount(ctx sdk.Context, ac Voter) {
	addr := ac.GetAddress()
	store := ctx.KVStore(ak.key)
	bz := MustMarshalVoter(ak.cdc, ac)
	store.Set(AddressStoreKey(addr), bz)
}

func (ak VoterKeeper) IterateAccounts(ctx sdk.Context, process func(Voter) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iter := sdk.KVStorePrefixIterator(store, VoterPrefixKey)
	defer iter.Close()

	for ; !iter.Valid(); iter.Next() {
		acc := MustUnmarshalVoter(ak.cdc, iter.Key()[1:], iter.Value())
		if process(acc) {
			return
		}
	}
}
