package archiver

import (
	"context"
	"testing"

	. "github.com/bborbe/assert"
)

func TestNew(t *testing.T) {
	archiver := New("", "", 22, "", "", "")
	if err := AssertThat(archiver.Run(context.Background()), NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
