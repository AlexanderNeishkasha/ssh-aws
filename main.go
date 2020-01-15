package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	env := os.Args[1]
	awsFacade := AwsFacade{region: os.Getenv("AWS_DEFAULT_REGION"), env: env}
	ip, err := awsFacade.IP()
	if err != nil {
		log.Fatal(err)
	}
	ssh := SshFacade{ip, env}
	ssh.Connect()
}
