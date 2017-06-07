package archiver

import (
	"context"
	"fmt"
	"github.com/bborbe/backup_rsync/model"
	"github.com/bborbe/backup_rsync/rsync"
	"github.com/golang/glog"
)

type backupArchiver struct {
	backupSourceDirectory model.BackupSourceDirectory
	remoteHost            model.RemoteHost
	remotePort            model.RemotePort
	remoteUser            model.RemoteUser
	linkDest              model.LinkDest
	remoteDirectory       model.RemoteTargetDirectory
}

func New(
	backupSourceDirectory model.BackupSourceDirectory,
	remoteHost model.RemoteHost,
	remotePort model.RemotePort,
	remoteUser model.RemoteUser,
	linkDest model.LinkDest,
	remoteDirectory model.RemoteTargetDirectory,
) *backupArchiver {
	b := new(backupArchiver)
	b.backupSourceDirectory = backupSourceDirectory
	b.remoteHost = remoteHost
	b.remotePort = remotePort
	b.remoteUser = remoteUser
	b.linkDest = linkDest
	b.remoteDirectory = remoteDirectory
	return b
}

func (b *backupArchiver) Run(ctx context.Context) error {
	glog.V(1).Info("archiv started")
	defer glog.V(1).Info("archiv finished")

	if err := b.runRsync(ctx); err != nil {
		glog.V(1).Infof("run rsync failed: %v", err)
		return err
	}

	return nil
}

func (b *backupArchiver) runRsync(ctx context.Context) error {
	rsyncCommand := rsync.New(
		"-azP",
		"--no-p",
		"--numeric-ids",
		"-e",
		fmt.Sprintf("ssh -T -x -o StrictHostKeyChecking=no -p %d", b.remotePort),
		"--delete",
		"--delete-excluded",
		fmt.Sprintf("--port=%d", b.remotePort),
		fmt.Sprintf("--link-dest=%s", b.linkDest),
		fmt.Sprintf("%s", b.backupSourceDirectory),
		fmt.Sprintf("ssh://%s@%s:%d/%s", b.remoteUser, b.remoteHost, b.remotePort, b.remoteDirectory),
	)

	return rsyncCommand.Run(ctx)
}
