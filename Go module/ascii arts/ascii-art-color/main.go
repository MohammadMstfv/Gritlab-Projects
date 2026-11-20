package main

import (
	"fmt"
	"os"
	"strings"
)

const resetColor = "\033[0m"

var colors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"purple":  "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"gray":    "\033[90m",
	"grey":    "\033[90m",
	"orange":  "\033[38;5;208m",
	"pink":    "\033[38;5;205m",
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return
	}

	bannerName := "standard"

	if isBanner(args[len(args)-1]) && len(args) > 1 { // Check if the last argument is a banner name and there are other arguments
		bannerName = args[len(args)-1]
		args = args[:len(args)-1]
	}

	banner, err := loadBanner(bannerName + ".txt")
	if err != nil {
		fmt.Println("Error reading banner:", err)
		return
	}

	colorCode := ""
	substring := ""
	text := ""

	if strings.HasPrefix(args[0], "--") { // Check if the first argument is a color option
		if !strings.HasPrefix(args[0], "--color=") {
			printUsage()
			return
		}

		colorName := strings.TrimPrefix(args[0], "--color=")

		var ok bool
		colorCode, ok = colors[strings.ToLower(colorName)] // Check if the color name is valid

		if !ok {
			printUsage()
			return
		}

		if len(args) == 2 {
			text = args[1]
		} else if len(args) == 3 {
			substring = args[1]
			text = args[2]
		} else {
			printUsage()
			return
		}
	} else {
		if len(args) != 1 {
			printUsage()
			return
		}

		text = args[0]
	}

	output, err := buildASCII(text, banner, colorCode, substring)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(output)
}

func printUsage() { // Print usage instructions and error message for invalid
	fmt.Println("Usage: go run . [OPTION] [STRING]")
	fmt.Println(`EX: go run . --color=<color> <substring to be colored> "something"`)
}

func isBanner(name string) bool {
	return name == "standard" ||
		name == "shadow" ||
		name == "thinkertoy"
}

func loadBanner(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	text := strings.ReplaceAll(string(data), "\r\n", "\n")

	return strings.Split(text, "\n"), nil
}

func buildASCII(input string, banner []string, colorCode string, substring string) (string, error) {
	if input == "" {
		return "", nil
	}

	input = strings.ReplaceAll(input, `\n`, "\n")

	lines := strings.Split(input, "\n")

	allEmpty := true

	for _, line := range lines {
		if line != "" {
			allEmpty = false
			break
		}
	}

	var result strings.Builder

	if allEmpty {
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

		art, err := renderLine(line, banner, colorCode, substring)

		if err != nil {
			return "", err
		}

		result.WriteString(art)
	}

	return result.String(), nil
}

func renderLine(text string, banner []string, colorCode string, substring string) (string, error) {
	var result strings.Builder

	colored := findColoredCharacters(text, substring, colorCode != "")

	for row := 0; row < 8; row++ {
		for i, char := range text {
			if char < 32 || char > 126 {
				return "", fmt.Errorf("unsupported character: %q", char)
			}

			start := 1 + (int(char)-32)*9
			index := start + row

			if index >= len(banner) {
				return "", fmt.Errorf("invalid banner file")
			}

			if colorCode != "" && colored[i] {
				result.WriteString(colorCode)
				result.WriteString(banner[index])
				result.WriteString(resetColor)
			} else {
				result.WriteString(banner[index])
			}
		}

		result.WriteString("\n")
	}

	return result.String(), nil
}

func findColoredCharacters(text string, substring string, useColor bool) []bool {
	colored := make([]bool, len(text))

	if !useColor {
		return colored
	}

	if substring == "" {
		for i := range colored {
			colored[i] = true
		}

		return colored
	}

	start := 0

	for start <= len(text)-len(substring) {
		index := strings.Index(text[start:], substring)

		if index == -1 {
			break
		}

		index += start

		for i := index; i < index+len(substring); i++ {
			colored[i] = true
		}

		start = index + 1
	}

	return colored
}
