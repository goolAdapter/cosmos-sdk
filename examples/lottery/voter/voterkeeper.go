package voter

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

var (
	VoterPrefixKey = []byte{0x02}
)

type VoterKeeper struct {

	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The wire codec for binary encoding/decoding of voters.
	cdc *wire.Codec

	accountMapper auth.AccountMapper
}

func NewVoterKeeper(cdc *wire.Codec, key sdk.StoreKey, am auth.AccountMapper) VoterKeeper {
	return VoterKeeper{
		key:           key,
		cdc:           cdc,
		accountMapper: am,
	}
}

// Turn an address to key used to get it from the voter store
func AddressStoreKey(addr sdk.AccAddress) []byte {
	return append(VoterPrefixKey, addr.Bytes()...)
}

func (ak VoterKeeper) GetVoter(ctx sdk.Context, addr sdk.AccAddress) (Voter, error) {
	store := ctx.KVStore(ak.key)
	bz := store.Get(AddressStoreKey(addr))
	if bz == nil {
		return Voter{}, fmt.Errorf("%v 's Voter not exist", addr)
	}

	ac := MustUnmarshalVoter(ak.cdc, addr, bz)
	return ac, nil
}

func (ak VoterKeeper) SetVoter(ctx sdk.Context, ac Voter) {
	addr := ac.GetAddress()

	//auto create account
	account := ak.accountMapper.GetAccount(ctx, ac.Address)
	if account == nil {
		ak.accountMapper.NewAccountWithAddress(ctx, addr)
	}

	store := ctx.KVStore(ak.key)
	bz := MustMarshalVoter(ak.cdc, ac)
	store.Set(AddressStoreKey(addr), bz)
}

func (ak VoterKeeper) DeleteVoter(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(ak.key)
	store.Delete(AddressStoreKey(addr))
}

func (ak VoterKeeper) IterateVoters(ctx sdk.Context, process func(Voter) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iter := sdk.KVStorePrefixIterator(store, VoterPrefixKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		acc := MustUnmarshalVoter(ak.cdc, iter.Key()[1:], iter.Value())
		if process(acc) {
			return
		}
	}
}
