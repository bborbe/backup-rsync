package rsync

import (
	"testing"

	"context"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	backupRsync := New()
	if err := AssertThat(backupRsync.Run(context.Background()), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
