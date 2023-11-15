package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/urfave/cli/v2"
)

func Update(c *cli.Context) error {
	fileName := ""
	withQuote := c.Bool("with-quote")

	switch c.String("file") {
	case "toml":
		fileName = file.TOML_FILE_NAME

	default:
		fileName = file.PLAIN_FILE_NAME
	}

	config := readConfigFromFile(fileName)
	exports := loadEnvironments(config,withQuote)

	os.WriteFile(getCacheFileName(fileName), []byte(strings.Join(exports, "")), 0644)
	fmt.Println("cache updated.")

	return nil
}
