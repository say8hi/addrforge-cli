package cli

import (
	"fmt"

	"github.com/say8hi/addrforge/internal/sol"
	"github.com/say8hi/addrforge/internal/util"
	"github.com/say8hi/addrforge/internal/worker"

	"github.com/spf13/cobra"
)

var solWorkers int

var solCmd = &cobra.Command{
	Use:   "sol <prefix>",
	Short: "Generate Solana address with prefix",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prefix := args[0]
		fmt.Println("Searching for Solana address with prefix:", prefix)

		return worker.Run(solWorkers, func(id int) (bool, error) {
			wallet, err := sol.GenerateWallet()
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
	solCmd.Flags().IntVarP(&solWorkers, "workers", "w", 4, "Workers amount")
}
