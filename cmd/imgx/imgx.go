package imgx

import (
	"fmt"

	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func NewRootCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "imgx",
		Short: "docker imgx",
		Long:  "Next generation of docker image management",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			_ = cmd.Help()
			return fmt.Errorf("unknown docker command: imgx %q", args[0])
		},
	}

	cmd.AddCommand(
		newListCommand(dockerCli),
	)

	return cmd
}
