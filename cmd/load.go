package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/gcp"
	"github.com/urfave/cli"
)

func Load(c *cli.Context) error {
	// check exist with stat
	if _, err := os.Stat(file.FILE_NAME); os.IsNotExist(err) {
		fmt.Printf("%s does not exist.\n", file.FILE_NAME)
		fmt.Println("Please run `sem init`.")
		return nil
	}

	config, err := file.ReadTomlFile()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	exports := []string{}

	gcpExports, err := gcp.Load(config)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	exports = append(exports, gcpExports...)

	awsExports, err := aws.Load(config)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	exports = append(exports, awsExports...)

	fmt.Println(strings.Join(exports, ""))
	return nil
}
