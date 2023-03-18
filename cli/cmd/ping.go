package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yakumo-saki/phantasma-flow-cli/ping"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping phantasma-flow server.",
	Long: `Ping does not do any authentication.
	Success of ping indicates "server is alive" only.`,
	Run: ping.Ping,
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
