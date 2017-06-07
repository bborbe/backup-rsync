package archiver

import (
	"context"
	"testing"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	archiver := New("", "", 22, "", "", "")
	if err := AssertThat(archiver.Archiv(context.Background()), NilValue()); err != nil {
		t.Fatal(err)
	}
}
