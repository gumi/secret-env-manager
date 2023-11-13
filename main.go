package main

import (
	"fmt"
	"os"

	"github.com/gumi-tsd/secret-env-manager/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{}
	app.Name = "secret-env-manager"
	app.Usage = "manage secret environment variables"

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "file",
			Aliases: []string{"f"},
			Usage:   "Load configuration from the specified format file. Available formats are Plain and Toml, and the default is Plain.",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "init",
			Usage:  fmt.Sprintf("Save the credentials stored in GCP Secret Manager as file."),
			Action: cmd.Init,
			Flags:  flags,
		},
		{
			Name:   "load",
			Usage:  fmt.Sprintf("Output a string to read credentials from SecretManager based on the file and export them as environment variables."),
			Action: cmd.Load,
			Flags:  flags,
		},
		{
			Name:   "update",
			Usage:  fmt.Sprintf("Forcefully update the cached information for the load command."),
			Action: cmd.Update,
			Flags:  flags,
		},
	}

	app.Run(os.Args)
}
