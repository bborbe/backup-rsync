package model

import (
	"fmt"
	"github.com/bborbe/io/util"
	"io/ioutil"
)

type BackupSourceDirectory string

type RemoteHost string

type RemotePort int

type RemoteUser string

func (r RemoteUser) String() string {
	return string(r)
}

type RemoteTargetDirectory string

type PrivateKey []byte

func PrivateKeyFromFile(path string) (PrivateKey, error) {
	path, err := util.NormalizePath(path)
	if err != nil {
		return nil, fmt.Errorf("normalize path '%s' failed: %v", path, err)
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file '%s' failed: %v", path, err)
	}
	return PrivateKey(content), nil
}
