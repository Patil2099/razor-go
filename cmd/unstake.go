package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math/big"
	"razor/core"
	"razor/core/types"
	"razor/utils"
)

var unstakeCmd = &cobra.Command{
	Use:   "unstake",
	Short: "Unstake your razors",
	Long: `unstake allows user to unstake their razors in the razor network
	For ex:
	unstake --address <address> --password <password>
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := GetConfigData()
		if err != nil {
			log.Fatal("Error in getting config: ", err)
		}
		address, _ := cmd.Flags().GetString("address")
		password := utils.PasswordPrompt()

		client := utils.ConnectToClient(config.Provider)

		balance, err := utils.FetchBalance(client, address)
		if err != nil {
			log.Fatal("Error in fetching balance: ", err)
		}
		if balance.Cmp(big.NewInt(0)) == 0 {
			log.Fatal("Account balance is 0. Cannot unstake...")
		}

		epoch, err := WaitForCommitState(client, address, "unstake")
		if err != nil {
			log.Fatal(err)
		}
		stakeManager := utils.GetStakeManager(client)
		txnOpts := utils.GetTxnOpts(types.TransactionOptions{
			Client:         client,
			Password:       password,
			AccountAddress: address,
			ChainId:        core.ChainId,
			GasMultiplier:  config.GasMultiplier,
		})
		log.Info("Unstaking coins")
		txn, err := stakeManager.Unstake(txnOpts, epoch)
		if err != nil {
			log.Fatal("Error in un-staking: ", err)
		}
		log.Info("Successfully unstaked all the tokens")
		log.Info("Transaction hash: ", txn.Hash())
		utils.WaitForBlockCompletion(client,fmt.Sprintf("%s",txn.Hash()))
	},
}

func init() {
	rootCmd.AddCommand(unstakeCmd)

	var Address  string

	unstakeCmd.Flags().StringVarP(&Address, "address", "", "", "address of the staker")

	unstakeCmd.MarkFlagRequired("address")
}
