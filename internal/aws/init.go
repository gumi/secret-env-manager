package aws

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/model"
	"github.com/gumi-tsd/secret-env-manager/internal/ui"
)

func Init(config *model.Config) error {
	profile := ui.TextField("Please Select AWS Profile (emplty is skip)", os.Getenv("AWS_PROFILE"))
	profile = strings.TrimSpace(profile)

	if profile == "" {
		fmt.Println("Profile is empty, skipped AWS init.")
		return nil
	}
	fmt.Printf("profile : %s\n", profile)

	// get
	secrets, err := ListSecrets(profile)
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
				Platform:   "aws",
				Service:    "secretsmanager",
				Account:    profile,
				SecretName: secret.Name,
				ExportName: exportName,
				Version:    secret.Version,
			})
		}
	}

	return nil
}
