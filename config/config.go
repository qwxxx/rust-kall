package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token    string `yaml:"token"`
	AppName  string `yaml:"appname"`
	Password string `yaml:"password"`
	AppKey   string `yaml:"appkey"`
	SqlConn  string `yaml:"sql_conn"`
}

var Cfg *Config

func New() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("read file failed: %v", err)
	}

	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err)
	}
	return c, nil
}
