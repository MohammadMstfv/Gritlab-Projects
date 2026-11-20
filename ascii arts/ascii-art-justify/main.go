package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

const usageMessage = "Usage: go run . [OPTION] [STRING] [BANNER]\n\nExample: go run . --align=right something standard"

type alignment string

const (
	alignLeft    alignment = "left"
	alignCenter  alignment = "center"
	alignRight   alignment = "right"
	alignJustify alignment = "justify"
)

var validBanners = map[string]struct{}{
	"standard":   {},
	"shadow":     {},
	"thinkertoy": {},
}

func main() {
	text, bannerName, align, err := parseArguments(os.Args[1:])
	if err != nil {
		fmt.Println(usageMessage)
		return
	}

	banner, err := loadBanner(bannerName + ".txt")
	if err != nil {
		fmt.Println("Error reading banner:", err)
		return
	}

	output, err := buildASCIIWithAlignment(text, banner, align, terminalWidth())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(output)
}

func parseArguments(args []string) (string, string, alignment, error) {
	if len(args) < 1 || len(args) > 3 {
		return "", "", "", fmt.Errorf("invalid arguments")
	}

	align := alignLeft
	if strings.HasPrefix(args[0], "--") {
		if !strings.HasPrefix(args[0], "--align=") {
			return "", "", "", fmt.Errorf("invalid alignment option")
		}
		align = alignment(strings.TrimPrefix(args[0], "--align="))
		if !isValidAlignment(align) {
			return "", "", "", fmt.Errorf("invalid alignment")
		}
		args = args[1:]
	}

	if len(args) < 1 || len(args) > 2 {
		return "", "", "", fmt.Errorf("invalid arguments")
	}

	bannerName := "standard"
	if len(args) == 2 {
		bannerName = args[1]
	}
	if _, ok := validBanners[bannerName]; !ok {
		return "", "", "", fmt.Errorf("invalid banner")
	}

	return args[0], bannerName, align, nil
}

func isValidAlignment(value alignment) bool {
	return value == alignLeft || value == alignCenter || value == alignRight || value == alignJustify
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

func normalizeLineEndings(s string) string { // make it same in windows and linux
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func buildASCII(input string, banner []string) (string, error) {
	return buildASCIIWithAlignment(input, banner, alignLeft, 0)
}

func buildASCIIWithAlignment(input string, banner []string, align alignment, width int) (string, error) {
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

		art, err := renderLineAligned(line, banner, align, width)
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
	return renderLineAligned(text, banner, alignLeft, 0)
}

func renderLineAligned(text string, banner []string, align alignment, width int) (string, error) {
	var result strings.Builder
	rowParts := make([][]string, 8)

	for row := 0; row < 8; row++ { // Each character is represented in 8 rows
		segments := make([]string, 0, len(text))
		for _, char := range text {
			if char < 32 || char > 126 {
				return "", fmt.Errorf("unsupported character: %q", char)
			}

			start := 1 + (int(char)-32)*9 // Each character takes up 9 lines in the banner (8 for the art + 1 empty line)
			index := start + row

			if index >= len(banner) {
				return "", fmt.Errorf("invalid banner file")
			}

			segments = append(segments, banner[index])
		}
		rowParts[row] = segments
	}

	for row, segments := range rowParts {
		result.WriteString(alignLine(segments, text, align, width))
		if row < len(rowParts)-1 {
			result.WriteString("\n")
		}
	}
	result.WriteString("\n")

	return result.String(), nil
}

func alignLine(segments []string, text string, align alignment, width int) string {
	lineWidth := 0
	for _, segment := range segments { // width of ascii art
		lineWidth += len(segment)
	}
	line := strings.Join(segments, "")
	if width <= lineWidth || align == alignLeft { // No alignment needed
		return line
	}

	padding := width - lineWidth // difference of terminal and ascii art
	switch align {
	case alignRight:
		return strings.Repeat(" ", padding) + line
	case alignCenter:
		return strings.Repeat(" ", padding/2) + line
	case alignJustify:
		return justifyLine(segments, text, padding)
	default:
		return line
	}
}

func justifyLine(segments []string, text string, padding int) string {
	spaceCount := strings.Count(text, " ")
	if spaceCount == 0 {
		return strings.Join(segments, "")
	}

	base, remainder := padding/spaceCount, padding%spaceCount
	var result strings.Builder
	spaceIndex := 0
	for index, segment := range segments {
		result.WriteString(segment)
		if index < len(text) && text[index] == ' ' { // Add spaces after each space character in the original text
			spaces := base
			if spaceIndex < remainder {
				spaces++
			}
			result.WriteString(strings.Repeat(" ", spaces))
			spaceIndex++
		}
	}
	return result.String()
}

func terminalWidth() int {
	var windowSize struct {
		rows    uint16
		columns uint16
		xpixel  uint16
		ypixel  uint16
	}
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(os.Stdout.Fd()), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&windowSize)), 0, 0, 0)
	if errno == 0 && windowSize.columns > 0 {
		return int(windowSize.columns)
	}
	return 80
}
