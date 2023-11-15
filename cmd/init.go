package cmd

import (
	"fmt"
	"log"

	"github.com/gumi-tsd/secret-env-manager/internal/aws"
	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/googlecloud"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) error {
	switch c.String("file") {

	case "toml":
		fileName := file.TOML_FILE_NAME
		config := initHandle(fileName)
		if err := file.WriteTomlFile(config, fileName); err != nil {
			log.Fatalln(err)
		}
		fmt.Println(fmt.Sprintf("%s has been saved.", fileName))

	default:
		fileName := file.PLAIN_FILE_NAME
		config := initHandle(fileName)
		if err := file.WritePlainFile(config, fileName); err != nil {
			log.Fatalln(err)
		}
		fmt.Println(fmt.Sprintf("%s has been saved.", fileName))

	}

	return nil

}

func initHandle(fileName string) *model.Config {
	if !file.IsWrite(fileName) {
		return nil
	}

	// convert
	config := model.Config{}

	println("------------------------------------------")
	if err := googlecloud.Init(&config); err != nil {
		log.Fatalln(err)
	}
	println("==========================================")
	if err := aws.Init(&config); err != nil {
		log.Fatalln(err)
	}
	println("------------------------------------------")

	if len(config.Environments) == 0 {
		fmt.Println("No environments found, init canceled.")
		return nil
	}

	return &config
}
