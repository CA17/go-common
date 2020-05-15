package conf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type WebConfig struct {
	Debug        bool   `yaml:"debug"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Secret       string `yaml:"secret"`
	CertFile     string `yaml:"cert_file"`
	KeyFile      string `yaml:"key_file"`
	AuthSkip     string `yaml:"auth_skip"`
	AllowOrigins string `yaml:"allow_origins"`
}

type DBConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	MaxConn int    `yaml:"max_conn"`
	MaxIdle int    `yaml:"max_idle"`
	Name    string `yaml:"name"`
	User    string `yaml:"user"`
	Passwd  string `yaml:"passwd"`
}

type AppConfig interface {
	GetWebConfig() *WebConfig
	GetDBConfig() *DBConfig
	GetAppName() string
}

func InitConfig(config AppConfig) error {
	cfgstr, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fmt.Sprintf("/etc/%s.yaml", config.GetAppName()), cfgstr, 664)
}
