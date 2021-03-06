package client

import (
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/stake/client/cli"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	stakeQueryCmd := &cobra.Command{
		Use:   "stake",
		Short: "Querying commands for the staking module",
	}
	stakeQueryCmd.AddCommand(client.GetCommands(
		cli.GetCmdQueryDelegation(mc.storeKey, mc.cdc),
		cli.GetCmdQueryDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryUnbondingDelegation(mc.storeKey, mc.cdc),
		cli.GetCmdQueryUnbondingDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryRedelegation(mc.storeKey, mc.cdc),
		cli.GetCmdQueryRedelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidator(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidators(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidatorDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidatorUnbondingDelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryValidatorRedelegations(mc.storeKey, mc.cdc),
		cli.GetCmdQueryParams(mc.storeKey, mc.cdc),
		cli.GetCmdQueryPool(mc.storeKey, mc.cdc))...)

	return stakeQueryCmd

}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	stakeTxCmd := &cobra.Command{
		Use:   "stake",
		Short: "Staking transaction subcommands",
	}

	stakeTxCmd.AddCommand(client.PostCommands(
		cli.GetCmdCreateValidator(mc.cdc),
		cli.GetCmdEditValidator(mc.cdc),
		cli.GetCmdDelegate(mc.cdc),
		cli.GetCmdRedelegate(mc.storeKey, mc.cdc),
		cli.GetCmdUnbond(mc.storeKey, mc.cdc),
	)...)

	return stakeTxCmd
}
