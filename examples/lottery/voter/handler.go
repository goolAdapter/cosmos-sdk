package voter

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "voter" type messages.
func NewHandler(k VoterKeeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateVoter:
			return handleMsgCreateVoter(ctx, k, msg)
		case MsgRevokeVoter:
			return handleMsgRevokeVoter(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized cool Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}

}

func handleMsgCreateVoter(ctx sdk.Context, k VoterKeeper, msg MsgCreateVoter) sdk.Result {
	ac := NewVoter(msg.Address, msg.Memo)
	k.SetVoter(ctx, ac)

	return sdk.Result{}
}

func handleMsgRevokeVoter(ctx sdk.Context, k VoterKeeper, msg MsgRevokeVoter) sdk.Result {
	k.DeleteVoter(ctx, msg.Address)

	return sdk.Result{}
}
