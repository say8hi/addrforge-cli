package cli

import (
	"fmt"
	"sync"

	"github.com/say8hi/addrforge/internal/eth"
	"github.com/say8hi/addrforge/internal/util"
	"github.com/say8hi/addrforge/internal/worker"
	"github.com/spf13/cobra"
)

var (
	ethWorkers         int
	ethCaseInsensitive bool
	ethOutputFile      string
)

var ethCmd = &cobra.Command{
	Use:   "eth <prefix>",
	Short: "Generate Ethereum address with prefix",
	Args:  cobra.RangeArgs(1, 5),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate prefixes
		for _, prefix := range args {
			if err := validateEthereumPrefix(prefix); err != nil {
				return fmt.Errorf("invalid prefix '%s': %w", prefix, err)
			}
		}

		fmt.Printf("Searching for Ethereum address with prefixes: %v\n", args)
		if ethCaseInsensitive {
			fmt.Println("Using case-insensitive matching")
		}

		var outputMutex sync.Mutex

		return worker.RunWithResult(ethWorkers, func(id int) (bool, error) {
			wallet, err := eth.GenerateWallet()
			if err != nil {
				return false, err
			}

			for _, prefix := range args {
				var matches bool
				if ethCaseInsensitive {
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

					if ethOutputFile != "" {
						util.SaveResult(ethOutputFile, result)
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

func validateEthereumPrefix(prefix string) error {
	// Remove 0x if present
	if len(prefix) > 2 && prefix[:2] == "0x" {
		prefix = prefix[2:]
	}

	// Check hex characters
	for _, char := range prefix {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return fmt.Errorf("contains invalid hex character: %c", char)
		}
	}
	return nil
}

func init() {
	ethCmd.Flags().IntVarP(&ethWorkers, "workers", "w", 4, "Number of workers")
	ethCmd.Flags().
		BoolVarP(&ethCaseInsensitive, "ignore-case", "i", false, "Case-insensitive matching")
	ethCmd.Flags().StringVarP(&ethOutputFile, "output", "o", "", "Output file for results")
}
