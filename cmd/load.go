package cmd

import (
	"fmt"
	"os"
	"strings"

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

	exports := []string{}

	config, err := file.ReadTomlFile()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, env := range config.GCP.Environments {
		secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", config.GCP.Project, env.SecretName, env.Version)
		data, err := gcp.AccessSecretVersion(secretName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		exports = append(exports, fmt.Sprintf("export %s=%s\n", env.ExportName, *data))
	}

	fmt.Println(strings.Join(exports, ""))
	return nil
}
