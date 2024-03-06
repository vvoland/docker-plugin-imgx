package app

import (
	"context"
	"fmt"

	"github.com/docker/cli/cli/command"
)

func List(ctx context.Context, dockerCLI command.Cli, args []string) error {
	out := dockerCLI.Out()
	fmt.Fprintln(out, "TODO")
	return nil
}
