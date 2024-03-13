package cmd

import (
	"fmt"
	"io/ioutil"
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
	exports := loadEnvironments(config, withQuote)

	os.WriteFile(getCacheFileName(fileName), []byte(strings.Join(exports, "")), 0600)

	setIgnore()

	fmt.Println("cache updated.")

	return nil
}

func setIgnore() {
	// Read the entire .gitignore file
	content, err := ioutil.ReadFile(".gitignore")
	if err != nil {
		panic(err)
	}

	// Check if .cache.env is already in .gitignore
	if !strings.Contains(string(content), ".cache.env") {
		// .cache.env is not in .gitignore, so we add it
		f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err = f.WriteString("\n# auto inserted by sem"); err != nil {
			panic(err)
		}
		if _, err = f.WriteString("\n.cache.env"); err != nil {
			panic(err)
		}
	}
}
