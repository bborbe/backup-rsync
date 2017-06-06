package main

import (
	"runtime"
	flag "github.com/bborbe/flagenv"
	"github.com/golang/glog"
	"time"
	"context"
	"github.com/bborbe/cron"
)

const (
	defaultWait = time.Minute * 5
	parameterWait = "wait"
	parameterOneTime = "one-time"
)

var (
	waitPtr = flag.Duration(parameterWait, defaultWait, "wait")
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
		cleanup,
	)
	return cron.Run(context.Background())
}

func cleanup(ctx context.Context) error {
	glog.V(1).Info("backup started")
	defer glog.V(1).Info("backup finished")
	return nil
}
