package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
)

const RunningCode = 16

type AwsFacade struct {
	region    string
	accessKey string
	secretKey string
	token     string
	env       string
	session   *session.Session
	instances []*ec2.Instance
}

func (awsFacade *AwsFacade) IP() (ip string, err error) {
	awsFacade.initSession()
	awsFacade.loadInstances()
	for _, instance := range awsFacade.instances {
		env := awsFacade.getEnv(instance)
		if *instance.State.Code == RunningCode && env != "" {
			return *instance.PrivateIpAddress, nil
		}
	}
	return "", errors.New("IP for env '" + awsFacade.env + "' not found")
}

func (awsFacade *AwsFacade) initSession() {
	if awsFacade.session != nil {
		return
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(awsFacade.region),
		Credentials: credentials.NewStaticCredentials(awsFacade.accessKey, awsFacade.secretKey, awsFacade.token),
	})
	if err != nil {
		log.Fatal(err)
	}
	awsFacade.session = awsSession
}

func (awsFacade *AwsFacade) loadInstances() {
	ec2Client := ec2.New(awsFacade.session)
	req, resp := ec2Client.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	if err := req.Send(); err != nil {
		log.Println("Error sending request")
		log.Fatal(err)
	}
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			awsFacade.instances = append(awsFacade.instances, instance)
		}
	}
}

func (awsFacade *AwsFacade) getEnv(instance *ec2.Instance) string {
	role, env := false, ""
	for _, tag := range instance.Tags {
		if *tag.Key == "role" && *tag.Value == "PhpServer" {
			role = true
		}
		if *tag.Key == "env" && *tag.Value == awsFacade.env {
			env = *tag.Value
		}
	}
	if role {
		return env
	}
	return ""
}
