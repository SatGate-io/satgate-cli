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

var (
	spendAgent  string
	spendPeriod string
)

var spendCmd = &cobra.Command{
	Use:   "spend",
	Short: "Show spend summary (org-wide or per-agent)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		c, err := client.New()
		if err != nil {
			return err
		}

		path := "/admin/spend"
		if spendAgent != "" {
			path += "?agent=" + spendAgent
		}
		if spendPeriod != "" {
			sep := "?"
			if spendAgent != "" {
				sep = "&"
			}
			path += sep + "period=" + spendPeriod
		}

		data, code, err := c.Get(path)
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

		// Try to parse as org summary
		var orgResp struct {
			TotalAllocated float64 `json:"total_allocated"`
			TotalConsumed  float64 `json:"total_consumed"`
			Agents         []struct {
				Name   string  `json:"name"`
				Spent  float64 `json:"spent"`
				Budget float64 `json:"budget"`
			} `json:"agents"`
		}
		if err := json.Unmarshal(data, &orgResp); err == nil && (orgResp.TotalAllocated > 0 || len(orgResp.Agents) > 0) {
			fmt.Println("Spend Summary")
			fmt.Println("─────────────────────────────")
			if orgResp.TotalAllocated > 0 {
				pct := 0.0
				if orgResp.TotalAllocated > 0 {
					pct = (orgResp.TotalConsumed / orgResp.TotalAllocated) * 100
				}
				fmt.Printf("  Allocated:  $%.2f\n", orgResp.TotalAllocated)
				fmt.Printf("  Consumed:   $%.2f (%.1f%%)\n", orgResp.TotalConsumed, pct)
				fmt.Println()
			}

			if len(orgResp.Agents) > 0 {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "AGENT\tSPENT\tBUDGET\tUTILIZATION")
				fmt.Fprintln(w, "─────\t─────\t──────\t───────────")
				for _, a := range orgResp.Agents {
					util := "—"
					if a.Budget > 0 {
						util = fmt.Sprintf("%.1f%%", (a.Spent/a.Budget)*100)
					}
					budget := "unlimited"
					if a.Budget > 0 {
						budget = fmt.Sprintf("$%.2f", a.Budget)
					}
					fmt.Fprintf(w, "%s\t$%.2f\t%s\t%s\n", a.Name, a.Spent, budget, util)
				}
				w.Flush()
			}
			return nil
		}

		// Fallback: just pretty-print whatever we got
		var raw interface{}
		json.Unmarshal(data, &raw)
		out, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println(string(out))
		return nil
	},
}

func init() {
	spendCmd.Flags().StringVar(&spendAgent, "agent", "", "filter by agent name")
	spendCmd.Flags().StringVar(&spendPeriod, "period", "", "time period (e.g. 7d, 30d)")
	rootCmd.AddCommand(spendCmd)
}
