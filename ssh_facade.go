package main

import (
	"errors"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type SshFacade struct {
	ip       string
	env      string
	stageKey string
	prodKey  string
}

func (ssh *SshFacade) Connect() {
	pathToKey, err := ssh.pathToKey()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("ssh", "-i"+pathToKey, "ubuntu@"+ssh.ip)
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

func (ssh *SshFacade) pathToKey() (key string, err error) {
	if ssh.env == "stg" {
		return ssh.stageKey, nil
	} else if ssh.env == "prod" {
		return ssh.prodKey, nil
	} else {
		return "", errors.New("Key for env '" + ssh.env + "' not found")
	}
}
