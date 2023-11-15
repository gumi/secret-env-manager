package aws

import (
	"encoding/json"
	"fmt"

	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

func Load(config *model.Config,withQuote bool) ([]string, error) {
	err := error(nil)

	exports := []string{}
	cache := map[string]string{}

	for _, env := range config.Environments {
		if env.Platform != "aws" {
			continue
		}
		cacheKey := fmt.Sprintf("%s:%s:%s:%s", env.Account, env.Service, env.SecretName, env.Version)

		data := ""
		if value, ok := cache[cacheKey]; ok {
			data = value
		} else {
			data, err = AccessSecret(env)
			if err != nil {
				return nil, err
			}
			cache[cacheKey] = data
		}

		// Key が指定されている場合は、json形式と判断し、Keyに対応する値を取得する
		if env.Key != "" {
			secretJson := map[string]string{}
			err = json.Unmarshal([]byte(data), &secretJson)

			if err != nil {
				return nil, err
			}

			value := secretJson[env.Key]
			if value == "" {
				fmt.Printf("Key %s is not found in %s\n", env.Key, env.SecretName)
				return nil, fmt.Errorf("Key %s is not found in %s", env.Key, env.SecretName)
			}

			export := ""
			if withQuote {
				export = fmt.Sprintf("%s='%s'\n", env.Key, value)
			}else{
				export = fmt.Sprintf("%s=%s\n", env.Key, value)
			}
			exports = append(exports, export)
			continue
		}

		export := ""
		if withQuote {
			export = fmt.Sprintf("%s='%s'\n", env.ExportName, data)
		}else{
			export = fmt.Sprintf("%s=%s\n", env.ExportName, data)
		}
		
		exports = append(exports, export)
	}

	return exports, nil
}
