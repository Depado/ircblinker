package conf

// Configuration is the main struct that represents a configuration.
import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Configuration is the configuration structure
type Configuration struct {
	Server      string   `yaml:"server"`
	Channel     string   `yaml:"channel"`
	Nick        string   `yaml:"nick"`
	TLS         bool     `yaml:"tls"`
	InsecureTLS bool     `yaml:"insecure_tls"`
	HL          []string `yaml:"hl"`
	User        string   `yaml:"user"`
}

// C is the exported configuration
var C Configuration

// Load parses the yml file passed as argument and fills the Config.
func Load(fp string) error {
	var content []byte
	var err error
	if content, err = ioutil.ReadFile(fp); err != nil {
		return nil
	}
	return yaml.Unmarshal(content, &C)
}
