package lottery

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteItem struct {
	Address sdk.AccAddress `json:"address"`
	Hash    []byte         `jsong:"hash"`
	Value   int64          `jsong:"value"`
}

type VoteItemList []VoteItem

//----------------------------------------
// Sort interface

//nolint
func (vl VoteItemList) Len() int           { return len(vl) }
func (vl VoteItemList) Less(i, j int) bool { return vl[i].Address.String() < vl[j].Address.String() }
func (vl VoteItemList) Swap(i, j int)      { vl[i], vl[j] = vl[j], vl[i] }

var _ sort.Interface = VoteItemList{}

// Sort is a helper function to sort the set of coins inplace
func (vl VoteItemList) Sort() VoteItemList {
	sort.Sort(vl)
	return vl
}
