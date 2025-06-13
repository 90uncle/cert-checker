package config

import (
	"cert-checker/internal/notifier"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Domains []string `yaml:"domains"`
	Alert   struct {
		Threshold      int `yaml:"threshold"`
		RepeatInterval int `yaml:"repeat_interval"`
	} `yaml:"alert"`
	Notifiers struct {
		Email    notifier.EmailNotifier    `yaml:"email"`
		DingTalk notifier.DingTalkNotifier `yaml:"dingtalk"`
		WeCom    notifier.WeComNotifier    `yaml:"wecom"`
	} `yaml:"notifiers"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
