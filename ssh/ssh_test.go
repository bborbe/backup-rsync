package ssh

import (
	"context"
	. "github.com/bborbe/assert"
	"github.com/bborbe/backup_rsync/model"
	"github.com/bborbe/io/util"
	"io/ioutil"
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	ssh := New("", "", 22, nil, "")
	_, err := ssh.Run(context.Background())
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestExecSuccess(t *testing.T) {
	if testing.Short() {
		return
	}
	path, err := util.NormalizePath("~/.ssh/id_rsa")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	privateKey, err := ioutil.ReadFile(path)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	ssh := New(model.RemoteUser(os.Getenv("USER")), "localhost", 22, privateKey, "ls")
	content, err := ssh.Run(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(len(content) > 0, Is(true)); err != nil {
		t.Fatal(err)
	}
}

func TestExecFail(t *testing.T) {
	if testing.Short() {
		return
	}
	path, err := util.NormalizePath("~/.ssh/id_rsa")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	privateKey, err := ioutil.ReadFile(path)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	ssh := New(model.RemoteUser(os.Getenv("USER")), "localhost", 22, privateKey, "cd /foo")
	_, err = ssh.Run(context.Background())
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
