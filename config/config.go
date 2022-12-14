package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

var Setting = new(config)

type config struct {
	DefaultFolder string    `yaml:"defaultFolder"`
	AppId         string    `yaml:"appId"`
	AppKey        string    `yaml:"appKey"`
	SecretKey     string    `yaml:"secretKey"`
	AccessToken   string    `yaml:"accessToken"`
	RefreshToken  string    `yaml:"refreshToken"`
	ExpiresIn     int       `yaml:"expiresIn"`
	LastRefresh   time.Time `yaml:"lastRefresh"`
}

func init() {
	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(data, Setting)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
