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

func TestRenderLineSupportsEveryPrintableASCIICharacter(t *testing.T) {
	banner := createTestBanner()

	for char := byte(32); char <= 126; char++ {
		t.Run(fmt.Sprintf("%q", char), func(t *testing.T) {
			got, err := renderLine(string(char), banner)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			expectedLines := make([]string, 8)
			for row := range expectedLines {
				expectedLines[row] = fmt.Sprintf("%c%d", char, row)
			}
			expected := strings.Join(expectedLines, "\n") + "\n"

			if got != expected {
				t.Fatalf("expected %q, got %q", expected, got)
			}
		})
	}
}

func TestRenderLineSupportsExclamationAndDoubleQuoteTogether(t *testing.T) {
	banner := createTestBanner()

	got, err := renderLine(`!"`, banner)
	if err != nil {
		t.Fatal(err)
	}

	expectedLines := []string{
		"!0\"0",
		"!1\"1",
		"!2\"2",
		"!3\"3",
		"!4\"4",
		"!5\"5",
		"!6\"6",
		"!7\"7",
	}
	expected := strings.Join(expectedLines, "\n") + "\n"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestRenderLineRejectsUnsupportedCharacters(t *testing.T) {
	banner := createTestBanner()

	for _, input := range []string{"\x00", "\x1f", "\x7f", "é", "😀"} {
		t.Run(fmt.Sprintf("%q", input), func(t *testing.T) {
			got, err := renderLine(input, banner)
			if err == nil {
				t.Fatalf("expected an error, got output %q", got)
			}
			if !strings.HasPrefix(err.Error(), "unsupported character:") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRenderLineRejectsShortBanner(t *testing.T) {
	_, err := renderLine("!", make([]string, 10))
	if err == nil || err.Error() != "invalid banner file" {
		t.Fatalf("expected invalid banner error, got %v", err)
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

func TestBuildASCII(t *testing.T) {
	banner := createTestBanner()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "single printable character", input: "!", want: "!0\n!1\n!2\n!3\n!4\n!5\n!6\n!7\n"},
		{name: "exclamation and double quote", input: `!"`, want: "!0\"0\n!1\"1\n!2\"2\n!3\"3\n!4\"4\n!5\"5\n!6\"6\n!7\"7\n"},
		{name: "leading empty line", input: `\nA`, want: "\nA0\nA1\nA2\nA3\nA4\nA5\nA6\nA7\n"},
		{name: "trailing empty line", input: `A\n`, want: "A0\nA1\nA2\nA3\nA4\nA5\nA6\nA7\n\n"},
		{name: "multiple empty lines", input: `\n\n`, want: "\n\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildASCII(tt.input, banner)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestBuildASCIIRejectsUnsupportedCharacter(t *testing.T) {
	got, err := buildASCII("A\x00B", createTestBanner())
	if err == nil {
		t.Fatalf("expected an error, got %q", got)
	}
}

func TestThinkertoyHelloThereExclamation(t *testing.T) {
	banner, err := loadBanner("thinkertoy.txt")
	if err != nil {
		t.Fatalf("could not load thinkertoy banner: %v", err)
	}

	got, err := buildASCII("Hello There!", banner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if lines := strings.Count(got, "\n"); lines != 8 {
		t.Fatalf("expected 8 output rows, got %d", lines)
	}
	rows := strings.Split(got, "\n")
	if !strings.HasSuffix(strings.TrimSpace(rows[5]), "O") {
		t.Fatalf("expected the exclamation glyph on row 5, got %q", got)
	}
}

func TestNormalizeLineEndings(t *testing.T) {
	got := normalizeLineEndings("first\r\nsecond\r\nthird")
	if got != "first\nsecond\nthird" {
		t.Fatalf("expected normalized line endings, got %q", got)
	}
}

func TestAllLinesEmpty(t *testing.T) {
	for _, tt := range []struct {
		name  string
		lines []string
		want  bool
	}{
		{name: "empty slice", lines: []string{}, want: true},
		{name: "one empty line", lines: []string{""}, want: true},
		{name: "several empty lines", lines: []string{"", "", ""}, want: true},
		{name: "contains text", lines: []string{"", "text", ""}, want: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := allLinesEmpty(tt.lines); got != tt.want {
				t.Fatalf("expected %t, got %t", tt.want, got)
			}
		})
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
		{name: "thinkertoy banner", args: []string{"prog", "hello", "thinkertoy"}, want: "thinkertoy"},
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
