package main

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func shell(cfg *config) error {

	client, err := connect(cfg)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	if err := session.Shell(); err != nil {
		return err
	}
	session.Wait()
	return nil

}

func connect(cfg *config) (*ssh.Client, error) {
	clientCfg := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{ssh.Password(cfg.Pwd)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", cfg.Addr, clientCfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func tryConnect(cfg *config) bool {
	client, err := connect(cfg)
	if err != nil {
		return false
	}
	defer client.Close()
	return true
}
