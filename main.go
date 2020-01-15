package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	env := getEnv()
	fmt.Println("Connecting...")
	awsFacade := AwsFacade{region: os.Getenv("AWS_DEFAULT_REGION"), env: env}
	ip, err := awsFacade.IP()
	if err != nil {
		log.Fatal(err)
	}
	ssh := SshFacade{ip, env}
	ssh.Connect()
}

func getEnv() string {
	if len(os.Args) > 1  {
		return os.Args[1]
	}
	var env string
	fmt.Printf("Enter the name of environment (stg/prod): ")
	fmt.Scanln(&env)
	return env
}
