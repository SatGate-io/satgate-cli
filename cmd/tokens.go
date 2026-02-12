package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/spf13/cobra"
)

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "List all tokens with status, spend, and budget remaining",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		data, code, err := c.Get("/admin/tokens")
		if err != nil {
			return err
		}
		if code != 200 {
			return fmt.Errorf("API returned HTTP %d: %s", code, string(data))
		}

		if flagJSON {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Tokens []struct {
				ID        string  `json:"id"`
				Name      string  `json:"name"`
				Status    string  `json:"status"`
				Spent     float64 `json:"spent"`
				Budget    float64 `json:"budget"`
				ExpiresAt string  `json:"expires_at"`
			} `json:"tokens"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			// Try as raw array
			var tokens []struct {
				ID        string  `json:"id"`
				Name      string  `json:"name"`
				Status    string  `json:"status"`
				Spent     float64 `json:"spent"`
				Budget    float64 `json:"budget"`
				ExpiresAt string  `json:"expires_at"`
			}
			if err2 := json.Unmarshal(data, &tokens); err2 != nil {
				return fmt.Errorf("parsing response: %w", err)
			}
			resp.Tokens = tokens
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tSPENT\tBUDGET\tEXPIRES")
		fmt.Fprintln(w, "──\t────\t──────\t─────\t──────\t───────")
		for _, t := range resp.Tokens {
			status := t.Status
			if status == "revoked" {
				status = "⛔ revoked"
			} else if status == "active" {
				status = "✓ active"
			}
			remaining := ""
			if t.Budget > 0 {
				remaining = fmt.Sprintf("$%.2f", t.Budget)
			} else {
				remaining = "unlimited"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t$%.2f\t%s\t%s\n",
				truncate(t.ID, 12), t.Name, status, t.Spent, remaining, truncate(t.ExpiresAt, 10))
		}
		w.Flush()

		fmt.Fprintf(os.Stderr, "\n%d tokens total\n", len(resp.Tokens))
		return nil
	},
}

var tokenCmd = &cobra.Command{
	Use:   "token [id]",
	Short: "Show token detail: caveats, delegation chain, spend history",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := client.New()
		if err != nil {
			return err
		}

		data, code, err := c.Get("/admin/tokens/" + args[0])
		if err != nil {
			return err
		}
		if code == 404 {
			return fmt.Errorf("token %s not found", args[0])
		}
		if code != 200 {
			return fmt.Errorf("API returned HTTP %d: %s", code, string(data))
		}

		if flagJSON {
			fmt.Println(string(data))
			return nil
		}

		// Pretty-print the token detail
		var resp map[string]interface{}
		json.Unmarshal(data, &resp)

		fmt.Println("Token Detail")
		fmt.Println("─────────────────────────────")
		for _, key := range []string{"id", "name", "status", "spent", "budget", "created_at", "expires_at", "parent_id", "caveats", "routes", "delegation_chain"} {
			if v, ok := resp[key]; ok {
				switch val := v.(type) {
				case []interface{}, map[string]interface{}:
					pretty, _ := json.MarshalIndent(val, "  ", "  ")
					fmt.Printf("  %-18s %s\n", key+":", string(pretty))
				default:
					fmt.Printf("  %-18s %v\n", key+":", val)
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tokensCmd)
	rootCmd.AddCommand(tokenCmd)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
