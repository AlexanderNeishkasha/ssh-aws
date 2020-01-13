package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	amazonSession := initSession()
	ec2Client := ec2.New(amazonSession)
	req, resp := ec2Client.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	if err := req.Send(); err != nil {
		log.Println("Error sending request")
		log.Fatal(err)
	}
	fmt.Println(resp)
}

func initSession() *session.Session {
	region := os.Getenv("AWS_DEFAULT_REGION")
	fmt.Println(region)
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		log.Println("Error with creating aws session")
		log.Fatal(err)
	}
	return sess
}
