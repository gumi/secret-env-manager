package gcp

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

func Load(config *model.Config,withQuote bool) ([]string, error) {
	exports := []string{}

	for _, env := range config.Environments {
		if env.Platform != "gcp" {
			continue
		}

		secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", env.Account, env.SecretName, env.Version)
		data, err := AccessSecretVersion(secretName)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		export := ""
		if withQuote {
			export = fmt.Sprintf("%s='%s'\n", env.ExportName, *data)
		}else{
			export = fmt.Sprintf("%s=%s\n", env.ExportName, *data)
		}
		
		exports = append(exports, export)
	}

	return exports, nil
}
