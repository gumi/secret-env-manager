package cmd

import (
	"fmt"
	"log"

	"github.com/gumi-tsd/secret-env-manager/internal/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/gcp"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) error {
	if !file.IsWrite() {
		return nil
	}

	// convert
	config := file.Config{}

	println("------------------------------------------")
	if err := gcp.Init(&config); err != nil {
		log.Fatalln(err)
	}
	println("==========================================")
	if err := aws.Init(&config); err != nil {
		log.Fatalln(err)
	}
	println("------------------------------------------")

	if len(config.AWS.Environments) == 0 && len(config.GCP.Environments) == 0 {
		fmt.Println("No environments found, init canceled.")
		return nil
	}

	// output
	if err := file.WriteTomlFile(config); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(fmt.Sprintf("%s has been saved.", file.FILE_NAME))
	fmt.Println("ExportName can be changed to any value you like.")

	return nil

}
