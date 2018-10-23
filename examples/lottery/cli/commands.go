package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
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
		Short: "Setup a new lottery status",
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
			var seq int64
			if err != nil {
				return err
			} else if len(res) != 0 {
				cdc.MustUnmarshalBinary(res, &seq)
			}

			actors := strings.Split(actorsStr, "#")

			memo := viper.GetString(client.FlagMemo)
			if len(memo) > 1024 {
				err = fmt.Errorf("flag %s too long, must less 1024.", client.FlagMemo)
				return err
			}

			msg := lottery.NewMsgSetupLottery(from, seq, actors, memo)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagActor, "", "actor list, like A#B#C#D")
	return cmd
}

const (
	flagAmount = "amount"
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

			res, err := cliCtx.QueryStore(lottery.GetInfoSequenceKey(from, cdc), "main")
			var seq int64
			if err != nil {
				return err
			} else if len(res) != 0 {
				cdc.MustUnmarshalBinary(res, &seq)
			}

			memo := viper.GetString(client.FlagMemo)
			if len(memo) > 1024 {
				err = fmt.Errorf("flag %s too long, must less 1024.", client.FlagMemo)
				return err
			}

			amount := viper.GetInt64(flagAmount)
			if amount <= 0 || amount > 10240 {
				err = fmt.Errorf("flag %s have invalid value %d", flagAmount, amount)
				return err
			}

			msg := lottery.NewMsgStartLotteryRound(from, seq+1, amount, memo)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int64(flagAmount, 0, "how many number to generate on this round")

	return cmd
}

var flagTargetAddress = "target"

func ResultCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "result",
		Short: "fetch result",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(lottery.GetAccountDecoder(cdc))

			bech32addr := viper.GetString(flagTargetAddress)
			target, err := sdk.AccAddressFromBech32(bech32addr)
			if err != nil {
				return err
			}

			res, err := cliCtx.QueryStore(lottery.GetInfoStatusKey(target, cdc), "main")
			var status int64
			if err != nil {
				return err
			} else if len(res) != 0 {
				cdc.MustUnmarshalBinary(res, &status)
			}

			if status != 3 {
				return fmt.Errorf("round not finish")
			}

			res, err = cliCtx.QueryStore(lottery.GetInfoSequenceKey(target, cdc), "main")
			var seq int64
			if err != nil {
				return err
			} else if len(res) != 0 {
				cdc.MustUnmarshalBinary(res, &seq)
			}

			res, err = cliCtx.QueryStore(lottery.GetInfoResultKey(target, seq, cdc), "main")
			var result int64
			if err != nil {
				return err
			} else if len(res) != 0 {
				cdc.MustUnmarshalBinary(res, &result)
			}

			fmt.Printf("result is %v", result)
			return nil
		},
	}

	cmd.Flags().String(flagTargetAddress, "", "target address")

	return cmd
}
