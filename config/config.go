package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/kardianos/osext"
	"github.com/kelseyhightower/envconfig"
)

// Config struct holds the current configuration
type Config struct {
	// Server settings
	Address string `yaml:"address" envconfig:"address"`
	Port    int    `yaml:"port" envconfig:"port"`

	// Logger settings
	LogLevel  string `yaml:"loglevel" envconfig:"loglevel"`
	LogFormat string `yaml:"logformat" envconfig:"logformat"`

	// CrawlInterval in seconds
	CrawlInterval int `yaml:"crawl_interval" envconfig:"crawlinterval"`

	// AWS settings
	AWSRegion string `yaml:"aws_region" envconfig:"awsregion"`

	// Service settings
	// - Owner of the service. For example the team running it.
	//   Defaulted to the current user.
	Owner string `yaml:"owner" envconfig:"owner"`
}

// Initialize a new Config
func Initialize() *Config {
	cfg := DefaultConfig()
	ReadConfigFile(cfg, getConfigFilePath())
	ReadEnvironment(cfg)

	return cfg
}

// DefaultConfig returns a Config struct with default values
func DefaultConfig() *Config {
	return &Config{
		Address: "0.0.0.0",
		Port:    7654,

		LogLevel:  "debug",
		LogFormat: "text",

		CrawlInterval: 600,

		AWSRegion: "eu-west-1",

		Owner: os.Getenv("USER"),
	}
}

// getConfigFilePath returns the location of the config file in order of priority:
// 1 ) File in same directory as the executable
// 2 ) Global file in /etc/melkor/melkor.yml
func getConfigFilePath() string {
	path, _ := osext.ExecutableFolder()
	path = fmt.Sprintf("%s/melkor.yml", path)
	if _, err := os.Open(path); err == nil {
		return path
	}
	globalPath := "/etc/melkor/config.yml"
	if _, err := os.Open(globalPath); err == nil {
		return globalPath
	}

	return ""
}

// ReadConfigFile reads the config file and overrides any values net in both it
// and the DefaultConfig
func ReadConfigFile(cfg *Config, path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	configFile, _ := ioutil.ReadAll(file)
	yaml.Unmarshal(configFile, cfg)
}

// ReadEnvironment overrides any configs set with settings from the environment
func ReadEnvironment(cfg *Config) {
	envconfig.Process("MELKOR", cfg)
}
