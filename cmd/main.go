package main

import (
	"os"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/vvoland/docker-plugin-imgx/cmd/imgx"
)

func main() {
	if plugin.RunningStandalone() {
		os.Args = append([]string{"docker"}, os.Args[1:]...)
	}

	plugin.Run(imgx.NewRootCommand, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Paweł Gronowski",
		Version:       "0.1.0",
	})
}
