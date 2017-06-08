package archiver

import (
	"context"
	. "github.com/bborbe/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	archiver := New("", "", "", 22, "", nil, "", time.Now())
	if err := AssertThat(archiver.Run(context.Background()), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestBackupName(t *testing.T) {
	archiver := New("", "", "", 22, "", nil, "", time.Unix(1496840099, 0))
	if err := AssertThat(archiver.backupName(), Is("2017-06-07")); err != nil {
		t.Fatal(err)
	}
}
