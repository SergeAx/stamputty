package main

import (
	"testing"

	"golang.org/x/sys/windows/registry"
)

func TestSettingStringFormatting(t *testing.T) {
	setting := Setting{
		DefaultValue: "abc\x00def",
		CurrentValue: uint32(42),
	}

	if got := setting.GetDefaultValueString(); got != "abcdef" {
		t.Fatalf("GetDefaultValueString() = %q, want %q", got, "abcdef")
	}

	if got := setting.GetCurrentValueString(); got != "42" {
		t.Fatalf("GetCurrentValueString() = %q, want %q", got, "42")
	}
}

func TestRegTypeToRegistryType(t *testing.T) {
	if got := regTypeToRegistryType(registry.DWORD); got != RegDWord {
		t.Fatalf("regTypeToRegistryType(DWORD) = %v, want %v", got, RegDWord)
	}

	if got := regTypeToRegistryType(registry.SZ); got != RegString {
		t.Fatalf("regTypeToRegistryType(SZ) = %v, want %v", got, RegString)
	}
}

