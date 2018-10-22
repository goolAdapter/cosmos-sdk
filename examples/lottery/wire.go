package lottery

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec for default AppAccount
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSetupLottery{}, "lottery/MsgSetupLottery", nil)
	cdc.RegisterConcrete(MsgStartLotteryRound{}, "lottery/MsgStartLotteryRound", nil)
	cdc.RegisterConcrete(MsgPreVote{}, "lottery/MsgPreVote", nil)
	cdc.RegisterConcrete(MsgVote{}, "lottery/MsgVote", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
