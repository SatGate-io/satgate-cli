package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	mintAgent    string
	mintBudget   float64
	mintCurrency string
	mintExpiry   string
	mintRoutes   string
	mintParent   string
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint a new capability token for an agent",
	Long: `Mint a new macaroon capability token with budget, expiry, and route restrictions.

Interactive mode (no flags): prompts for all fields.
Non-interactive: provide --agent, --budget, --expiry flags.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		printTarget(cfg)

		// Interactive mode if no agent flag provided
		if mintAgent == "" {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("  Agent name: ")
			mintAgent, _ = reader.ReadString('\n')
			mintAgent = strings.TrimSpace(mintAgent)

			fmt.Print("  Budget (USD, 0 for unlimited): ")
			budgetStr, _ := reader.ReadString('\n')
			budgetStr = strings.TrimSpace(budgetStr)
			if budgetStr != "" && budgetStr != "0" {
				mintBudget, _ = strconv.ParseFloat(budgetStr, 64)
			}

			fmt.Print("  Expiry (e.g. 30d, 24h, or blank for none): ")
			mintExpiry, _ = reader.ReadString('\n')
			mintExpiry = strings.TrimSpace(mintExpiry)

			fmt.Print("  Allowed routes (comma-separated, or * for all): ")
			mintRoutes, _ = reader.ReadString('\n')
			mintRoutes = strings.TrimSpace(mintRoutes)
		}

		if mintAgent == "" {
			return fmt.Errorf("agent name is required")
		}

		c, err := client.New()
		if err != nil {
			return err
		}

		// Build request
		req := map[string]interface{}{
			"name": mintAgent,
		}

		if c.Surface() == "cloud" {
			// Cloud uses credits (cents) and DelegateRequest format
			if mintBudget > 0 {
				req["budget_limit_credits"] = int64(mintBudget * 100) // dollars → credits
			}
			scope := map[string]interface{}{}
			if mintRoutes != "" && mintRoutes != "*" {
				routes := strings.Split(mintRoutes, ",")
				for i := range routes {
					routes[i] = strings.TrimSpace(routes[i])
				}
				scope["routes"] = routes
			} else {
				scope["routes"] = []string{"*"}
			}
			req["scope"] = scope
			// Parent ID: default to root if not specified
			if mintParent != "" {
				req["parent_id"] = mintParent
			}
		} else {
			// Gateway admin API format
			if mintBudget > 0 {
				req["budget"] = mintBudget
				if mintCurrency != "" {
					req["currency"] = mintCurrency
				} else {
					req["currency"] = "USD"
				}
			}
			if mintRoutes != "" && mintRoutes != "*" {
				routes := strings.Split(mintRoutes, ",")
				for i := range routes {
					routes[i] = strings.TrimSpace(routes[i])
				}
				req["routes"] = routes
			}
		}
		if mintExpiry != "" {
			req["expiry"] = mintExpiry
		}

		if flagDry {
			out, _ := json.MarshalIndent(req, "", "  ")
			fmt.Printf("[DRY RUN] Would mint token:\n%s\n", string(out))
			return nil
		}

		// Confirm
		fmt.Fprintf(os.Stderr, "\n  Minting token for agent %q", mintAgent)
		if mintBudget > 0 {
			fmt.Fprintf(os.Stderr, " (budget: $%.2f)", mintBudget)
		}
		if mintExpiry != "" {
			fmt.Fprintf(os.Stderr, " (expires: %s)", mintExpiry)
		}
		fmt.Fprintln(os.Stderr)

		if !confirmAction("⚠️  Proceed?") {
			fmt.Fprintln(os.Stderr, "Cancelled.")
			return nil
		}

		path := "/admin/tokens/mint"
		if c.Surface() == "cloud" {
			path = "/cloud/delegation-v2/delegate"
		}
		data, code, err := c.Post(path, req)
		if err != nil {
			return err
		}
		if code != 200 && code != 201 {
			return fmt.Errorf("API returned HTTP %d: %s", code, string(data))
		}

		if flagJSON {
			fmt.Println(string(data))
			return nil
		}

		var resp map[string]interface{}
		json.Unmarshal(data, &resp)

		fmt.Println("\n✓ Token minted successfully")
		fmt.Println("─────────────────────────────")

		// Cloud response wraps token in {"token": {...}, "macaroon_token": "..."}
		tokenData := resp
		if nested, ok := resp["token"].(map[string]interface{}); ok {
			tokenData = nested
		}

		if id, ok := tokenData["id"]; ok {
			fmt.Printf("  ID:       %v\n", id)
		}
		fmt.Printf("  Agent:    %s\n", mintAgent)
		if status, ok := tokenData["status"]; ok {
			fmt.Printf("  Status:   %v\n", status)
		}
		if mintBudget > 0 {
			fmt.Printf("  Budget:   $%.2f\n", mintBudget)
		}
		if scope, ok := tokenData["scope"].(map[string]interface{}); ok {
			if routes, ok := scope["routes"].([]interface{}); ok {
				routeStrs := make([]string, len(routes))
				for i, r := range routes {
					routeStrs[i] = fmt.Sprintf("%v", r)
				}
				fmt.Printf("  Routes:   %s\n", strings.Join(routeStrs, ", "))
			}
		}
		if exp, ok := tokenData["expires_at"]; ok {
			fmt.Printf("  Expires:  %v\n", exp)
		}
		if mac, ok := resp["macaroon_token"]; ok && mac != "" {
			fmt.Printf("  Macaroon: %v\n", mac)
		} else if mac, ok := resp["macaroon"]; ok && mac != "" {
			fmt.Printf("  Macaroon: %v\n", mac)
		} else if token, ok := resp["token"]; ok {
			if _, isMap := token.(map[string]interface{}); !isMap {
				fmt.Printf("  Token:    %v\n", token)
			}
		}

		fmt.Println("\n⚠️  Save the token/macaroon now — it won't be shown again.")

		return nil
	},
}

func init() {
	mintCmd.Flags().StringVar(&mintAgent, "agent", "", "agent name")
	mintCmd.Flags().Float64Var(&mintBudget, "budget", 0, "budget ceiling in currency units")
	mintCmd.Flags().StringVar(&mintCurrency, "currency", "USD", "budget currency")
	mintCmd.Flags().StringVar(&mintExpiry, "expiry", "", "token expiry (e.g. 30d, 24h)")
	mintCmd.Flags().StringVar(&mintRoutes, "routes", "", "allowed routes (comma-separated)")
	mintCmd.Flags().StringVar(&mintParent, "parent", "", "parent token ID (cloud surface, for delegation)")
	rootCmd.AddCommand(mintCmd)
}
