package rsync

import (
	"os"
	"os/exec"
	"strings"

	"syscall"

	"context"

	"github.com/golang/glog"
)

type backupRsync struct {
	args []string
}

func New(args ...string) *backupRsync {
	b := new(backupRsync)
	b.args = args
	return b
}

func (b *backupRsync) Run(ctx context.Context) error {
	glog.V(2).Infof("run: rsync %s", strings.Join(b.args, " "))
	cmd := exec.Command("rsync", b.args...)
	if glog.V(4) {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	glog.V(2).Infof("rsync started")
	if err := cmd.Wait(); err != nil {
		if msg, ok := err.(*exec.ExitError); ok {
			glog.V(2).Infof("rsync closed with exit error")
			if waitstatus, ok := msg.Sys().(syscall.WaitStatus); ok {
				glog.V(2).Infof("rsync closed with exit error: %d", waitstatus.ExitStatus())
				if waitstatus.ExitStatus() == 24 {
					glog.V(2).Infof("rsync finished with vanished file error")
					return nil
				}
			}
		}
		return err
	}
	glog.V(2).Infof("rsync finished")
	return nil
}
