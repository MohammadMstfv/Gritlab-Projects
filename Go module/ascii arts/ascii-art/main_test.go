package main

import (
	"fmt"
	"strings"
	"testing"
)

func createTestBanner() []string {
	banner := make([]string, 1+95*9)

	for char := 32; char <= 126; char++ {
		start := 1 + (char-32)*9

		for row := 0; row < 8; row++ {
			banner[start+row] = fmt.Sprintf("%c%d", rune(char), row)
		}
	}

	return banner
}

func TestRenderLine(t *testing.T) {
	banner := createTestBanner()

	got, err := renderLine("AB", banner)
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := []string{
		"A0B0",
		"A1B1",
		"A2B2",
		"A3B3",
		"A4B4",
		"A5B5",
		"A6B6",
		"A7B7",
	}

	expected := strings.Join(expectedLines, "\n") + "\n"

	if got != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestEmptyInput(t *testing.T) {
	banner := createTestBanner()

	got, err := buildASCII("", banner)
	if err != nil {
		t.Fatal(err)
	}

	if got != "" {
		t.Errorf("expected empty output, got %q", got)
	}
}

func TestNewLine(t *testing.T) {
	banner := createTestBanner()

	got, err := buildASCII(`\n`, banner)
	if err != nil {
		t.Fatal(err)
	}

	if got != "\n" {
		t.Errorf("expected one newline, got %q", got)
	}
}

func TestParseBannerName(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{name: "default banner", args: []string{"prog", "hello"}, want: "standard"},
		{name: "custom banner", args: []string{"prog", "hello", "shadow"}, want: "shadow"},
		{name: "invalid banner", args: []string{"prog", "hello", "invalid"}, wantErr: "Banner must be: standard, shadow, or thinkertoy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBannerName(tt.args)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
