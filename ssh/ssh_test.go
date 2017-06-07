package ssh

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestRun(t *testing.T) {
	err := New("", "", 22, "", "")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
