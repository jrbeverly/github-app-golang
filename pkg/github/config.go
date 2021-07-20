package igithub

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type GithubConfig struct {
	GithubApp GithubConfigApplication `yaml:"github"`
}

type GithubConfigApplication struct {
	GithubAppIdentifier int64  `yaml:"github-app-identifier"`
	GithubPrivateKey    string `yaml:"github-private-key"`
	GithubWebhookSecret string `yaml:"github-webhook-secret"`
}

func ReadGithubConfig(filePath string) (GithubConfig, error) {
	// Can be swapped out with Viper (or something similar)
	var config GithubConfig

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	if err := yaml.UnmarshalStrict(bytes, &config); err != nil {
		return config, err
	}

	return config, nil
}
