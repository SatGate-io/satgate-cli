package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SatGate-io/satgate-cli/internal/client"
	"github.com/SatGate-io/satgate-cli/internal/config"
	"github.com/spf13/cobra"
)

var revokeCmd = &cobra.Command{
	Use:   "revoke [token-id]",
	Short: "Immediately revoke a capability token",
	Long:  `Revoke a token, instantly killing the agent's access. This is irreversible.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		printTarget(cfg)

		tokenID := args[0]

		if flagDry {
			fmt.Fprintf(os.Stderr, "[DRY RUN] Would revoke token %s\n", tokenID)
			return nil
		}

		// Try to get token name for confirmation message
		c, err := client.New()
		if err != nil {
			return err
		}

		// Try to get token name for confirmation
		tokenName := tokenID
		detailPath := "/admin/tokens/" + tokenID
		if c.Surface() == "cloud" {
			detailPath = "/cloud/delegation-v2/token/" + tokenID
		}
		detailData, detailCode, _ := c.Get(detailPath)
		if detailCode == 200 {
			var detail map[string]interface{}
			json.Unmarshal(detailData, &detail)
			if name, ok := detail["name"].(string); ok && name != "" {
				tokenName = fmt.Sprintf("%s (%s)", tokenID, name)
			}
		}

		if !confirmAction(fmt.Sprintf("⚠️  Revoke token %s?\n   This is immediate and irreversible. The agent will lose all access.", tokenName)) {
			fmt.Fprintln(os.Stderr, "Cancelled.")
			return nil
		}

		var data []byte
		var code int
		if c.Surface() == "cloud" {
			data, code, err = c.Post("/cloud/delegation-v2/revoke/"+tokenID, nil)
		} else {
			data, code, err = c.Delete("/admin/tokens/" + tokenID + "/revoke")
		}
		if err != nil {
			return err
		}

		if code == 404 {
			return fmt.Errorf("token %s not found", tokenID)
		}
		if code != 200 && code != 204 {
			return fmt.Errorf("API returned HTTP %d: %s", code, string(data))
		}

		if flagJSON {
			fmt.Println(string(data))
			return nil
		}

		fmt.Fprintf(os.Stderr, "✓ Token %s revoked.\n", tokenName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)
}
