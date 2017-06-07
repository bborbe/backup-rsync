package main

import (
	. "github.com/bborbe/assert"
	"testing"
)

func TestCleanup(t *testing.T) {
	err := rsync(nil)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
