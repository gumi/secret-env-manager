package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

const FILE_NAME = "env.toml"

type Config struct {
	AWS AWS `toml:"aws"`
	GCP GCP `toml:"gcp"`
}

type AWS struct {
	Profile      string `toml:"Profile"`
	Environments []Env  `toml:"env"`
}

type GCP struct {
	Project      string `toml:"Project"`
	Environments []Env  `toml:"env"`
}

type Env struct {
	SecretName string `toml:"SecretName"`
	ExportName string `toml:"ExportName"`
	Version    string `toml:"Version"`
}

func ReadTomlFile() (*Config, error) {
	var config Config

	_, err := toml.DecodeFile(FILE_NAME, &config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &config, nil
}

func WriteTomlFile(config Config) error {
	f, err := os.Create(FILE_NAME)
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

func IsWrite() bool {
	if _, err := os.Stat(FILE_NAME); err == nil {
		fmt.Printf("File %s already exists. Do you want to overwrite it? (yes/no): ", FILE_NAME)
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
