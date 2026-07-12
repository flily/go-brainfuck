package config

import (
	"testing"
)

func TestConfigureIsStandard(t *testing.T) {
	c1 := Configure(0x01020304)
	if c1.IsStandard() {
		t.Fatalf("expected c1 to not be standard")
	}

	c2 := Configure(0x00010203)
	if !c2.IsStandard() {
		t.Fatalf("expected c2 to be standard")
	}
}
