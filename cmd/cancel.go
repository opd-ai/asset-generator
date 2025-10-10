package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cancelAll bool
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel ongoing image generation",
	Long: `Cancel the current generation or all queued generations.

By default, cancels the current generation in progress. Use --all to cancel
all queued generations.

The cancel command sends an interrupt signal to the SwarmUI server to stop
generation immediately. This is useful for:
  - Stopping long-running generations (e.g., Flux models)
  - Clearing a backlog of queued generations
  - Recovering from stuck generations

Examples:
  # Cancel the current generation
  asset-generator cancel
  
  # Cancel all queued generations
  asset-generator cancel --all
  
  # Cancel with verbose output
  asset-generator cancel -v`,
	RunE: runCancel,
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	cancelCmd.Flags().BoolVar(&cancelAll, "all", false, "Cancel all queued generations instead of just the current one")
}

func runCancel(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if cancelAll {
		if !quiet {
			fmt.Println("Cancelling all queued generations...")
		}

		if err := assetClient.InterruptAll(ctx); err != nil {
			return fmt.Errorf("failed to cancel generations: %w", err)
		}

		if !quiet {
			fmt.Println("✓ Successfully cancelled all queued generations")
		}
	} else {
		if !quiet {
			fmt.Println("Cancelling current generation...")
		}

		if err := assetClient.Interrupt(ctx); err != nil {
			return fmt.Errorf("failed to cancel generation: %w", err)
		}

		if !quiet {
			fmt.Println("✓ Successfully cancelled current generation")
		}
	}

	return nil
}
