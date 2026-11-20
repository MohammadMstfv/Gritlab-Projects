package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type asciiFont struct {
	height int
	glyphs map[rune][]string
}

const binaryBannerName = "Binary"

var bannerFiles = map[string]string{
	"standard":   "standard.txt",
	"shadow":     "shadow.txt",
	"thinkertoy": "thinkertoy.txt",
}

func main() {
	if len(os.Args) == 1 {
		http.HandleFunc("/", homeHandler)
		http.HandleFunc("/ascii-art", asciiArtHandler)
		http.HandleFunc("/export", exportHandler)
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

		fmt.Println("Server running at http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Server error:", err)
		}
		return
	}

	if len(os.Args) > 3 {
		fmt.Println("Usage: go run . <text> [banner]")
		return
	}

	bannerName, err := parseBannerName(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	output, err := buildBannerASCII(os.Args[1], bannerName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Print(output)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := os.ReadFile("templates/index.html")
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "unable to read template", http.StatusInternalServerError)
		return
	}

	page, err := template.New("index.html").Parse(string(data))
	if err != nil {
		http.Error(w, "unable to parse template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.Execute(w, nil); err != nil {
		return
	}
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	input := r.FormValue("text")
	bannerName := r.FormValue("banner")
	if input == "" || bannerName == "" {
		http.Error(w, "text and banner are required", http.StatusBadRequest)
		return
	}

	output, err := buildBannerASCII(input, bannerName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "unable to read banner: ") {
			if errors.Is(err, os.ErrNotExist) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, "unable to read banner", http.StatusInternalServerError)
			return
		}
		if !isSupportedBanner(bannerName) {
			http.Error(w, "unknown banner", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(output))
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Query().Get("format") != "txt" {
		http.Error(w, "unsupported export format", http.StatusBadRequest)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	output := r.FormValue("output")
	if output == "" {
		http.Error(w, "output is required", http.StatusBadRequest)
		return
	}

	data := []byte(output)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Content-Disposition", `attachment; filename="ascii-art.txt"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func parseBannerName(args []string) (string, error) {
	bannerName := "standard"
	if len(args) == 3 {
		bannerName = args[2]
	}

	if !isSupportedBanner(bannerName) {
		return "", fmt.Errorf("unknown banner %q", bannerName)
	}

	return bannerName, nil
}

func isSupportedBanner(bannerName string) bool {
	if bannerName == binaryBannerName {
		return true
	}
	_, ok := bannerFiles[bannerName]
	return ok
}

func buildBannerASCII(input, bannerName string) (string, error) {
	if bannerName == binaryBannerName {
		return buildBinaryASCII(input)
	}

	filename, ok := bannerFiles[bannerName]
	if !ok {
		return "", fmt.Errorf("unknown banner %q", bannerName)
	}
	font, err := loadFont(filename)
	if err != nil {
		return "", fmt.Errorf("unable to read banner: %w", err)
	}
	return buildASCII(input, font)
}

func loadFont(filename string) (asciiFont, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return asciiFont{}, err
	}

	lines := strings.Split(normalizeLineEndings(string(data)), "\n")
	if !strings.HasSuffix(strings.ToLower(filename), ".txt") {
		return asciiFont{}, fmt.Errorf("unsupported font format: %s", filename)
	}
	return parseTextFont(lines)
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func parseTextFont(lines []string) (asciiFont, error) {
	font := asciiFont{height: 8, glyphs: make(map[rune][]string, 95)}
	for char := 32; char <= 126; char++ {
		start := 1 + (char-32)*9
		if start+font.height > len(lines) {
			return asciiFont{}, fmt.Errorf("invalid text banner file")
		}
		font.glyphs[rune(char)] = append([]string(nil), lines[start:start+font.height]...)
	}
	return font, nil
}

func buildASCII(input string, font asciiFont) (string, error) {
	if input == "" {
		return "", nil
	}

	input = normalizeInputLineEndings(input)
	lines := strings.Split(input, "\n")

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

		art, err := renderLine(line, font)
		if err != nil {
			return "", err
		}

		result.WriteString(art)
	}

	return result.String(), nil
}

func buildBinaryASCII(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	lines := strings.Split(normalizeInputLineEndings(input), "\n")
	var result strings.Builder
	if allLinesEmpty(lines) {
		for range len(lines) - 1 {
			result.WriteByte('\n')
		}
		return result.String(), nil
	}

	for _, line := range lines {
		for _, char := range line {
			if char < 32 || char > 126 {
				return "", fmt.Errorf("unsupported character: %q", char)
			}
			if char == ' ' {
				result.WriteByte(' ')
			} else {
				result.WriteString(fmt.Sprintf("%08b", char))
			}
		}
		result.WriteByte('\n')
	}
	return result.String(), nil
}

func normalizeInputLineEndings(input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	input = strings.ReplaceAll(input, `\r\n`, "\n")
	input = strings.ReplaceAll(input, `\n`, "\n")
	input = strings.ReplaceAll(input, `\r`, "\n")
	return input
}

func allLinesEmpty(lines []string) bool {
	for _, line := range lines {
		if line != "" {
			return false
		}
	}

	return true
}

func renderLine(text string, font asciiFont) (string, error) {
	rows := make([]string, font.height)
	for _, char := range text {
		if char < 32 || char > 126 {
			return "", fmt.Errorf("unsupported character: %q", char)
		}
		glyph, ok := font.glyphs[char]
		if !ok || len(glyph) != font.height {
			return "", fmt.Errorf("invalid font character %q", char)
		}
		for row := range rows {
			rows[row] += glyph[row]
		}
	}

	var result strings.Builder
	for _, row := range rows {
		result.WriteString(row)
		result.WriteByte('\n')
	}

	return result.String(), nil
}
