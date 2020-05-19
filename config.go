package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Region         string `yaml:"Region"`
	AccessKey      string `yaml:"AccessKey"`
	SecretKey      string `yaml:"SecretKey"`
	PathToStageKey string `yaml:"PathToStageKey"`
	PathToProdKey  string `yaml:"PathToProdKey"`
}

const configPath = ".ssh-aws/config.yaml"

func NewConfig() *Config {
	config, err := loadConfigFromFile()
	if err != nil {
		config = createConfig()
		config.store()
	}
	return config
}

func (c *Config) store() {
	file := createConfigFile()
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(&c); err != nil {
		log.Fatal(err)
	}
}

func loadConfigFromFile() (config *Config, err error) {
	file, err := os.Open(getConfigPath())
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(file)
	defer func() {
		if parsingError := recover(); parsingError != nil {
			err = errors.New("error while parsing config file")
		}
	}()
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}
	return config, err
}

func createConfig() *Config {
	return &Config{
		Region:         inputField("aws region"),
		AccessKey:      inputField("access key"),
		SecretKey:      inputField("secret key"),
		PathToStageKey: inputField("path to stg key"),
		PathToProdKey:  inputField("path to prod key"),
	}
}

func createConfigFile() *os.File {
	path := getConfigPath()
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func getConfigPath() string {
	homedir, _ := os.UserHomeDir()
	return homedir + "/" + configPath
}

func inputField(fieldName string) string {
	var result string
	fmt.Printf("Enter " + fieldName + ": ")
	if _, err := fmt.Scanln(&result); err != nil {
		log.Fatal(err)
	}
	return result
}
