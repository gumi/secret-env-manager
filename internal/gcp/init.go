package gcp

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/ui"
)

func Init(config *file.Config) error {
	project := ui.TextField("Please Select GCP Project (emplty is skip)", os.Getenv("GCP_PROJECT"))
	if project == "" {
		fmt.Println("Project is empty, skipped GCP init.")
		return nil
	}
	fmt.Printf("project : %s\n", project)

	// get
	secrets, err := ListSecrets(project)
	if err != nil {
		return err
	}

	uiModel := ui.CheckBoxList(*secrets)

	// convert
	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)

	config.GCP.Project = project
	for i, secret := range uiModel.Secrets.Secrets {
		if uiModel.Selected[i] {
			exportName := strings.ToUpper(secret.Name)
			exportName = re.ReplaceAllString(exportName, "_")

			config.GCP.Environments = append(config.GCP.Environments, file.Env{
				SecretName: secret.Name,
				ExportName: exportName,
				Version:    secret.Version,
			})
		}
	}

	return nil
}
