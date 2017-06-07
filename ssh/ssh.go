package ssh

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bborbe/backup_rsync/model"
	"golang.org/x/crypto/ssh"
	"time"
)

type backupSsh struct {
	user       model.RemoteUser
	addr       model.RemoteHost
	port       model.RemotePort
	privateKey model.PrivateKey
	cmd        string
}

func New(
	user model.RemoteUser,
	addr model.RemoteHost,
	port model.RemotePort,
	privateKey model.PrivateKey,
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
	key, err := ssh.ParsePrivateKey(b.privateKey)
	if err != nil {
		return "", fmt.Errorf("parse private key failed: %v", err)
	}
	config := &ssh.ClientConfig{
		User: b.user.String(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		Timeout: 5 * time.Second,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", b.addr, b.port), config)
	if err != nil {
		return "", fmt.Errorf("ssh connect failed: %v", err)
	}
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("create ssh session failed: %v", err)
	}
	defer session.Close()
	var buffer bytes.Buffer
	session.Stdout = &buffer
	//      session.Stdin = bytes.NewBufferString("My input")
	err = session.Run(b.cmd)
	return buffer.String(), err
}
