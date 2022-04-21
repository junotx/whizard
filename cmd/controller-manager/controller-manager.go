package main

import (
	"os"

	"github.com/kubesphere/paodin-monitoring/cmd/controller-manager/app"
)

func main() {
	command := app.NewControllerManagerCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
