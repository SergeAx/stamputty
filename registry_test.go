package main

import "testing"

func TestDecodePuttySessionName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "decodes common percent escapes",
			input: "My%20Session%2FProd",
			want:  "My Session/Prod",
		},
		{
			name:  "invalid escape is preserved",
			input: "Bad%ZZName",
			want:  "Bad%ZZName",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodePuttySessionName(tt.input)
			if got != tt.want {
				t.Fatalf("decodePuttySessionName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
