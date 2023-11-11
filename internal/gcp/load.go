package gcp

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
)

func Load(config *file.Config) ([]string, error) {
	exports := []string{}

	for _, env := range config.GCP.Environments {
		secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", config.GCP.Project, env.SecretName, env.Version)
		data, err := AccessSecretVersion(secretName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		export := fmt.Sprintf("export %s=\"%s\"\n", env.ExportName, *data)

		exports = append(exports, export)
	}

	return exports, nil
}
