package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func createTestBanner() asciiFont {
	font := asciiFont{height: 8, glyphs: make(map[rune][]string, 95)}

	for char := 32; char <= 126; char++ {
		glyph := make([]string, font.height)
		for row := 0; row < 8; row++ {
			glyph[row] = fmt.Sprintf("%c%d", rune(char), row)
		}
		font.glyphs[rune(char)] = glyph
	}

	return font
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

func TestRenderLineUsesFIGletSmushing(t *testing.T) {
	font := asciiFont{
		height: 1,
		layout: 1,
		glyphs: map[rune][]string{
			'A': {"##"},
			'B': {"#!"},
		},
	}

	got, err := renderLine("AB", font)
	if err != nil {
		t.Fatal(err)
	}
	if got != "##!\n" {
		t.Fatalf("expected smushed glyphs %q, got %q", "##!\n", got)
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

func TestCarriageReturnsActAsNewLines(t *testing.T) {
	font := createTestBanner()
	want, err := buildASCII("A\nB", font)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{name: "escaped carriage return", input: `A\rB`},
		{name: "escaped CRLF", input: `A\r\nB`},
		{name: "actual carriage return", input: "A\rB"},
		{name: "actual CRLF", input: "A\r\nB"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := buildASCII(test.input, font)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("expected line break output %q, got %q", want, got)
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
		{name: "invalid banner", args: []string{"prog", "hello", "invalid"}, wantErr: `unknown banner "invalid"`},
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

func TestLoadFIGletFonts(t *testing.T) {
	for _, name := range []string{"Binary"} {
		t.Run(name, func(t *testing.T) {
			font, err := loadFont(bannerFiles[name])
			if err != nil {
				t.Fatalf("loadFont(%q): %v", name, err)
			}
			if font.height < 1 {
				t.Fatalf("expected positive font height, got %d", font.height)
			}
			output, err := buildASCII("A", font)
			if err != nil {
				t.Fatalf("buildASCII(%q): %v", name, err)
			}
			if output == "" {
				t.Fatal("expected rendered ASCII art")
			}
		})
	}
}

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	homeHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "ASCII Art") {
		t.Fatal("expected home page to contain its title")
	}
}

func TestASCIIArtHandler(t *testing.T) {
	form := url.Values{"text": {"A"}, "banner": {"standard"}}
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	asciiArtHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.Len() == 0 {
		t.Fatal("expected ASCII art output")
	}
}

func TestASCIIArtHandlerAcceptsCarriageReturns(t *testing.T) {
	font, err := loadFont("standard.txt")
	if err != nil {
		t.Fatal(err)
	}
	want, err := buildASCII("A\nB", font)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{name: "escaped carriage return", input: `A\rB`},
		{name: "escaped CRLF", input: `A\r\nB`},
		{name: "actual carriage return", input: "A\rB"},
		{name: "actual CRLF", input: "A\r\nB"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{"text": {test.input}, "banner": {"standard"}}
			req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			recorder := httptest.NewRecorder()

			asciiArtHandler(recorder, req)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
			}
			if recorder.Body.String() != want {
				t.Errorf("expected rendered output %q, got %q", want, recorder.Body.String())
			}
		})
	}
}

func TestASCIIArtHandlerRejectsInvalidBanner(t *testing.T) {
	form := url.Values{"text": {"A"}, "banner": {"invalid"}}
	req := httptest.NewRequest(http.MethodPost, "/ascii-art", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	asciiArtHandler(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}
