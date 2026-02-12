package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var modeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Show current policy mode (read-only in Phase 1)",
	Long:  `Display the current policy mode for each route. Mode switching comes in a future release.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		c, err := client.New()
		if err != nil {
			return err
		}

		data, code, err := c.Get("/admin/routes")
		if err != nil {
			return fmt.Errorf("cannot reach gateway at %s: %w", cfg.Gateway, err)
		}
		if code != 200 {
			return fmt.Errorf("API returned HTTP %d: %s", code, string(data))
		}

		if flagJSON {
			fmt.Println(string(data))
			return nil
		}

		var routes []struct {
			Path   string `json:"path"`
			Policy string `json:"policy"`
			Name   string `json:"name"`
		}
		json.Unmarshal(data, &routes)

		// Also try wrapped response
		if len(routes) == 0 {
			var wrapped struct {
				Routes []struct {
					Path   string `json:"path"`
					Policy string `json:"policy"`
					Name   string `json:"name"`
				} `json:"routes"`
			}
			json.Unmarshal(data, &wrapped)
			routes = wrapped.Routes
		}

		fmt.Println("Policy Modes")
		fmt.Println("─────────────────────────────")

		icons := map[string]string{
			"observe":    "👁  Observe",
			"chargeback": "👁  Observe",
			"control":    "🎛  Control",
			"fiat402":    "🎛  Control",
			"charge":     "💲 Charge",
			"l402":       "💲 Charge",
			"public":     "🔓 Public",
		}

		for _, r := range routes {
			mode := r.Policy
			display, ok := icons[mode]
			if !ok {
				display = mode
			}
			name := r.Path
			if r.Name != "" {
				name = r.Name + " (" + r.Path + ")"
			}
			fmt.Printf("  %-40s %s\n", name, display)
		}

		if len(routes) == 0 {
			fmt.Println("  No routes configured")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(modeCmd)
}
