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

func TestFindColoredCharacters(t *testing.T) {
	got := findColoredCharacters("kitten kit", "kit", true)

	expected := []bool{
		true, true, true,
		false, false, false,
		false,
		true, true, true,
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("position %d: expected %v, got %v", i, expected[i], got[i])
		}
	}
}

func TestWholeStringColor(t *testing.T) {
	got := findColoredCharacters("Hello", "", true)

	for i, value := range got {
		if !value {
			t.Errorf("position %d should be colored", i)
		}
	}
}

func TestNoColor(t *testing.T) {
	got := findColoredCharacters("Hello", "", false)

	for i, value := range got {
		if value {
			t.Errorf("position %d should not be colored", i)
		}
	}
}

func TestRenderLineWithColor(t *testing.T) {
	banner := createTestBanner()

	got, err := renderLine("AB", banner, colors["red"], "A")
	if err != nil {
		t.Fatal(err)
	}

	var expected strings.Builder

	for row := 0; row < 8; row++ {
		expected.WriteString(colors["red"])
		expected.WriteString(fmt.Sprintf("A%d", row))
		expected.WriteString(resetColor)

		expected.WriteString(fmt.Sprintf("B%d", row))
		expected.WriteString("\n")
	}

	if got != expected.String() {
		t.Errorf("expected:\n%q\ngot:\n%q", expected.String(), got)
	}
}

func TestEmptyInput(t *testing.T) {
	banner := createTestBanner()

	got, err := buildASCII("", banner, colors["blue"], "")
	if err != nil {
		t.Fatal(err)
	}

	if got != "" {
		t.Errorf("expected empty output, got %q", got)
	}
}
