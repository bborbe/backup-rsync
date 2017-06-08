package model

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestRemoteTargetDirectoryStringWithOutSlash(t *testing.T) {
	if err := AssertThat(RemoteTargetDirectory("/foo").String(), Is("/foo/")); err != nil {
		t.Fatal(err)
	}
}

func TestRemoteTargetDirectoryStringWithSlash(t *testing.T) {
	if err := AssertThat(RemoteTargetDirectory("/foo/").String(), Is("/foo/")); err != nil {
		t.Fatal(err)
	}
}

func TestBackupSourceDirectoryStringWithOutSlash(t *testing.T) {
	if err := AssertThat(BackupSourceDirectory("/foo").String(), Is("/foo/")); err != nil {
		t.Fatal(err)
	}
}

func TestBackupSourceDirectoryStringWithSlash(t *testing.T) {
	if err := AssertThat(BackupSourceDirectory("/foo/").String(), Is("/foo/")); err != nil {
		t.Fatal(err)
	}
}
