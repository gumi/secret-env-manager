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

	fileFlags := &cli.StringFlag{
		Name:    "file",
		Aliases: []string{"f"},
		Usage:   "Load configuration from the specified format file. Available formats are Plain and Toml, and the default is Plain.",
	}

	exportFlags := &cli.BoolFlag{
		Name:    "with-export",
		Aliases: []string{"e"},
		Usage:   "When this option is enabled, 'export' is added when displaying to standard output.",
	}

	quoteFlags := &cli.BoolFlag{
		Name:    "with-quote",
		Aliases: []string{"q"},
		Usage:   "When this option is enabled, single quotes are added to the cached environment variable values. It is generally recommended to use this option when the value contains spaces.",
	}


	app.Commands = []*cli.Command{
		{
			Name:   "init",
			Usage:  fmt.Sprintf("Save the credentials stored in GCP Secret Manager as file."),
			Action: cmd.Init,
			Flags:  []cli.Flag{
				fileFlags,
				},
		},
		{
			Name:   "load",
			Usage:  fmt.Sprintf("Output a string to read credentials from SecretManager based on the file and export them as environment variables."),
			Action: cmd.Load,
			Flags:  []cli.Flag{
				fileFlags,
				exportFlags,
				},
		},
		{
			Name:   "update",
			Usage:  fmt.Sprintf("Forcefully update the cached information for the load command."),
			Action: cmd.Update,
			Flags:  []cli.Flag{
				fileFlags,
				quoteFlags,
				},
		},
	}

	app.Run(os.Args)
}
