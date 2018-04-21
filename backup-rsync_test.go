package main

import (
	"context"
	. "github.com/bborbe/assert"
	"testing"
)

func TestRsync(t *testing.T) {
	err := rsync(context.Background())
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
