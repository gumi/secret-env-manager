package main

import (
	"fmt"
	"os"

	"github.com/gumi-tsd/secret-env-manager/cmd"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "secret-env-manager"
	app.Usage = "manage secret environment variables"
	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  fmt.Sprintf("Save the credentials stored in GCP Secret Manager as %s.", file.FILE_NAME),
			Action: cmd.Init,
		},
		{
			Name:   "load",
			Usage:  fmt.Sprintf("Output a string to read credentials from SecretManager based on the %s file and export them as environment variables.", file.FILE_NAME),
			Action: cmd.Load,
		},
	}

	app.Run(os.Args)
}
