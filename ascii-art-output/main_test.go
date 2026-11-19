package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func TestBuildASCIIMultipleLines(t *testing.T) {
	banner := createTestBanner()

	got, err := buildASCII("A\\n\\nB", banner)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := renderLine("A", banner)
	if err != nil {
		t.Fatal(err)
	}
	expected += "\n"
	secondLine, err := renderLine("B", banner)
	if err != nil {
		t.Fatal(err)
	}
	expected += secondLine

	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestBuildASCIIErrors(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		banner []string
	}{
		{name: "unsupported character", input: "é", banner: createTestBanner()},
		{name: "invalid banner", input: "A", banner: []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := buildASCII(tt.input, tt.banner); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

func TestAllLinesEmpty(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{name: "all empty", lines: []string{"", ""}, want: true},
		{name: "contains text", lines: []string{"", "A"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allLinesEmpty(tt.lines); got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
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

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    commandOptions
		wantErr bool
	}{
		{
			name: "single string",
			args: []string{"hello"},
			want: commandOptions{input: "hello", bannerName: "standard"},
		},
		{
			name: "output with banner",
			args: []string{"--output=banner.txt", "hello", "shadow"},
			want: commandOptions{outputFile: "banner.txt", input: "hello", bannerName: "shadow"},
		},
		{
			name:    "invalid output format",
			args:    []string{"--output=banner", "hello"},
			wantErr: true,
		},
		{
			name:    "unknown option",
			args:    []string{"--file=banner.txt", "hello"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if tt.wantErr {
				if err == nil || err.Error() != usageMessage {
					t.Fatalf("expected usage error, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestParseArgsInvalidArguments(t *testing.T) {
	tests := [][]string{
		{},
		{"hello", "standard", "extra"},
		{"--output=file.txt"},
		{"--output=.txt", "hello"},
		{"--output=file", "hello"},
		{"-output=file.txt", "hello"},
		{"hello", "invalid"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, ","), func(t *testing.T) {
			if _, err := parseArgs(args); err == nil || err.Error() != usageMessage {
				t.Fatalf("expected usage error, got %v", err)
			}
		})
	}
}

func TestParseArgsOptionalBanners(t *testing.T) {
	for _, banner := range []string{"standard", "shadow", "thinkertoy"} {
		options, err := parseArgs([]string{"hello", banner})
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", banner, err)
		}
		if options.bannerName != banner {
			t.Errorf("expected banner %q, got %q", banner, options.bannerName)
		}
	}
}

func TestIsOutputFlag(t *testing.T) {
	tests := map[string]bool{
		"--output=file.txt":  true,
		"--output=.txt":      false,
		"--output=file":      false,
		"--output=file.txt=": false,
		"-output=file.txt":   false,
		"--output=file.TXT":  false,
	}

	for argument, want := range tests {
		if got := isOutputFlag(argument); got != want {
			t.Errorf("isOutputFlag(%q) = %v, want %v", argument, got, want)
		}
	}
}

func TestNormalizeLineEndings(t *testing.T) {
	input := "first\r\nsecond\r\n"
	if got := normalizeLineEndings(input); got != "first\nsecond\n" {
		t.Errorf("expected normalized line endings, got %q", got)
	}
}

func TestLoadBanner(t *testing.T) {
	directory := t.TempDir()
	filename := directory + "/banner.txt"
	if err := os.WriteFile(filename, []byte("first\r\nsecond"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := loadBanner(filename)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"first", "second"}; strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("expected %v, got %v", want, got)
	}

	if _, err := loadBanner(directory + "/missing.txt"); err == nil {
		t.Fatal("expected missing banner error")
	}
}

func runMain(t *testing.T, args ...string) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	originalArgs := os.Args
	originalStdout := os.Stdout
	os.Args = append([]string{"ascii-art"}, args...)
	os.Stdout = writer
	main()
	writer.Close()
	os.Args = originalArgs
	os.Stdout = originalStdout

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	reader.Close()
	return string(output)
}

func TestMainOutputModes(t *testing.T) {
	stdout := runMain(t, "A")
	if stdout == "" {
		t.Fatal("expected ASCII art on stdout")
	}

	outputFile := filepath.Join(t.TempDir(), "output.txt")
	stdout = runMain(t, "--output="+outputFile, "A", "thinkertoy")
	if stdout != "" {
		t.Errorf("expected no stdout when writing to a file, got %q", stdout)
	}
	fileOutput, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(fileOutput) == 0 {
		t.Fatal("expected ASCII art in output file")
	}
}

func TestMainErrors(t *testing.T) {
	usage := runMain(t, "--output=invalid", "A")
	if usage != usageMessage+"\n" {
		t.Errorf("expected usage message, got %q", usage)
	}

	buildError := runMain(t, "é")
	if !strings.Contains(buildError, "Error: unsupported character") {
		t.Errorf("expected rendering error, got %q", buildError)
	}

	writeErrorPath := filepath.Join(t.TempDir(), "output.txt")
	if err := os.Mkdir(writeErrorPath, 0755); err != nil {
		t.Fatal(err)
	}
	writeError := runMain(t, "--output="+writeErrorPath, "A")
	if !strings.Contains(writeError, "Error writing output:") {
		t.Errorf("expected file-writing error, got %q", writeError)
	}
}

func TestMainMissingBanner(t *testing.T) {
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDirectory); err != nil {
			t.Error(err)
		}
	})

	output := runMain(t, "A")
	if !strings.Contains(output, "Error reading banner:") {
		t.Errorf("expected banner-reading error, got %q", output)
	}
}
