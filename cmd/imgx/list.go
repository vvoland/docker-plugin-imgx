package imgx

import (
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"github.com/vvoland/docker-plugin-imgx/app"
)

func newListCommand(dockerCli command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List images",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			err := app.List(ctx, dockerCli, args)
			if err != nil {
				cmd.PrintErrln(err)
			}
		},
	}
}
