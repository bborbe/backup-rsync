package archiver

import (
	"context"
	"fmt"
	"github.com/bborbe/backup_rsync/model"
	"github.com/bborbe/backup_rsync/remote"
	"github.com/bborbe/backup_rsync/rsync"
	"github.com/golang/glog"
	"time"
)

type backupArchiver struct {
	backupSourceDirectory model.BackupSourceDirectory
	remoteHost            model.RemoteHost
	remotePort            model.RemotePort
	remoteUser            model.RemoteUser
	privateKey            model.PrivateKey
	remoteDirectory       model.RemoteTargetDirectory
	today                 time.Time
}

func New(
	backupSourceDirectory model.BackupSourceDirectory,
	remoteHost model.RemoteHost,
	remotePort model.RemotePort,
	remoteUser model.RemoteUser,
	privateKey model.PrivateKey,
	remoteDirectory model.RemoteTargetDirectory,
	today time.Time,
) *backupArchiver {
	b := new(backupArchiver)
	b.backupSourceDirectory = backupSourceDirectory
	b.remoteHost = remoteHost
	b.remotePort = remotePort
	b.remoteUser = remoteUser
	b.privateKey = privateKey
	b.remoteDirectory = remoteDirectory
	b.today = today
	return b
}

func (b *backupArchiver) Run(ctx context.Context) error {
	glog.V(1).Info("archiv started")
	defer glog.V(1).Info("archiv finished")

	if err := b.validate(); err != nil {
		glog.V(1).Infof("validate failed: %v", err)
		return err
	}

	exists, err := b.backupExists()
	if err != nil {
		glog.V(1).Infof("validate failed: %v", err)
		return err
	}
	if exists {
		glog.V(2).Infof("backup already exists")
		return nil
	}
	if err := b.createIncompleteIfNotExists(); err != nil {
		glog.V(1).Infof("create incomplete directory failed: %v", err)
		return err
	}
	if err := b.createCurrentIfNotExists(); err != nil {
		glog.V(1).Infof("create current directory failed: %v", err)
		return err
	}
	if err := b.rsync(ctx); err != nil {
		glog.V(1).Infof("run rsync failed: %v", err)
		return err
	}
	if err := b.renameIncomplete(); err != nil {
		return fmt.Errorf("rename incomplete failed: %v", err)
	}
	if err := b.updateCurrentSymlink(); err != nil {
		return fmt.Errorf("update current symlink failed: %v", err)
	}
	return nil
}

func (b *backupArchiver) backupName() string {
	return b.today.Format("2006-01-02")
}

func (b *backupArchiver) backupExists() (bool, error) {
	return false, nil
}

func (b *backupArchiver) createIncompleteIfNotExists() error {
	return nil
}

func (b *backupArchiver) renameIncomplete() error {
	return nil
}

func (b *backupArchiver) updateCurrentSymlink() error {
	return nil
}

func (b *backupArchiver) createCurrentIfNotExists() error {
	return nil
}

func (b *backupArchiver) createRemoteExecutor() remote.CommandExecutor {
	return remote.NewCommandExecutor(b.remoteUser, b.remoteHost, b.remotePort, b.privateKey)
}

func (b *backupArchiver) validate() error {
	if len(b.backupSourceDirectory) == 0 {
		return fmt.Errorf("backup directory invalid")
	}
	if len(b.remoteHost) == 0 {
		return fmt.Errorf("remote host invalid")
	}
	if b.remotePort <= 0 {
		return fmt.Errorf("remote port invalid")
	}
	if len(b.remoteUser) == 0 {
		return fmt.Errorf("remote user invalid")
	}
	if len(b.remoteDirectory) == 0 {
		return fmt.Errorf("remote directory invalid")
	}
	return nil
}

func (b *backupArchiver) rsync(ctx context.Context) error {
	rsyncCommand := rsync.New(
		"-azP",
		"--no-p",
		"--numeric-ids",
		"-e",
		fmt.Sprintf("ssh -T -x -o StrictHostKeyChecking=no -p %d", b.remotePort),
		"--delete",
		"--delete-excluded",
		fmt.Sprintf("--port=%d", b.remotePort),
		//fmt.Sprintf("--link-dest=%s", "current"),
		fmt.Sprintf("%s", b.backupSourceDirectory),
		fmt.Sprintf("%s@%s:%s", b.remoteUser, b.remoteHost, b.remoteDirectory),
	)
	return rsyncCommand.Run(ctx)
}
