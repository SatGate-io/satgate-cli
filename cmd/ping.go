package cmd

import (
	"fmt"
	"os"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Quick liveness check (exit code 0 = healthy, 1 = unreachable)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		c, err := client.New()
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s unreachable: %v\n", cfg.Gateway, err)
			os.Exit(1)
			return nil
		}

		path := "/admin/ping"
		if c.Surface() == "cloud" {
			path = "/healthz"
		}

		_, code, err := c.Get(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ %s unreachable: %v\n", cfg.Gateway, err)
			os.Exit(1)
			return nil
		}

		if code == 200 {
			fmt.Printf("✓ %s is healthy\n", cfg.Gateway)
		} else {
			fmt.Fprintf(os.Stderr, "✗ %s returned HTTP %d\n", cfg.Gateway, code)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
