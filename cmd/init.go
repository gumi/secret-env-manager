package cmd

import (
	"fmt"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/gcp"
	"github.com/gumi-tsd/secret-env-manager/internal/ui"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) error {
	if !file.IsWrite() {
		return nil
	}

	project := ui.TextField()
	if project == "" {
		fmt.Println("Project is empty.")
		fmt.Println("init canceled.")
		return nil
	}
	fmt.Printf("project: %s\n", project)

	// get
	secrets, err := gcp.ListSecrets(fmt.Sprintf("projects/%s", project))
	if err != nil {
		fmt.Println(err)
		return err
	}

	// select
	uiModel := ui.CheckBoxList(*secrets)

	// convert
	config := file.Config{}
	config.GCP.Project = project
	for i, secret := range uiModel.Secrets.Secrets {
		if uiModel.Selected[i] {
			exportName := strings.ToUpper(secret.Name)
			config.GCP.Environments = append(config.GCP.Environments, file.Env{
				SecretName: secret.Name,
				ExportName: exportName,
				Version:    "latest",
			})
		}
	}

	// output
	file.WriteTomlFile(config)
	fmt.Println(fmt.Sprintf("%s has been saved.", file.FILE_NAME))
	fmt.Println("ExportName can be changed to any value you like.")

	return nil
}
