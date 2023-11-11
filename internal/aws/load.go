package aws

import (
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/file"
)

func Load(config *file.Config) ([]string, error) {
	exports := []string{}

	for _, env := range config.AWS.Environments {

		data, err := AccessSecret(env.SecretName)
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

		export := fmt.Sprintf("export %s=\"%s\"\n", env.ExportName, data)
		exports = append(exports, export)
	}

	return exports, nil
}
