package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestCleanup(t *testing.T) {
	err := cleanup(nil)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
