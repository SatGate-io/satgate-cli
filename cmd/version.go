package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print CLI version and build info",
	Run: func(cmd *cobra.Command, args []string) {
		if flagJSON {
			out, _ := json.MarshalIndent(map[string]string{
				"version":    version,
				"build_time": buildTime,
			}, "", "  ")
			fmt.Println(string(out))
			return
		}
		fmt.Printf("satgate %s (built %s)\n", version, buildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
