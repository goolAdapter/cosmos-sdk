package lottery

import (
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handle all "lottery" type messages.
func NewHandler(k LotteryKeeper) sdk.Handler {

	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetupLottery:
			return handleMsgSetupLottery(ctx, k, msg)
		case MsgStartLotteryRound:
			return handleMsgStartLotteryRound(ctx, k, msg)
		case MsgPreVote:
			return handleMsgPreVote(ctx, k, msg)
		case MsgVote:
			return handleMsgVote(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized cool Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}

}

func handleMsgSetupLottery(ctx sdk.Context, k LotteryKeeper, msg MsgSetupLottery) sdk.Result {
	err := k.SetupLottery(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgStartLotteryRound(ctx sdk.Context, k LotteryKeeper, msg MsgStartLotteryRound) sdk.Result {
	err := k.StartLotteryRound(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgPreVote(ctx sdk.Context, k LotteryKeeper, msg MsgPreVote) sdk.Result {
	err := k.HandlePreVoteMsg(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgVote(ctx sdk.Context, k LotteryKeeper, msg MsgVote) sdk.Result {
	err := k.HandleVoteMsg(ctx, msg)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
