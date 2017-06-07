package main

import (
	"context"
	"runtime"
	"time"

	"github.com/bborbe/backup_rsync/archiver"
	"github.com/bborbe/backup_rsync/model"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

const (
	defaultWait      = time.Minute * 5
	parameterWait    = "wait"
	parameterOneTime = "one-time"
	parameterSource  = "source"
	parameterHost    = "host"
	parameterPort    = "port"
	parameterUser    = "user"
	parameterLink    = "link"
	parameterTarget  = "target"
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

	backupSourceDirectory := model.BackupSourceDirectory(*backupSourceDirectoryPtr)
	remoteHost := model.RemoteHost(*remoteHostPtr)
	remotePort := model.RemotePort(*remotePortPtr)
	remoteUser := model.RemoteUser(*remoteUserPtr)
	linkDest := model.LinkDest(*linkDestPtr)
	remoteTargetDirectory := model.RemoteTargetDirectory(*remoteTargetDirectoryPtr)

	backupArchiver := archiver.New(backupSourceDirectory, remoteHost, remotePort, remoteUser, linkDest, remoteTargetDirectory)
	return backupArchiver.Run(ctx)
}
