package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate reports (threats, spend, compliance)",
}

var reportThreatsCmd = &cobra.Command{
	Use:   "threats",
	Short: "Show blocked requests, anomalies, and threat summary",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		c, err := client.New()
		if err != nil {
			return err
		}

		data, code, err := c.Get("/admin/reports/threats")
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

		var resp struct {
			TotalBlocked int `json:"total_blocked"`
			Categories   []struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			} `json:"categories"`
			RecentThreats []struct {
				Time   string `json:"time"`
				Type   string `json:"type"`
				Agent  string `json:"agent"`
				Route  string `json:"route"`
				Action string `json:"action"`
			} `json:"recent_threats"`
		}
		if err := json.Unmarshal(data, &resp); err == nil && resp.TotalBlocked > 0 {
			fmt.Println("Threat Report")
			fmt.Println("─────────────────────────────")
			fmt.Printf("  Total Blocked: %d\n\n", resp.TotalBlocked)

			if len(resp.Categories) > 0 {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "CATEGORY\tCOUNT")
				fmt.Fprintln(w, "────────\t─────")
				for _, cat := range resp.Categories {
					fmt.Fprintf(w, "%s\t%d\n", cat.Name, cat.Count)
				}
				w.Flush()
				fmt.Println()
			}

			if len(resp.RecentThreats) > 0 {
				fmt.Println("Recent Threats")
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "TIME\tTYPE\tAGENT\tROUTE\tACTION")
				for _, t := range resp.RecentThreats {
					fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.Time, t.Type, t.Agent, t.Route, t.Action)
				}
				w.Flush()
			}
			return nil
		}

		// Fallback
		var raw interface{}
		json.Unmarshal(data, &raw)
		out, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	reportCmd.AddCommand(reportThreatsCmd)
	rootCmd.AddCommand(reportCmd)
}
