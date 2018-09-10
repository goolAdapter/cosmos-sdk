package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "lottery" type messages.
func NewHandler() sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		return sdk.Result{}
	}

}
