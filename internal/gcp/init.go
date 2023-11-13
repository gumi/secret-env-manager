package gcp

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/model"
	"github.com/gumi-tsd/secret-env-manager/internal/ui"
)

func Init(config *model.Config) error {
	project := ui.TextField("Please Select GCP Project (emplty is skip)", os.Getenv("GCP_PROJECT"))
	project = strings.TrimSpace(project)

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

	for i, secret := range uiModel.Secrets.Secrets {
		if uiModel.Selected[i] {
			exportName := strings.ToUpper(secret.Name)
			exportName = re.ReplaceAllString(exportName, "_")

			config.Environments = append(config.Environments, model.Env{
				Platform:   "gcp",
				Service:    "secretmanager",
				Account:    project,
				SecretName: secret.Name,
				ExportName: exportName,
				Version:    secret.Version,
			})
		}
	}

	return nil
}
