package main

import (
	"context"
	"runtime"
	"time"

	"github.com/bborbe/backup_rsync/archiver"
	"github.com/bborbe/cron"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
)

const (
	defaultWait      = time.Minute * 5
	parameterWait    = "wait"
	parameterOneTime = "one-time"
)

var (
	waitPtr    = flag.Duration(parameterWait, defaultWait, "wait")
	oneTimePtr = flag.Bool(parameterOneTime, false, "exit after first fetch")
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

	backupArchiver := archiver.New()
	return backupArchiver.Archiv(ctx)
}
