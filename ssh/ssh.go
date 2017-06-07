package ssh

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type backupSsh struct {
	user       string
	addr       string
	port       int
	privateKey string
	cmd        string
}

func New(
	user string,
	addr string,
	port int,
	privateKey string,
	cmd string,
) *backupSsh {
	b := new(backupSsh)
	b.user = user
	b.addr = addr
	b.port = port
	b.privateKey = privateKey
	b.cmd = cmd
	return b
}

func (b *backupSsh) Run(ctx context.Context) (string, error) {
	key, err := ssh.ParsePrivateKey([]byte(b.privateKey))
	if err != nil {
		return "", err
	}
	config := &ssh.ClientConfig{
		User: b.user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", b.addr, b.port), config)
	if err != nil {
		return "", err
	}
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	var buffer bytes.Buffer
	session.Stdout = &buffer
	//      session.Stdin = bytes.NewBufferString("My input")
	err = session.Run(b.cmd)
	return buffer.String(), err
}
