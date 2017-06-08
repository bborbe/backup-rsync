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
	backupSourceDirectory     model.BackupSourceDirectory
	backupSourceBaseDirectory model.BackupSourceBaseDirectory
	remoteHost                model.RemoteHost
	remotePort                model.RemotePort
	remoteUser                model.RemoteUser
	privatePath               model.PrivatePath
	privateKey                model.PrivateKey
	remoteDirectory           model.RemoteTargetDirectory
	today                     time.Time
}

func New(
	backupSourceDirectory model.BackupSourceDirectory,
	backupSourceBaseDirectory model.BackupSourceBaseDirectory,
	remoteHost model.RemoteHost,
	remotePort model.RemotePort,
	remoteUser model.RemoteUser,
	privatePath model.PrivatePath,
	privateKey model.PrivateKey,
	remoteDirectory model.RemoteTargetDirectory,
	today time.Time,
) *backupArchiver {
	b := new(backupArchiver)
	b.backupSourceDirectory = backupSourceDirectory
	b.backupSourceBaseDirectory = backupSourceBaseDirectory
	b.remoteHost = remoteHost
	b.remotePort = remotePort
	b.remoteUser = remoteUser
	b.privatePath = privatePath
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
	exists, err := b.backupExists(ctx)
	if err != nil {
		glog.V(1).Infof("validate failed: %v", err)
		return err
	}
	if exists {
		glog.V(2).Infof("backup already exists")
		return nil
	}
	if err := b.createIncompleteIfNotExists(ctx); err != nil {
		glog.V(1).Infof("create incomplete directory failed: %v", err)
		return err
	}
	if err := b.createCurrentIfNotExists(ctx); err != nil {
		glog.V(1).Infof("create current directory failed: %v", err)
		return err
	}
	if err := b.rsync(ctx); err != nil {
		glog.V(1).Infof("run rsync failed: %v", err)
		return err
	}
	if err := b.renameIncomplete(ctx); err != nil {
		return fmt.Errorf("rename incomplete failed: %v", err)
	}
	if err := b.updateCurrentSymlink(ctx); err != nil {
		return fmt.Errorf("update current symlink failed: %v", err)
	}
	if err := b.remoteEmpty(ctx); err != nil {
		return fmt.Errorf("remove empty failed: %v", err)
	}
	return nil
}

func (b *backupArchiver) backupName() string {
	return b.today.Format("2006-01-02")
}

func (b *backupArchiver) remoteBackupPath() string {
	return b.remoteDirectory.Join(b.backupName())
}

func (b *backupArchiver) remoteIncompletePath() string {
	return b.remoteDirectory.Join("incomplete")
}

func (b *backupArchiver) remoteEmptyPath() string {
	return b.remoteDirectory.Join("empty")
}

func (b *backupArchiver) remoteCurrentPath() string {
	return b.remoteDirectory.Join("current")
}

func (b *backupArchiver) backupExists(ctx context.Context) (bool, error) {
	dir := b.remoteBackupPath()
	glog.V(4).Infof("check if directory %s exists", dir)
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("cd %s", dir)); err != nil {
		glog.V(4).Infof("directory %s does not exists", dir)
		return false, nil
	}
	glog.V(4).Infof("directory %s does exists", dir)
	return true, nil
}

func (b *backupArchiver) createIncompleteIfNotExists(ctx context.Context) error {
	dir := b.remoteIncompletePath() + b.backupSourceDirectory.String()
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("mkdir -p %s", dir)); err != nil {
		glog.V(4).Infof("create directory %s failed: %v", dir, err)
		return err
	}
	glog.V(4).Infof("create directory %s created", dir)
	return nil
}

func (b *backupArchiver) renameIncomplete(ctx context.Context) error {
	glog.V(4).Infof("rename incomplete to date")
	_, err := b.remoteSudo(ctx, fmt.Sprintf("mv %s %s", b.remoteIncompletePath(), b.remoteBackupPath()))
	if err != nil {
		glog.V(4).Infof("rename incomplete to date failed: %v", err)
		return err
	}
	glog.V(4).Infof("rename incomplete to date completed")
	return nil
}

func (b *backupArchiver) updateCurrentSymlink(ctx context.Context) error {
	glog.V(4).Infof("update current symlink")
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("rm %s", b.remoteCurrentPath())); err != nil {
		glog.V(4).Infof("rename incomplete to date failed: %v", err)
		return err
	}
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("ln -s %s %s", b.backupName(), b.remoteCurrentPath())); err != nil {
		glog.V(4).Infof("link backup to current failed: %v", err)
		return err
	}
	return nil
}

func (b *backupArchiver) remoteCurrentExists(ctx context.Context) (bool, error) {
	dir := b.remoteCurrentPath()
	glog.V(4).Infof("check if directory %s exists", dir)
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("cd %s", dir)); err != nil {
		glog.V(4).Infof("directory %s does not exists", dir)
		return false, nil
	}
	glog.V(4).Infof("directory %s does exists", dir)
	return true, nil
}

func (b *backupArchiver) createCurrentIfNotExists(ctx context.Context) error {
	glog.V(4).Infof("create current if not exists started")
	exists, err := b.remoteCurrentExists(ctx)
	if err != nil {
		glog.V(4).Infof("check current exists failed: %v", err)
		return err
	}
	if exists {
		glog.V(4).Infof("current already exists")
		return nil
	}
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("mkdir -p %s", b.remoteEmptyPath())); err != nil {
		glog.V(4).Infof("create directory %s failed: %v", b.remoteEmptyPath(), err)
		return err
	}
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("ln -s empty %s", b.remoteCurrentPath())); err != nil {
		glog.V(4).Infof("link empty to current failed: %v", err)
		return err
	}
	glog.V(4).Infof("create current if not exists finished")
	return nil
}

func (b *backupArchiver) createRemoteExecutor() remote.CommandExecutor {
	return remote.NewCommandExecutor(b.remoteUser, b.remoteHost, b.remotePort, b.privateKey)
}

func (b *backupArchiver) remoteEmpty(ctx context.Context) error {
	glog.V(4).Infof("remove empty directory")
	if _, err := b.remoteSudo(ctx, fmt.Sprintf("rmdir %s", b.remoteEmptyPath())); err != nil {
		glog.V(4).Infof("remove empty dir failed: %v", err)
		return nil
	}
	return nil
}

func (b *backupArchiver) remoteSudo(ctx context.Context, cmd string) (string, error) {
	glog.V(4).Infof("sudo exec '%s'", cmd)
	content, err := b.createRemoteExecutor().ExecuteCommand(ctx, fmt.Sprintf("sudo %s", cmd))
	if err != nil {
		glog.V(4).Infof("sudo exec '%s' failed: %v", cmd, err)
	}
	return content, err
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
		"--rsync-path",
		"sudo rsync",
		"-a",
		"--progress",
		"--compress",
		"--numeric-ids",
		"-e",
		fmt.Sprintf("ssh -T -x -o StrictHostKeyChecking=no -p %d -i %s", b.remotePort, b.privatePath.String()),
		"--delete",
		"--delete-excluded",
		fmt.Sprintf("--port=%d", b.remotePort),
		fmt.Sprintf("--link-dest=%s%s", b.remoteCurrentPath(), b.backupSourceDirectory.String()),
		fmt.Sprintf("%s%s", b.backupSourceBaseDirectory, b.backupSourceDirectory.String()),
		fmt.Sprintf("%s@%s:%s", b.remoteUser, b.remoteHost, b.remoteIncompletePath()+b.backupSourceDirectory.String()),
	)
	return rsyncCommand.Run(ctx)
}
