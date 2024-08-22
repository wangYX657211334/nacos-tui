package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type Api interface {
	GetNacosContexts() []NacosContext
	GetNacosContext() NacosContext
	SetNacosContext(name string) error
	SetNacosContextNamespace(namespace string, namespaceName string) error
	GetStringProperty(key string, defaultValue string) string
	GetIntProperty(key string, defaultValue string) int
	GetBoolProperty(key string, defaultValue string) bool
	SetProperty(string, string) error
}

func LoadApplicationConfig() (*ApplicationConfig, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Join(errors.New("get user home dir error"), err)
	}
	configPath := filepath.Join(userHome, ".nacos-tui", "nacos-tui.yaml")
	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		defaultConfig := DefaultApplicationConfig()
		err := SaveApplicationConfig(&defaultConfig)
		if err != nil {
			return nil, err
		}
	}
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.Join(errors.New("application read config error"), err)
	}
	var config ApplicationConfig
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, errors.Join(errors.New("application config unmarshal error"), err)
	}
	if config.Properties == nil {
		config.Properties = make(map[string]string)
	}
	return &config, nil
}

func SaveApplicationConfig(config *ApplicationConfig) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return errors.Join(errors.New("get user home dir error"), err)
	}
	configPath := filepath.Join(userHome, ".nacos-tui", "nacos-tui.yaml")
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return errors.Join(errors.New("application config marshal error"), err)
	}
	err = os.WriteFile(configPath, configYaml, 0666)
	if err != nil {
		return errors.Join(errors.New("application write config error"), err)
	}
	return nil
}

func DefaultApplicationConfig() ApplicationConfig {
	return ApplicationConfig{
		UseServer: "local",
		Servers: []NacosContext{
			{
				Name:             "local",
				Url:              "http://127.0.0.1:8848/nacos",
				User:             "nacos",
				Password:         "nacos",
				UseNamespaceName: "public",
			},
		},
	}
}

type ApplicationConfig struct {
	Properties map[string]string
	UseServer  string
	Servers    []NacosContext
}
type NacosContext struct {
	Name             string
	Url              string
	User             string
	Password         string
	UseNamespace     string
	UseNamespaceName string
}

func (c *ApplicationConfig) GetNacosContexts() []NacosContext {
	return c.Servers
}
func (c *ApplicationConfig) GetStringProperty(key string, defaultValue string) string {
	if len(c.Properties[key]) == 0 {
		return defaultValue
	}
	return c.Properties[key]
}
func (c *ApplicationConfig) GetIntProperty(key string, defaultValue string) int {
	value, err := strconv.Atoi(c.GetStringProperty(key, defaultValue))
	if err != nil {
		value, _ = strconv.Atoi(defaultValue)
		return value
	}
	return value
}
func (c *ApplicationConfig) GetBoolProperty(key string, defaultValue string) bool {
	value, err := strconv.ParseBool(c.GetStringProperty(key, defaultValue))
	if err != nil {
		value, _ = strconv.ParseBool(defaultValue)
		return value
	}
	return value
}
func (c *ApplicationConfig) SetProperty(key string, value string) error {
	c.Properties[key] = value
	return SaveApplicationConfig(c)
}

func (c *ApplicationConfig) GetNacosContext() NacosContext {
	var serverConfig NacosContext
	for _, server := range c.Servers {
		if strings.EqualFold(server.Name, c.UseServer) {
			serverConfig = server
		}
	}
	return serverConfig
}

func (c *ApplicationConfig) SetNacosContext(name string) error {
	c.UseServer = name
	return SaveApplicationConfig(c)
}
func (c *ApplicationConfig) SetNacosContextNamespace(namespace string, namespaceName string) error {
	for i := range c.Servers {
		if strings.EqualFold(c.Servers[i].Name, c.UseServer) {
			c.Servers[i].UseNamespace = namespace
			c.Servers[i].UseNamespaceName = namespaceName
		}
	}
	return SaveApplicationConfig(c)
}
