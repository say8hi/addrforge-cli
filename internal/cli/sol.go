package cli

import (
	"fmt"
	"sync"

	"github.com/say8hi/addrforge/internal/sol"
	"github.com/say8hi/addrforge/internal/util"
	"github.com/say8hi/addrforge/internal/worker"
	"github.com/spf13/cobra"
)

var (
	solWorkers         int
	solCaseInsensitive bool
	solOutputFile      string
)

var solCmd = &cobra.Command{
	Use:   "sol <prefix>",
	Short: "Generate Solana address with prefix",
	Args:  cobra.RangeArgs(1, 5),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate prefixes
		for _, prefix := range args {
			if err := validateSolanaPrefix(prefix); err != nil {
				return fmt.Errorf("invalid prefix '%s': %w", prefix, err)
			}
		}

		fmt.Printf("Searching for Solana address with prefixes: %v\n", args)
		if solCaseInsensitive {
			fmt.Println("Using case-insensitive matching")
		}

		var outputMutex sync.Mutex

		return worker.RunWithResult(solWorkers, func(id int) (bool, error) {
			wallet, err := sol.GenerateWallet()
			if err != nil {
				return false, err
			}

			for _, prefix := range args {
				var matches bool
				if solCaseInsensitive {
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

					if solOutputFile != "" {
						util.SaveResult(solOutputFile, result)
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

func validateSolanaPrefix(prefix string) error {
	// Base58 alphabet
	const base58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

	for _, char := range prefix {
		found := false
		for _, b58char := range base58 {
			if char == b58char {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("contains invalid Base58 character: %c", char)
		}
	}
	return nil
}

func init() {
	solCmd.Flags().IntVarP(&solWorkers, "workers", "w", 4, "Number of workers")
	solCmd.Flags().
		BoolVarP(&solCaseInsensitive, "ignore-case", "i", false, "Case-insensitive matching")
	solCmd.Flags().StringVarP(&solOutputFile, "output", "o", "", "Output file for results")
}
