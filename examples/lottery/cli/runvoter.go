package cli

import (
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/examples/lottery"
	wire "github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/tendermint/tendermint/libs/log"
	"golang.org/x/crypto/ripemd160"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type runVoterCommander struct {
	cdc    *wire.Codec
	logger log.Logger

	Value int64
}

var flagTargetAddress = "target"

func RunVoterCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := &runVoterCommander{
		cdc:    cdc,
		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}

	cmd := &cobra.Command{
		Use: "run-voter",
		Run: cmdr.runVoter,
	}

	cmd.Flags().String(flagTargetAddress, "", "target address")

	return cmd
}

func (c *runVoterCommander) runVoter(cmd *cobra.Command, args []string) {
	cliCtx := context.NewCLIContext().WithCodec(c.cdc).WithLogger(os.Stdout).WithAccountDecoder(lottery.GetAccountDecoder(c.cdc))

	passphrase, err := keys.ReadPassphraseFromStdin(cliCtx.FromAddressName)
	if err != nil {
		panic(err)
	}

	bech32addr := viper.GetString(flagTargetAddress)
	target, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		panic(err)
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		panic(err)
	}
OUTER:
	for {
		time.Sleep(1 * time.Second)

		res, err := cliCtx.QueryStore(lottery.GetInfoStatusKey(target, c.cdc), "main")
		var status int64
		if err != nil {
			c.logger.Error("error querying outgoing packet list length", "err", err)
			continue OUTER //TODO replace with continue (I think it should just to the correct place where OUTER is now)
		} else if len(res) != 0 {
			c.cdc.MustUnmarshalBinary(res, &status)
		}

		res, err = cliCtx.QueryStore(lottery.GetInfoSequenceKey(target, c.cdc), "main")
		var seq int64
		if err != nil {
			c.logger.Error("error querying outgoing packet list length", "err", err)
			continue OUTER //TODO replace with continue (I think it should just to the correct place where OUTER is now)
		} else if len(res) != 0 {
			c.cdc.MustUnmarshalBinary(res, &seq)
		}

		switch status {
		case lottery.WaitforPreVotePhase:
			err = c.preVote(cliCtx, target, from, seq, passphrase)
		case lottery.WaitforVotePhase:
			err = c.vote(cliCtx, target, from, seq, passphrase)
		default:
			c.logger.Debug("nothing to do, sleep 1 second")
		}

		if err != nil {
			c.logger.Error("error voting", "err", err)
			continue OUTER //TODO replace with continue (I think it should just to the correct place where OUTER is now)
		}
	}
}

func (c *runVoterCommander) preVote(cliCtx context.CLIContext, target, from sdk.AccAddress, seq int64, passphrase string) error {
	c.logger.Debug("do preVote")
	res, err := cliCtx.QueryStore(lottery.GetInfoVoteKey(target, c.cdc), "main")
	var vl lottery.VoteItemList
	if err != nil {
		c.logger.Error("error querying outgoing packet list length", "err", err)
		return err
	} else if len(res) != 0 {
		c.cdc.MustUnmarshalBinary(res, &vl)
	}

	for _, v := range vl {
		if v.Address.String() == from.String() && len(v.Hash) == 0 {
			c.Value = time.Now().UnixNano()

			hash := c.preVoteHash()
			c.logger.Debug("lock prevote", "value", c.Value)
			c.logger.Debug("send prevote", "hash", hash)
			msg := lottery.NewMsgPreVote(target, from, seq, hash)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return c.sendTxWithPassphrase(cliCtx, passphrase, []sdk.Msg{msg})
		}
	}

	return nil
}

func (c *runVoterCommander) vote(cliCtx context.CLIContext, target, from sdk.AccAddress, seq int64, passphrase string) error {
	c.logger.Debug("do vote")
	res, err := cliCtx.QueryStore(lottery.GetInfoVoteKey(target, c.cdc), "main")
	var vl lottery.VoteItemList
	if err != nil {
		c.logger.Error("error querying outgoing packet list length", "err", err)
		return err
	} else if len(res) != 0 {
		c.cdc.MustUnmarshalBinary(res, &vl)
	}

	for _, v := range vl {
		if v.Address.String() == from.String() && len(v.Hash) != 0 && v.Value == 0 {
			c.logger.Debug("send vote", "value", c.Value)
			msg := lottery.NewMsgVote(target, from, seq, c.Value)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return c.sendTxWithPassphrase(cliCtx, passphrase, []sdk.Msg{msg})
		}
	}

	return nil
}

func (c *runVoterCommander) preVoteHash() []byte {
	// Doesn't write Name, since merkle.SimpleHashFromMap() will
	// include them via the keys.
	bz, _ := c.cdc.MarshalBinary(c.Value) // Does not error
	hasher := ripemd160.New()
	_, err := hasher.Write(bz)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}
	return hasher.Sum(nil)
}

func (c *runVoterCommander) sendTxWithPassphrase(cliCtx context.CLIContext, passphrase string, msgs []sdk.Msg) error {
	if err := cliCtx.EnsureAccountExists(); err != nil {
		return err
	}

	from, err := cliCtx.GetFromAddress()
	if err != nil {
		return err
	}
	txCtx := authctx.NewTxContextFromCLI().WithCodec(c.cdc)

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txCtx.AccountNumber == 0 {
		accNum, err := cliCtx.GetAccountNumber(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if txCtx.Sequence == 0 {
		accSeq, err := cliCtx.GetAccountSequence(from)
		if err != nil {
			return err
		}

		txCtx = txCtx.WithSequence(accSeq)
	}

	// build and sign the transaction
	txBytes, err := txCtx.BuildAndSign(cliCtx.FromAddressName, passphrase, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	return cliCtx.EnsureBroadcastTx(txBytes)
}
