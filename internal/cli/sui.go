package cli

import (
	"fmt"
	"sync"

	"github.com/say8hi/addrforge/internal/sui"
	"github.com/say8hi/addrforge/internal/util"
	"github.com/say8hi/addrforge/internal/worker"
	"github.com/spf13/cobra"
)

var (
	suiWorkers         int
	suiCaseInsensitive bool
	suiOutputFile      string
)

var suiCmd = &cobra.Command{
	Use:   "sui <prefix>",
	Short: "Generate Sui address with prefix",
	Args:  cobra.RangeArgs(1, 5),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate prefixes
		for _, prefix := range args {
			if err := sui.ValidatePrefix(prefix); err != nil {
				return fmt.Errorf("invalid prefix '%s': %w", prefix, err)
			}
		}

		fmt.Printf("Searching for Sui address with prefixes: %v\n", args)
		if suiCaseInsensitive {
			fmt.Println("Using case-insensitive matching")
		}

		var outputMutex sync.Mutex

		return worker.RunWithResult(suiWorkers, func(id int) (bool, error) {
			wallet, err := sui.GenerateWallet()
			if err != nil {
				return false, err
			}

			for _, prefix := range args {
				var matches bool
				if suiCaseInsensitive {
					matches = util.Match(wallet.Address, prefix)
				} else {
					matches = util.MatchCaseSensitive(wallet.Address, prefix)
				}

				if matches {
					outputMutex.Lock()
					result := fmt.Sprintf(
						"MATCH FOUND!\nAddress: %s\nPrivate: %s\nWorker: %d\n",
						wallet.Address,
						wallet.PrivateKey,
						id,
					)

					if suiOutputFile != "" {
						util.SaveResult(suiOutputFile, result)
					} else {
						fmt.Print(result)
					}
					outputMutex.Unlock()
					return true, nil
				}
			}
			return false, nil
		}, nil)
	},
}

func init() {
	suiCmd.Flags().IntVarP(&suiWorkers, "workers", "w", 4, "Number of workers")
	suiCmd.Flags().
		BoolVarP(&suiCaseInsensitive, "ignore-case", "i", false, "Case-insensitive matching")
	suiCmd.Flags().StringVarP(&suiOutputFile, "output", "o", "", "Output file for results")
}
