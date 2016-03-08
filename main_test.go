package main

import (
	"testing"
)

func TestTrimArguments(t *testing.T) {
	args := []string{"program", "b", "c", "d", "e"}
	trimmed := trimArguments(args)
	if len(trimmed) != 2 {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
	if trimmed[0] != "b" || trimmed[1] != "e" {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
}
