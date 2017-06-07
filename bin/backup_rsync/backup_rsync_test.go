package main

import (
	"context"
	. "github.com/bborbe/assert"
	"testing"
)

func TestCleanup(t *testing.T) {
	err := rsync(context.Background())
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
}
