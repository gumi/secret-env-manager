package cmd

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/urfave/cli/v2"
)

func Update(c *cli.Context) error {
	fileName := ""
	switch c.String("file") {
	case "toml":
		fileName = file.TOML_FILE_NAME

	default:
		fileName = file.PLAIN_FILE_NAME
	}

	cache(fileName)
	fmt.Println("cache updated.")

	return nil
}
