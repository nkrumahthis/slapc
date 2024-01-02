package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// from https://medium.com/@marcus.murray/go-ssh-client-shell-session-c4d40daa46cd

func ConnectToServer(host, user, pwd string) {
	pKey, err := os.ReadFile("config/privatekey")
	
	if err != nil {
		panic("Couldn't read config file")
	}

	var signer ssh.Signer

	signer, err = ssh.ParsePrivateKey(pKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	var hostkeyCallback ssh.HostKeyCallback
	hostkeyCallback, err = knownhosts.New("/Users/zp-tch-024/.ssh/known_hosts")
	if err != nil {
		fmt.Println(err.Error())
	}

	conf := &ssh.ClientConfig{
		User: user,
		HostKeyCallback: hostkeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
			ssh.PublicKeys(signer),
		},
	}

	var conn *ssh.Client

	fmt.Println("Connecting to SSH server")
	conn, err = ssh.Dial("tcp", host, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()
}
