package file

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gumi-tsd/secret-env-manager/internal/model"
)

const TOML_FILE_NAME = "env.toml"

func ReadTomlFile(fileName string) (*model.Config, error) {
	var config model.Config

	_, err := toml.DecodeFile(fileName, &config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &config, nil
}

func WriteTomlFile(config *model.Config, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return nil
}
