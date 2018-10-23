package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/examples/lottery"
	"github.com/cosmos/cosmos-sdk/examples/lottery/app"
	lotterycmd "github.com/cosmos/cosmos-sdk/examples/lottery/cli"
	votercmd "github.com/cosmos/cosmos-sdk/examples/lottery/voter/cli"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "lotterycli",
		Short: "lottery light-client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	// add standard rpc, and tx commands
	rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(
			authcmd.GetAccountCmd("acc", cdc, lottery.GetAccountDecoder(cdc)),
			votercmd.QueryVotersCmd("main", cdc),
		)...)

	rootCmd.AddCommand(
		client.PostCommands(
			lotterycmd.SetupCmd(cdc),
			lotterycmd.RoundCmd(cdc),
			lotterycmd.ResultCmd(cdc),
			lotterycmd.RunVoterCmd(cdc),
			votercmd.CreateVoterCmd(cdc),
			votercmd.RevokeVoterCmd(cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		lcd.ServeCommand(cdc),
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "LT", os.ExpandEnv("$HOME/.lotterycli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}
