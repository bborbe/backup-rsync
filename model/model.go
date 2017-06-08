package model

import (
	"fmt"
	"github.com/bborbe/io/util"
	"io/ioutil"
	"strings"
)

type BackupSourceBaseDirectory string

type BackupSourceDirectory string

func (b BackupSourceDirectory) String() string {
	return appendSlash(string(b))
}

type RemoteHost string

type RemotePort int

type RemoteUser string

func (r RemoteUser) String() string {
	return string(r)
}

type RemoteTargetDirectory string

func (r RemoteTargetDirectory) Join(name string) string {
	return fmt.Sprintf("%s%s", r.String(), name)
}

func (r RemoteTargetDirectory) String() string {
	return appendSlash(string(r))
}

func appendSlash(name string) string {
	if strings.HasSuffix(name, "/") {
		return name
	}
	return name + "/"
}

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
