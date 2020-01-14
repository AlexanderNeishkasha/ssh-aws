package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/creack/pty"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

const RunningCode = 16

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	amazonSession := initSession()
	resp := loadInstances(amazonSession)
	ips := make(map[string]string, 2)
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			env := env(instance)
			if *instance.State.Code == RunningCode && env != "" {
				ips[env] = *instance.PrivateIpAddress
			}
		}
	}
	pathToKey := os.Getenv("STAGE_KEY")
	cmd := exec.Command("ssh", "-i"+pathToKey, "ubuntu@"+ips["stg"])
	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}

func initSession() *session.Session {
	region := os.Getenv("AWS_DEFAULT_REGION")
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

func loadInstances(amazonSession *session.Session) *ec2.DescribeInstancesOutput {
	ec2Client := ec2.New(amazonSession)
	req, resp := ec2Client.DescribeInstancesRequest(&ec2.DescribeInstancesInput{})
	if err := req.Send(); err != nil {
		log.Println("Error sending request")
		log.Fatal(err)
	}
	return resp
}

func env(instance *ec2.Instance) string {
	role, env := false, ""
	for _, tag := range instance.Tags {
		if *tag.Key == "role" && *tag.Value == "PhpServer" {
			role = true
		}
		if *tag.Key == "env" && (*tag.Value == "prod" || *tag.Value == "stg") {
			env = *tag.Value
		}
	}
	if role {
		return env
	}
	return ""
}
