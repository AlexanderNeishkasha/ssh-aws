package main

import (
	"errors"
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
	env := os.Args[1]
	awsFacade := AwsFacade{region: os.Getenv("AWS_DEFAULT_REGION"), env: env}
	ip, err := awsFacade.IP()
	if err != nil {
		log.Fatal(err)
	}
	pathToKey, err := pathToKey(env)
	if err != nil {
		log.Fatal(err)
	}
	runSsh(ip, pathToKey)
}

func pathToKey(environment string) (key string, err error) {
	if environment == "stg" {
		return os.Getenv("STAGE_KEY"), nil
	} else if environment == "prod" {
		return os.Getenv("PROD_KEY"), nil
	} else {
		return "", errors.New("Key for env '" + environment + "' not found")
	}
}

func runSsh(ip string, pathToKey string) {
	cmd := exec.Command("ssh", "-i"+pathToKey, "ubuntu@"+ip)
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
