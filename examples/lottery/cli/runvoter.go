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
)

type runVoterCommander struct {
	cdc    *wire.Codec
	logger log.Logger
	txCtx  authctx.TxContext
	cliCtx context.CLIContext

	Value int64
}

func RunVoterCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := runVoterCommander{
		cdc:    cdc,
		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		txCtx:  authctx.NewTxContextFromCLI().WithCodec(cdc),
		cliCtx: context.NewCLIContext().WithCodec(cdc).WithLogger(os.Stdout).WithAccountDecoder(lottery.GetAccountDecoder(cdc)),
	}

	cmd := &cobra.Command{
		Use: "run-voter",
		Run: cmdr.runVoter,
	}

	return cmd
}

func (c runVoterCommander) runVoter(cmd *cobra.Command, args []string) {
	passphrase, err := keys.ReadPassphraseFromStdin(c.cliCtx.FromAddressName)
	if err != nil {
		panic(err)
	}

	from, err := c.cliCtx.GetFromAddress()
	if err != nil {
		panic(err)
	}
OUTER:
	for {
		time.Sleep(1 * time.Second)

		res, err := c.cliCtx.QueryStore(lottery.GetInfoStatusKey(from, c.cdc), "main")
		var status int64
		if err != nil {
			c.logger.Error("error querying outgoing packet list length", "err", err)
			continue OUTER //TODO replace with continue (I think it should just to the correct place where OUTER is now)
		} else if len(res) != 0 {
			c.cdc.MustUnmarshalBinary(res, &status)
		}

		switch status {
		case lottery.WaitforPreVotePhase:
			err = c.preVote(from, passphrase)
		case lottery.WaitforVotePhase:
			err = c.vote(from, passphrase)
		}

		if err != nil {
			c.logger.Error("error voting", "err", err)
			continue OUTER //TODO replace with continue (I think it should just to the correct place where OUTER is now)
		}
	}
}

func (c runVoterCommander) preVote(from sdk.AccAddress, passphrase string) error {
	res, err := c.cliCtx.QueryStore(lottery.GetInfoPrevoteKey(from, c.cdc), "main")
	var pvl lottery.PreVoteItemList
	if err != nil {
		c.logger.Error("error querying outgoing packet list length", "err", err)
		return err
	} else if len(res) != 0 {
		c.cdc.MustUnmarshalBinary(res, &pvl)
	}

	for _, v := range pvl {
		if v.Address.String() == from.String() && len(v.Hash) == 0 {
			c.Value = time.Now().UnixNano()

			hash := c.preVoteHash()
			c.logger.Debug("lock prevote value: %v", c.Value)
			c.logger.Debug("send prevote hash: %v", hash)
			msg := lottery.NewMsgPreVote(from, hash)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return c.sendTxWithPassphrase(passphrase, []sdk.Msg{msg})
		}
	}

	return nil
}

func (c runVoterCommander) vote(from sdk.AccAddress, passphrase string) error {
	res, err := c.cliCtx.QueryStore(lottery.GetInfoVoteKey(from, c.cdc), "main")
	var pvl lottery.PreVoteItemList
	if err != nil {
		c.logger.Error("error querying outgoing packet list length", "err", err)
		return err
	} else if len(res) != 0 {
		c.cdc.MustUnmarshalBinary(res, &pvl)
	}

	for _, v := range pvl {
		if v.Address.String() == from.String() && len(v.Hash) == 0 {
			c.logger.Debug("send prevote value: %v", c.Value)
			msg := lottery.NewMsgVote(from, c.Value)

			// Build and sign the transaction, then broadcast to a Tendermint node.
			return c.sendTxWithPassphrase(passphrase, []sdk.Msg{msg})
		}
	}

	return nil
}

func (c runVoterCommander) preVoteHash() []byte {
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

func (c runVoterCommander) sendTxWithPassphrase(passphrase string, msgs []sdk.Msg) error {
	if err := c.cliCtx.EnsureAccountExists(); err != nil {
		return err
	}

	from, err := c.cliCtx.GetFromAddress()
	if err != nil {
		return err
	}

	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if c.txCtx.AccountNumber == 0 {
		accNum, err := c.cliCtx.GetAccountNumber(from)
		if err != nil {
			return err
		}

		c.txCtx = c.txCtx.WithAccountNumber(accNum)
	}

	// TODO: (ref #1903) Allow for user supplied account sequence without
	// automatically doing a manual lookup.
	if c.txCtx.Sequence == 0 {
		accSeq, err := c.cliCtx.GetAccountSequence(from)
		if err != nil {
			return err
		}

		c.txCtx = c.txCtx.WithSequence(accSeq)
	}

	// build and sign the transaction
	txBytes, err := c.txCtx.BuildAndSign(c.cliCtx.FromAddressName, passphrase, msgs)
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	return c.cliCtx.EnsureBroadcastTx(txBytes)
}
