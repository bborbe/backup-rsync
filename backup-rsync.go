package main

import (
	"context"
	"github.com/bborbe/backup_rsync/archiver"
	"github.com/bborbe/backup_rsync/model"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
	"runtime"
	"time"
)

const (
	defaultWait             = time.Minute * 5
	parameterWait           = "wait"
	parameterOneTime        = "one-time"
	parameterSource         = "source"
	parameterHost           = "host"
	parameterPort           = "port"
	parameterUser           = "user"
	parameterTarget         = "target"
	parameterPrivateKeyPath = "privatekey"
	parameterBasedir        = "basedir"
)

var (
	waitPtr                  = flag.Duration(parameterWait, defaultWait, "wait")
	oneTimePtr               = flag.Bool(parameterOneTime, false, "exit after first fetch")
	backupSourceDirectoryPtr = flag.String(parameterSource, "", "directory to backup")
	remoteHostPtr            = flag.String(parameterHost, "", "remote host name")
	remotePortPtr            = flag.Int(parameterPort, 22, "remote ssh port")
	remoteUserPtr            = flag.String(parameterUser, "", "remote user name")
	remoteTargetDirectoryPtr = flag.String(parameterTarget, "", "remote target directory")
	privateKeyPathPtr        = flag.String(parameterPrivateKeyPath, "~/.ssh/id_rsa", "path to private key")
	parameterBasedirPtr      = flag.String(parameterBasedir, "", "backup base directory")
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
	var c cron.Cron
	if *oneTimePtr {
		c = cron.NewOneTimeCron(rsync)
	} else {
		c = cron.NewWaitCron(
			*waitPtr,
			rsync,
		)
	}
	return c.Run(context.Background())
}

func rsync(ctx context.Context) error {
	glog.V(1).Info("backup started")
	defer glog.V(1).Info("backup finished")
	privateKeyPath := model.PrivatePath(*privateKeyPathPtr)
	privateKey, err := privateKeyPath.PrivateKey()
	if err != nil {
		glog.V(4).Infof("read private key failed: %v", err)
		return err
	}
	backupArchiver := archiver.New(
		model.BackupSourceDirectory(*backupSourceDirectoryPtr),
		model.BackupSourceBaseDirectory(*parameterBasedirPtr),
		model.RemoteHost(*remoteHostPtr),
		model.RemotePort(*remotePortPtr),
		model.RemoteUser(*remoteUserPtr),
		privateKeyPath,
		privateKey,
		model.RemoteTargetDirectory(*remoteTargetDirectoryPtr),
		time.Now(),
	)

	return backupArchiver.Run(ctx)
}
