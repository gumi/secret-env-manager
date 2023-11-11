package aws

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
	"github.com/gumi-tsd/secret-env-manager/internal/ui"
)

func Init(config *file.Config) error {
	profile := ui.TextField("Please Select AWS Profile (emplty is skip)", os.Getenv("AWS_PROFILE"))
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

	config.AWS.Profile = profile
	for i, secret := range uiModel.Secrets.Secrets {
		if uiModel.Selected[i] {
			exportName := strings.ToUpper(secret.Name)
			exportName = re.ReplaceAllString(exportName, "_")

			config.AWS.Environments = append(config.AWS.Environments, file.Env{
				SecretName: secret.Name,
				ExportName: exportName,
				Version:    secret.Version,
			})
		}
	}

	return nil
}
