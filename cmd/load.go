package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/gcp"
	"github.com/urfave/cli/v2"
)

func Load(c *cli.Context) error {
	switch c.String("type") {

	case "toml":
		loadHandle(file.TOML_FILE_NAME)

	default:
		loadHandle(file.PLAIN_FILE_NAME)

	}

	return nil
}

func loadHandle(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf("%s does not exist.\n", fileName)
		fmt.Println("Please run `sem init`.")
		return
	}

	config, err := file.ReadPlainFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}

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

	fmt.Println(strings.Join(exports, ""))
}
