package main

import (
	"testing"
)

func TestTrimArgumentsWithConfig(t *testing.T) {
	args := []string{"program", "b", "c", "d", "e"}
	trimmed, config := trimArguments(args)
	if len(trimmed) != 2 {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
	if trimmed[0] != "b" || trimmed[1] != "e" {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
	if len(config) != 3 {
		t.Errorf("Bad config arguments returned %+v", config)
	}
	if config[0] != "program" || config[1] != "c" || config[2] != "d" {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
}

func TestTrimArgumentsWithoutConfig(t *testing.T) {
	args := []string{"program", "b", "e"}
	trimmed, config := trimArguments(args)
	if len(trimmed) != 2 {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
	if trimmed[0] != "b" || trimmed[1] != "e" {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
	if len(config) != 1 {
		t.Errorf("Bad config arguments returned %+v", config)
	}
	if config[0] != "program" {
		t.Errorf("Bad arguments returned %+v", trimmed)
	}
}
