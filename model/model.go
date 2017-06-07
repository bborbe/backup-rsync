package model

type BackupSourceDirectory string

type RemoteHost string

type RemotePort int

type RemoteUser string

func (r RemoteUser) String() string {
	return string(r)
}

type LinkDest string

type RemoteTargetDirectory string

type PrivateKey []byte
