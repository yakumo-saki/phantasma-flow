package ping

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yakumo-saki/phantasma-flow-cli/client"
	"github.com/yakumo-saki/phantasma-flow-cli/out"
)

func Ping(cmd *cobra.Command, args []string) {
	cli := client.NewClient(":5000")
	_, err := cli.Ping(":5000")
	if err != nil {
		out.Msg("Ping failed.")
		out.Msg(err)
		os.Exit(16)
	}

	// normal end
	out.Msg("Ping success.")
}
