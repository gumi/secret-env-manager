package file

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

const PLAIN_FILE_NAME = "env"

func ReadPlainFile(fileName string) (*model.Config, error) {
	var config model.Config

	// ファイルを読み込む
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	// 改行で分割する
	lines := strings.Split(string(data), "\n")
	// 環境変数をマップに格納する
	for _, line := range lines {
		if line == "" || line[0] == '#' {
			continue
		}

		env, error := parseEnv(line)
		if error != nil {
			return nil, error
		}

		config.Environments = append(config.Environments, *env)
	}

	return &config, nil
}

func parseEnv(line string) (*model.Env, error) {
	x := strings.SplitN(line, "=", 2)
	// 環境変数名と値を取得する
	name := x[0]
	uri := x[1]

	re := regexp.MustCompile(`sem://(?P<platform>[^/:]+):(?P<service>[^/]+)/(?P<account>[^/]+)/(?P<secretName>[^?]+)(\?version=(?P<version>[^&]+))?`)

	if !re.MatchString(uri) {
		fmt.Println("Invalid URI")
		return nil, nil
	}

	// 正規表現にマッチするグループを取得する
	matches := re.FindStringSubmatch(uri)
	groups := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = matches[i]
		}
	}

	if groups["version"] == "" {
		switch groups["platform"] {
		case "aws":
			groups["version"] = "AWSCURRENT"
		case "gcp":
			groups["version"] = "latest"
		default:
			fmt.Println("Invalid URI")
			return nil, nil
		}
	}

	env := model.Env{
		Platform:   groups["platform"],
		Service:    groups["service"],
		Account:    groups["account"],
		SecretName: groups["secretName"],
		ExportName: name,
		Version:    groups["version"],
	}

	return &env, nil
}

func WritePlainFile(config *model.Config, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	lines := []string{}
	for _, env := range config.Environments {
		exportName := env.ExportName
		if exportName == "" {
			exportName = strings.ToUpper(env.SecretName)
			exportName = re.ReplaceAllString(exportName, "_")
		}

		uri := fmt.Sprintf("sem://%s:%s/%s/%s", env.Platform, env.Service, env.Account, env.SecretName)
		if (env.Platform == "aws" && env.Version != "AWSCURRENT") || (env.Platform == "gcp" && env.Version != "latest") {
			uri = fmt.Sprintf("%s?version=%s", uri, env.Version)
		}

		line := fmt.Sprintf("%s=%s", exportName, uri)
		lines = append(lines, line)
	}

	_, err = f.WriteString(strings.Join(lines, "\n"))
	if err != nil {
		return err
	}

	return nil
}

func IsWrite(file_name string) bool {
	if _, err := os.Stat(file_name); err == nil {
		fmt.Printf("File `%s` already exists. Do you want to overwrite it? (yes/no): ", file_name)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		if input != "yes" {
			fmt.Println("Canceled.")
			return false
		}
	}
	// file not exists
	return true
}
