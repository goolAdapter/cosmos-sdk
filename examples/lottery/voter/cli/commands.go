package cli

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/examples/lottery"
	"github.com/cosmos/cosmos-sdk/examples/lottery/voter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagVoterMemo = "voterMemo"
)

func CreateVoterCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-voter",
		Short: "Create a new voter",
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

			memo := viper.GetString(flagVoterMemo)
			if len(memo) > 1024 {
				err = fmt.Errorf("flag %s too long, must less 1024.", flagVoterMemo)
				return err
			}

			msg := voter.NewMsgCreateVoter(from, memo)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagVoterMemo, "", "voter's description")
	return cmd
}

func RevokeVoterCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke-voter",
		Short: "Revocation a voter",
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

			msg := voter.NewMsgRevokeVoter(from)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func QueryVotersCmd(storeName string, cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voters",
		Short: "Query for all voters",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := voter.VoterPrefixKey
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, err := cliCtx.QuerySubspace(key, storeName)
			if err != nil {
				return err
			}

			// parse out the voters
			var voters []voter.Voter
			for _, kv := range resKVs {
				addr := kv.Key[1:]
				voter := voter.MustUnmarshalVoter(cdc, addr, kv.Value)
				voters = append(voters, voter)
			}

			switch viper.Get(cli.OutputFlag) {
			case "text":
				for _, voter := range voters {
					resp, err := voter.HumanReadableString()
					if err != nil {
						return err
					}

					fmt.Println(resp)
				}
			case "json":
				output, err := wire.MarshalJSONIndent(cdc, voters)
				if err != nil {
					return err
				}

				fmt.Println(string(output))
				return nil
			}

			// TODO: output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}
