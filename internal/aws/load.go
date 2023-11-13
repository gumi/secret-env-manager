package aws

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

func Load(config *model.Config) ([]string, error) {
	exports := []string{}

	for _, env := range config.Environments {
		if env.Platform != "aws" {
			continue
		}

		data, err := AccessSecret(env.SecretName, env.Account)
		if err != nil {
			return nil, err
		}

		// secretJson := map[string]string{}
		// err = json.Unmarshal([]byte(data), &secretJson)
		// if err != nil {
		// 	export := fmt.Sprintf("export %s=\"%s\"\n", env.ExportName, data)
		// 	exports = append(exports, export)
		// 	continue
		// }

		// for k, v := range secretJson {
		// 	exportName := strings.ToUpper(k)
		// 	export := fmt.Sprintf("export %s=\"%s\"\n", exportName, v)
		// 	exports = append(exports, export)
		// }

		export := fmt.Sprintf("export %s='%s'\n", env.ExportName, data)
		exports = append(exports, export)
	}

	return exports, nil
}
