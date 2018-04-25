package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type yamlConfig struct {
	DBConfig     DBCfgTemplate     `yaml:"database"`
	InFluxConfig InFluxTemplate    `yaml:"influx"`
	ServerConfig ServerCfgTemplate `yaml:"serverDetails"`
}

type DBCfgTemplate struct {
	Server        string `yaml:"server"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Port          string `yaml:"port"`
	Schema        string `yaml:"schema"`
	MaxConnection int    `yaml:"connection_max"`
}

type InFluxTemplate struct {
	Server        string `yaml:"server"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Port          string `yaml:"port"`
	Schema        string `yaml:"schema"`
	MaxConnection int    `yaml:"connection_max"`
}

type ServerCfgTemplate struct {
	Port  string `yaml: port`
	Debug bool   `yaml: debug`
}

var (
	Cfg yamlConfig
)

func InitializeConfig(path string) {
	fmt.Println("Config Path ", path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Unable Read YAML File ")
	}
	if err := yaml.Unmarshal(yamlFile, &Cfg); err != nil {
		log.Fatalln("Error Initializing ", err)
	}
	fmt.Println("============", Cfg.DBConfig.Schema)
}

func Get(configName string) interface{} {
	if configName == "database" {
		return Cfg.DBConfig
	} else if configName == "influx" {
		return Cfg.InFluxConfig
	} else if configName == "server" {
		return Cfg.ServerConfig
	}
	return nil
}
