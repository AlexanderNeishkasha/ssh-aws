package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	env := getEnv()
	fmt.Println("Connecting...")
	awsFacade := AwsFacade{region: config.Region, env: env, accessKey: config.AccessKey, secretKey: config.SecretKey, token: config.SessionToken}
	ip, err := awsFacade.IP()
	if err != nil {
		log.Fatal(err)
	}
	ssh := SshFacade{ip: ip, env: env, prodKey: config.PathToProdKey, stageKey: config.PathToStageKey}
	ssh.Connect()
}

func getEnv() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	var env string
	fmt.Printf("Enter the name of environment (stg/prod): ")
	fmt.Scanln(&env)
	return env
}
