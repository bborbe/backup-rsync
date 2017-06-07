package main

import (
	"context"
	"runtime"
	"time"

	"fmt"
	"github.com/bborbe/backup_rsync/archiver"
	"github.com/bborbe/backup_rsync/model"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

const (
	defaultWait         = time.Minute * 5
	parameterWait       = "wait"
	parameterOneTime    = "one-time"
	parameterSource     = "source"
	parameterHost       = "host"
	parameterPort       = "port"
	parameterUser       = "user"
	parameterLink       = "link"
	parameterTarget     = "target"
	parameterPrivateKey = "privatekey"
)

var (
	waitPtr                  = flag.Duration(parameterWait, defaultWait, "wait")
	oneTimePtr               = flag.Bool(parameterOneTime, false, "exit after first fetch")
	backupSourceDirectoryPtr = flag.String(parameterSource, "", "directory to backup")
	remoteHostPtr            = flag.String(parameterHost, "", "remote host name")
	remotePortPtr            = flag.Int(parameterPort, 22, "remote ssh port")
	remoteUserPtr            = flag.String(parameterUser, "", "remote user name")
	linkDestPtr              = flag.String(parameterLink, "", "link dest")
	remoteTargetDirectoryPtr = flag.String(parameterTarget, "", "remote target directory")
	privateKeyPtr            = flag.String(parameterPrivateKey, "", "private key")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := do(); err != nil {
		glog.Exit(err)
	}
}

func do() error {
	cron := cron.New(
		*oneTimePtr,
		*waitPtr,
		rsync,
	)
	return cron.Run(context.Background())
}

func rsync(ctx context.Context) error {
	glog.V(1).Info("backup started")
	defer glog.V(1).Info("backup finished")

	privateKey, err := model.PrivateKeyFromFile(*privateKeyPtr)
	if err != nil {
		return fmt.Errorf("read private key failed: %v", err)
	}

	backupArchiver := archiver.New(
		model.BackupSourceDirectory(*backupSourceDirectoryPtr),
		model.RemoteHost(*remoteHostPtr),
		model.RemotePort(*remotePortPtr),
		model.RemoteUser(*remoteUserPtr),
		privateKey,
		model.LinkDest(*linkDestPtr),
		model.RemoteTargetDirectory(*remoteTargetDirectoryPtr),
		time.Now(),
	)

	return backupArchiver.Run(ctx)
}
