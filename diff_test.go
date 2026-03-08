package main

import "testing"

func TestBuildSettingsDiffOnlyIncludesOnlyDifferent(t *testing.T) {
	defaultSettings := map[string]interface{}{
		"B": uint32(2),
		"A": "same",
		"C": "default-only",
	}
	defaultTypes := map[string]RegistryType{
		"A": RegString,
		"B": RegDWord,
		"C": RegString,
	}
	sessionSettings := map[string]interface{}{
		"A": "same",
		"B": uint32(3),
	}

	got := buildSettings(defaultSettings, defaultTypes, sessionSettings, false)

	if len(got) != 1 {
		t.Fatalf("len(buildSettings(..., false)) = %d, want 1", len(got))
	}

	if got[0].Name != "B" {
		t.Fatalf("got[0].Name = %q, want %q", got[0].Name, "B")
	}
	if got[0].Type != RegDWord {
		t.Fatalf("got[0].Type = %v, want %v", got[0].Type, RegDWord)
	}
	if !got[0].IsDifferent {
		t.Fatalf("got[0].IsDifferent = false, want true")
	}
}

func TestBuildSettingsShowAllIncludesUnchangedAndMissingAndSorted(t *testing.T) {
	defaultSettings := map[string]interface{}{
		"B": uint32(2),
		"A": "same",
		"C": "default-only",
	}
	defaultTypes := map[string]RegistryType{
		"A": RegString,
		"B": RegDWord,
		"C": RegString,
	}
	sessionSettings := map[string]interface{}{
		"A": "same",
		"B": uint32(3),
	}

	got := buildSettings(defaultSettings, defaultTypes, sessionSettings, true)

	if len(got) != 3 {
		t.Fatalf("len(buildSettings(..., true)) = %d, want 3", len(got))
	}

	if got[0].Name != "A" || got[1].Name != "B" || got[2].Name != "C" {
		t.Fatalf("names = [%s, %s, %s], want [A, B, C]", got[0].Name, got[1].Name, got[2].Name)
	}

	if got[0].IsDifferent {
		t.Fatalf("A IsDifferent = true, want false")
	}
	if !got[1].IsDifferent {
		t.Fatalf("B IsDifferent = false, want true")
	}
	if got[2].IsDifferent {
		t.Fatalf("C IsDifferent = true, want false")
	}

	if got[2].CurrentValue != got[2].DefaultValue {
		t.Fatalf("C CurrentValue = %#v, want default %#v", got[2].CurrentValue, got[2].DefaultValue)
	}

	if got[0].Type != RegString || got[1].Type != RegDWord || got[2].Type != RegString {
		t.Fatalf("unexpected types: got [%v %v %v]", got[0].Type, got[1].Type, got[2].Type)
	}
}
