package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type (
	ConfigTemplate struct {
		Repositories []Repository     `yaml:"repos"`
		Slacks       map[string]Slack `yaml:"slacks"`
	}

	Repository struct {
		Name           string `yaml:"name"`
		Owner          string `yaml:"owner"`
		IsEnterprise   bool   `yaml:"is_enterprise"`
		EnterpriseHost string `yaml:"enterprise_host"`
		AccessToken    string `yaml:"access_token"`
		SlackID        string `yaml:"slack_id"`
	}

	Slack struct {
		ID        string   `yaml:"id"`
		Channel   string   `yaml:"channel"`
		Username  string   `yaml:"username"`
		IconEmoji string   `yaml:"icon_emoji"`
		UserMap   []string `yaml:"user_map,omitempty"`
		WebHook   string   `yaml:"web_hook"`
	}
)

func LoadConfigFile() *ConfigTemplate {
	f, err := os.Open(".gocr.yml")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	ct := ConfigTemplate{}
	err = yaml.Unmarshal(b, &ct)
	if err != nil {
		fmt.Println(err)
	}
	return &ct
}
