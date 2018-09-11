package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/examples/lottery"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagActor = "actor"
)

func SetupCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup a new lottery status",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(lottery.GetAccountDecoder(cdc))

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			actorsStr := viper.GetString(flagActor)
			if actorsStr == "" {
				err = fmt.Errorf("flag %s have not specific.", flagActor)
				return err
			}

			res, err := cliCtx.QueryStore(lottery.GetInfoSequenceKey(from, cdc), "main")
			if err != nil {
				return err
			} else if len(res) == 0 {
				return err
			}

			actors := strings.Split(actorsStr, "#")
			msg := lottery.NewMsgSetupLottery(from, actors)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagActor, "", "actor list, like A#B#C#D")
	return cmd
}

const (
	flagNumber    = "number"
	flagRoundMemo = "roundmemo"
)

func RoundCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "round",
		Short: "start a new round",
		RunE: func(cmd *cobra.Command, args []string) error {
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(lottery.GetAccountDecoder(cdc))

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			memo := viper.GetString(flagRoundMemo)

			num := viper.GetInt64(flagNumber)
			if num <= 0 {
				err = fmt.Errorf("flag %s have invalid value %d", flagNumber, num)
				return err
			}

			msg := lottery.NewMsgStartLotteryRound(from, 0, num, memo)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int64(flagNumber, 0, "how many prize number to generate on this round")
	cmd.Flags().String(flagRoundMemo, "", "memo to this round")

	return cmd
}
