package main

import (
	"fmt"
	"os"
	"strings"
)

var validBanners = map[string]struct{}{
	"standard":   {},
	"shadow":     {},
	"thinkertoy": {},
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: go run . [STRING] [BANNER]\n\nEX: go run . something standard")
		return
	}

	bannerName, err := parseBannerName(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	banner, err := loadBanner(bannerName + ".txt")
	if err != nil {
		fmt.Println("Error reading banner:", err)
		return
	}

	output, err := buildASCII(os.Args[1], banner)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(output)
}

func parseBannerName(args []string) (string, error) {
	bannerName := "standard"
	if len(args) == 3 {
		bannerName = args[2]
	}

	if _, ok := validBanners[bannerName]; !ok {
		return "", fmt.Errorf("Banner must be: standard, shadow, or thinkertoy")
	}

	return bannerName, nil
}

func loadBanner(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return strings.Split(normalizeLineEndings(string(data)), "\n"), nil
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func buildASCII(input string, banner []string) (string, error) {
	if input == "" {
		return "", nil
	}

	input = strings.ReplaceAll(input, `\n`, "\n") // Handle escaped newlines
	lines := strings.Split(input, "\n")           // Split input into lines

	var result strings.Builder

	if allLinesEmpty(lines) {
		for i := 0; i < len(lines)-1; i++ {
			result.WriteString("\n")
		}
		return result.String(), nil
	}

	for _, line := range lines {
		if line == "" {
			result.WriteString("\n")
			continue
		}

		art, err := renderLine(line, banner)
		if err != nil {
			return "", err
		}

		result.WriteString(art)
	}

	return result.String(), nil
}

func allLinesEmpty(lines []string) bool {
	for _, line := range lines {
		if line != "" {
			return false
		}
	}

	return true
}

func renderLine(text string, banner []string) (string, error) {
	var result strings.Builder

	for row := 0; row < 8; row++ { // Each character is represented in 8 rows
		for _, char := range []byte(text) {
			if char < 32 || char > 126 {
				return "", fmt.Errorf("unsupported character: %q", char)
			}

			start := 1 + (int(char)-32)*9 // Each character takes up 9 lines in the banner (8 for the art + 1 empty line)
			index := start + row

			if index >= len(banner) {
				return "", fmt.Errorf("invalid banner file")
			}

			result.WriteString(banner[index])
		}

		result.WriteString("\n")
	}

	return result.String(), nil
}
