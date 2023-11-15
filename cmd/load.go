package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/gcp"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
	"github.com/urfave/cli/v2"
)

func Load(c *cli.Context) error {
	fileName := ""

	switch c.String("file") {
	case "toml":
		fileName = file.TOML_FILE_NAME

	default:
		fileName = file.PLAIN_FILE_NAME
	}

	if err := loadCache(fileName); err == nil {
		return nil
	}

	cache(fileName)

	return nil
}

func getCacheFileName(fileName string) string {
	return fmt.Sprintf(".cache.%s", fileName)
}

func loadCache(fileName string) error {
	cacheFileName := getCacheFileName(fileName)

	data, err := os.ReadFile(cacheFileName)
	if err == nil {
		printExports(strings.Split(string(data), "\n"))
		return nil
	}

	return err
}

func printExports(exports []string) {
	for _, export := range exports {
		if export == "" {
			continue
		}

		fmt.Printf("%s\n", export)
	}
}

func cache(fileName string) {
	config := readConfigFromFile(fileName)
	exports := loadEnvironments(config)

	os.WriteFile(getCacheFileName(fileName), []byte(strings.Join(exports, "")), 0644)
	printExports(exports)
}

func readConfigFromFile(fileName string) *model.Config {
	config := &model.Config{}
	err := error(nil)

	switch fileName {
	case file.TOML_FILE_NAME:
		config, err = file.ReadTomlFile(fileName)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		config, err = file.ReadPlainFile(fileName)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return config
}

func loadEnvironments(config *model.Config) []string {
	exports := []string{}

	gcpExports, err := gcp.Load(config)
	if err != nil {
		log.Fatalln(err)
	}
	exports = append(exports, gcpExports...)

	awsExports, err := aws.Load(config)
	if err != nil {
		log.Fatalln(err)
	}
	exports = append(exports, awsExports...)

	return exports

}
