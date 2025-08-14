package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "addrforge",
	Short: "Ethereum, Solana, and Sui vanity address generator",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(ethCmd)
	rootCmd.AddCommand(solCmd)
	rootCmd.AddCommand(suiCmd)
}
