package cmd

import (
	"fmt"
	"os"

	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	version   string
	buildTime string
	cfgFile   string
	flagJSON  bool
	flagYes   bool
	flagDry   bool
)

func SetVersionInfo(v, b string) {
	version = v
	buildTime = b
}

var rootCmd = &cobra.Command{
	Use:   "satgate",
	Short: "SatGate CLI — Manage your API's economic firewall",
	Long: `SatGate CLI wraps the SatGate Admin API for server-side API operators.

Mint tokens, track spend, revoke agents, and view security reports
from the terminal. The server-side counterpart to lnget.

They're the wallet. We're the register.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load config before every command
		config.Load(cfgFile)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ~/.satgate/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&flagYes, "yes", false, "skip confirmation prompts")
	rootCmd.PersistentFlags().BoolVar(&flagDry, "dry-run", false, "show what would happen without executing")
}

// printTarget prints the target gateway info before mutating commands
func printTarget(cfg *config.Config) {
	fmt.Fprintf(os.Stderr, "⚡ Target: %s (%s)", cfg.Gateway, cfg.Surface)
	if cfg.Tenant != "" && cfg.Tenant != "default" {
		fmt.Fprintf(os.Stderr, " tenant=%s", cfg.Tenant)
	}
	fmt.Fprintln(os.Stderr)
}

// confirmAction asks for y/N confirmation. Returns true if confirmed.
func confirmAction(prompt string) bool {
	if flagYes {
		return true
	}
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}
