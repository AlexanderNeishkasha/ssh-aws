package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Region         string
	AccessKey      string
	SecretKey      string
	PathToStageKey string
	PathToProdKey  string
}

const ConfigName = ".ssh-aws-config"
const ConfigPath = "$HOME/"

func getConfig() Config {
	var config Config
	viper.SetConfigName(ConfigName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		viper.Set("Region", inputField("aws region"))
		viper.Set("AccessKey", inputField("access key"))
		viper.Set("SecretKey", inputField("secret key"))
		viper.Set("PathToStageKey", inputField("path to stg key"))
		viper.Set("PathToProdKey", inputField("path to prod key"))
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatal(err)
		}
	}
	viper.Unmarshal(&config)
	return config
}

func inputField(fieldName string) string {
	var result string
	fmt.Printf("Enter " + fieldName + ": ")
	fmt.Scanln(&result)
	return result
}
