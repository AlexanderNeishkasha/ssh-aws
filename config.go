package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

type Config struct {
	Region         string
	AccessKey      string
	SecretKey      string
	SessionToken   string
	PathToStageKey string
	PathToProdKey  string
}

const credentialsFile = ".aws/credentials"

func NewConfig() (*Config, error) {
	path := getConfigPath()
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	pathToStageKey := cfg.Section("certs").Key("stg").String()
	if pathToStageKey == "" {
		pathToStageKey = inputField("path to stg key")
		cfg.Section("certs").Key("stg").SetValue(pathToStageKey)
		cfg.SaveTo(path)
	}

	pathToProdKey := cfg.Section("certs").Key("prod").String()
	if pathToProdKey == "" {
		pathToProdKey = inputField("path to prod key")
		cfg.Section("certs").Key("prod").SetValue(pathToProdKey)
		cfg.SaveTo(path)
	}

	return &Config{
		Region:         cfg.Section("default").Key("region").String(),
		AccessKey:      cfg.Section("default").Key("aws_access_key_id").String(),
		SecretKey:      cfg.Section("default").Key("aws_secret_access_key").String(),
		SessionToken:   cfg.Section("default").Key("aws_session_token").String(),
		PathToStageKey: pathToStageKey,
		PathToProdKey:  pathToProdKey,
	}, nil
}

func getConfigPath() string {
	homedir, _ := os.UserHomeDir()
	return homedir + "/" + credentialsFile
}

func inputField(fieldName string) string {
	var result string
	fmt.Printf("Enter " + fieldName + ": ")
	if _, err := fmt.Scanln(&result); err != nil {
		log.Fatal(err)
	}
	return result
}
