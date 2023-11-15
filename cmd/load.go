package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/googlecloud"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
	"github.com/urfave/cli/v2"
)

func Load(c *cli.Context) error {
	fileName := ""
	withExport := c.Bool("with-export")

	switch c.String("file") {
	case "toml":
		fileName = file.TOML_FILE_NAME

	default:
		fileName = file.PLAIN_FILE_NAME
	}
	
	exports, err := loadCache(fileName)
	if err != nil {
		log.Fatalf("%s\nPlaese run `sem update` or `sem update -q` before `sem load`", err)
	}

	printExports(exports,withExport)

	return nil
}

func getCacheFileName(fileName string) string {
	return fmt.Sprintf(".cache.%s", fileName)
}

func loadCache(fileName string) ([]string, error) {
	cacheFileName := getCacheFileName(fileName)

	data, err := os.ReadFile(cacheFileName)
	if err == nil {
		return strings.Split(string(data), "\n"), nil
	}

	return nil, err
}

func printExports(exports []string,withExport bool) {
	for _, export := range exports {
		if export == "" {
			continue
		}
		if withExport {
			export = fmt.Sprintf("export %s", export)
		}else {
			export = export
		}
		fmt.Printf("%s\n", export)
	}
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

func loadEnvironments(config *model.Config, withQuote bool) []string {
	exports := []string{}

	gcpExports, err := googlecloud.Load(config,withQuote)
	if err != nil {
		log.Fatalln(err)
	}
	exports = append(exports, gcpExports...)

	awsExports, err := aws.Load(config,withQuote)
	if err != nil {
		log.Fatalln(err)
	}
	exports = append(exports, awsExports...)

	return exports

}
