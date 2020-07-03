package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	ConfigTemplate struct {
		RepositorySet RepositorySet `yaml:"repositories"`
		SenderSet     SenderSet     `yaml:"senders"`
	}

	RepositorySet map[string]map[string]interface{}

	SenderSet map[string]map[string]interface{}
)

func loadConfigFile(configPath string) *ConfigTemplate {
	cp := strings.Replace(configPath, "yaml", "yml", 1)
	f, err := os.Open(filepath.Clean(cp))
	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = f.Close() }()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	var ct ConfigTemplate
	err = yaml.Unmarshal(b, &ct)
	if err != nil {
		fmt.Println(err)
	}
	return &ct
}
