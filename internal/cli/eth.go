package cli

import (
	"fmt"

	"github.com/say8hi/addrforge/internal/eth"
	"github.com/say8hi/addrforge/internal/util"
	"github.com/say8hi/addrforge/internal/worker"

	"github.com/spf13/cobra"
)

var ethWorkers int

var ethCmd = &cobra.Command{
	Use:   "eth <prefix>",
	Short: "Generate Ethereum address with prefix",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix := args[0]
		fmt.Println("Searching for Ethereum address with prefix:", prefix)

		return worker.Run(ethWorkers, func(id int) (bool, error) {
			wallet, err := eth.GenerateWallet()
			if err != nil {
				return false, err
			}
			if util.Match(wallet.Address, prefix) {
				fmt.Printf(
					"MATCH FOUND!\nAddress: %s\nPrivate: %s\n",
					wallet.Address,
					wallet.PrivateKey,
				)
				return true, nil
			}
			return false, nil
		})
	},
}

func init() {
	ethCmd.Flags().IntVarP(&ethWorkers, "workers", "w", 4, "Workers amount")
}
