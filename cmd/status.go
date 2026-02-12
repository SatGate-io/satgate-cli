package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show gateway health, version, and uptime",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		c, err := client.New()
		if err != nil {
			return err
		}

		data, code, err := c.Get("/admin/ping")
		if err != nil {
			return fmt.Errorf("cannot reach gateway at %s: %w", cfg.Gateway, err)
		}

		if flagJSON {
			// Enrich with CLI metadata
			var resp map[string]interface{}
			json.Unmarshal(data, &resp)
			resp["cli_version"] = version
			resp["cli_build_time"] = buildTime
			resp["surface"] = cfg.Surface
			resp["gateway"] = cfg.Gateway
			out, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		var resp map[string]interface{}
		json.Unmarshal(data, &resp)

		fmt.Println("SatGate Gateway Status")
		fmt.Println("─────────────────────────────")
		fmt.Printf("  Gateway:     %s\n", cfg.Gateway)
		fmt.Printf("  Surface:     %s\n", cfg.Surface)
		fmt.Printf("  HTTP Status: %d\n", code)

		if v, ok := resp["version"]; ok {
			fmt.Printf("  Version:     %v\n", v)
		}
		if v, ok := resp["uptime"]; ok {
			fmt.Printf("  Uptime:      %v\n", v)
		}
		if v, ok := resp["status"]; ok {
			fmt.Printf("  Status:      %v\n", v)
		}
		if v, ok := resp["mode"]; ok {
			fmt.Printf("  Mode:        %v\n", v)
		}

		fmt.Println("─────────────────────────────")
		fmt.Printf("  CLI Version: %s (%s)\n", version, buildTime)

		if code != 200 {
			fmt.Fprintf(os.Stderr, "\n⚠️  Gateway returned HTTP %d\n", code)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
