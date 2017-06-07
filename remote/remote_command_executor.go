package remote

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bborbe/backup_rsync/model"
	"golang.org/x/crypto/ssh"
	"time"
)

type CommandExecutor interface {
	ExecuteCommand(ctx context.Context, cmd string) (string, error)
}

type remoteExecutor struct {
	user       model.RemoteUser
	addr       model.RemoteHost
	port       model.RemotePort
	privateKey model.PrivateKey
	timeout    time.Duration
}

func NewCommandExecutor(
	user model.RemoteUser,
	addr model.RemoteHost,
	port model.RemotePort,
	privateKey model.PrivateKey,
) *remoteExecutor {
	b := new(remoteExecutor)
	b.user = user
	b.addr = addr
	b.port = port
	b.privateKey = privateKey
	b.timeout = 5 * time.Second
	return b
}

func (r *remoteExecutor) ExecuteCommand(ctx context.Context, cmd string) (string, error) {
	if len(cmd) == 0 {
		return "", fmt.Errorf("cmd empty")
	}
	key, err := ssh.ParsePrivateKey(r.privateKey)
	if err != nil {
		return "", fmt.Errorf("parse private key failed: %v", err)
	}
	config := &ssh.ClientConfig{
		User: r.user.String(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         r.timeout,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", r.addr, r.port), config)
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
	err = session.Run(cmd)
	return buffer.String(), err
}
