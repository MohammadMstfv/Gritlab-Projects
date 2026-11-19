package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type asciiFont struct {
	height    int
	layout    int
	hardblank rune
	glyphs    map[rune][]string
}

var bannerFiles = map[string]string{
	"standard":   "standard.txt",
	"shadow":     "shadow.txt",
	"thinkertoy": "thinkertoy.txt",
	"Binary":     "Binary.flf",
}

func main() {
	if len(os.Args) == 1 {
		http.HandleFunc("/", homeHandler)
		http.HandleFunc("/ascii-art", asciiArtHandler)
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

	font, err := loadFont(bannerFiles[bannerName])
	if err != nil {
		fmt.Println("Error reading font:", err)
		return
	}

	output, err := buildASCII(os.Args[1], font)
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

	filename, ok := bannerFiles[bannerName]
	if !ok {
		http.Error(w, "unknown banner", http.StatusBadRequest)
		return
	}

	font, err := loadFont(filename)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "unable to read banner", http.StatusInternalServerError)
		return
	}

	output, err := buildASCII(input, font)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(output))
}

func parseBannerName(args []string) (string, error) {
	bannerName := "standard"
	if len(args) == 3 {
		bannerName = args[2]
	}

	if _, ok := bannerFiles[bannerName]; !ok {
		return "", fmt.Errorf("unknown banner %q", bannerName)
	}

	return bannerName, nil
}

func loadFont(filename string) (asciiFont, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return asciiFont{}, err
	}

	lines := strings.Split(normalizeLineEndings(string(data)), "\n")
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".txt":
		return parseTextFont(lines)
	case ".flf":
		return parseFLFFont(lines)
	default:
		return asciiFont{}, fmt.Errorf("unsupported font format: %s", filepath.Ext(filename))
	}
}

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func parseTextFont(lines []string) (asciiFont, error) {
	font := asciiFont{height: 8, layout: -1, glyphs: make(map[rune][]string, 95)}
	for char := 32; char <= 126; char++ {
		start := 1 + (char-32)*9
		if start+font.height > len(lines) {
			return asciiFont{}, fmt.Errorf("invalid text banner file")
		}
		font.glyphs[rune(char)] = append([]string(nil), lines[start:start+font.height]...)
	}
	return font, nil
}

func parseFLFFont(lines []string) (asciiFont, error) {
	if len(lines) == 0 {
		return asciiFont{}, fmt.Errorf("invalid FIGlet font header")
	}
	header := strings.Fields(lines[0])
	if len(header) < 6 || !strings.HasPrefix(header[0], "flf2a") || len([]rune(header[0])) < 6 {
		return asciiFont{}, fmt.Errorf("invalid FIGlet font header")
	}

	height, err := strconv.Atoi(header[1])
	if err != nil || height < 1 {
		return asciiFont{}, fmt.Errorf("invalid FIGlet font height")
	}
	commentLines, err := strconv.Atoi(header[5])
	if err != nil || commentLines < 0 {
		return asciiFont{}, fmt.Errorf("invalid FIGlet comment count")
	}
	layout, err := strconv.Atoi(header[4])
	if err != nil {
		return asciiFont{}, fmt.Errorf("invalid FIGlet layout")
	}

	lineIndex := 1 + commentLines
	if lineIndex >= len(lines) || lines[lineIndex] == "" {
		return asciiFont{}, fmt.Errorf("FIGlet font has no character data")
	}
	endmarkRunes := []rune(strings.TrimRight(lines[lineIndex], " \t"))
	if len(endmarkRunes) == 0 {
		return asciiFont{}, fmt.Errorf("invalid FIGlet endmark")
	}
	endmark := endmarkRunes[len(endmarkRunes)-1]
	hardblank := []rune(header[0])[5]
	font := asciiFont{height: height, layout: layout, hardblank: hardblank, glyphs: make(map[rune][]string, 95)}

	for char := 32; char <= 126; char++ {
		glyph := make([]string, height)
		for row := 0; row < height; row++ {
			if lineIndex >= len(lines) {
				return asciiFont{}, fmt.Errorf("FIGlet font is missing character %d", char)
			}
			glyphLine := []rune(lines[lineIndex])
			lineIndex++
			if len(glyphLine) == 0 || glyphLine[len(glyphLine)-1] != endmark {
				return asciiFont{}, fmt.Errorf("invalid FIGlet character %d", char)
			}
			glyphLine = glyphLine[:len(glyphLine)-1]
			if row == height-1 && len(glyphLine) > 0 && glyphLine[len(glyphLine)-1] == endmark {
				glyphLine = glyphLine[:len(glyphLine)-1]
			}
			glyph[row] = strings.TrimRight(string(glyphLine), " \t")
		}
		font.glyphs[rune(char)] = glyph
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
			rows[row] = mergeGlyphRow(rows[row], glyph[row], font.layout, font.hardblank)
		}
	}

	var result strings.Builder
	for _, row := range rows {
		for _, char := range row {
			if char == font.hardblank {
				result.WriteRune(' ')
			} else {
				result.WriteRune(char)
			}
		}
		result.WriteByte('\n')
	}

	return result.String(), nil
}

func mergeGlyphRow(left, right string, layout int, hardblank rune) string {
	if left == "" {
		return right
	}
	if layout < 0 {
		return left + right
	}

	rules := layout & 63
	if layout == 64 {
		rules = 63
	}
	leftRunes := []rune(left)
	rightRunes := []rune(right)
	maxOverlap := len(leftRunes)
	if len(rightRunes) < maxOverlap {
		maxOverlap = len(rightRunes)
	}
	for overlap := maxOverlap; overlap >= 0; overlap-- {
		start := len(leftRunes) - overlap
		valid := true
		for i := 0; i < overlap; i++ {
			leftChar, rightChar := leftRunes[start+i], rightRunes[i]
			if leftChar != ' ' && rightChar != ' ' && !canSmush(leftChar, rightChar, rules, hardblank) {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		merged := append([]rune(nil), leftRunes[:start]...)
		for i := 0; i < overlap; i++ {
			merged = append(merged, smush(leftRunes[start+i], rightRunes[i], rules, hardblank))
		}
		merged = append(merged, rightRunes[overlap:]...)
		return string(merged)
	}
	return left + right
}

func canSmush(left, right rune, rules int, hardblank rune) bool {
	if rules&1 != 0 && left == right {
		return true
	}
	if rules&2 != 0 && (left == '_' || right == '_') {
		other := left
		if other == '_' {
			other = right
		}
		return strings.ContainsRune("|/\\[]{}()<>", other)
	}
	if rules&4 != 0 && hierarchyRank(left) >= 0 && hierarchyRank(right) >= 0 && hierarchyRank(left) != hierarchyRank(right) {
		return true
	}
	if rules&8 != 0 && isOppositePair(left, right) {
		return true
	}
	if rules&16 != 0 && isBigXPair(left, right) {
		return true
	}
	return rules&32 != 0 && left == hardblank && right == hardblank
}

func smush(left, right rune, rules int, hardblank rune) rune {
	if left == ' ' {
		return right
	}
	if right == ' ' {
		return left
	}
	if left == right && rules&1 != 0 {
		return left
	}
	if rules&2 != 0 && (left == '_' || right == '_') {
		if left == '_' {
			return right
		}
		return left
	}
	if rules&4 != 0 {
		leftRank, rightRank := hierarchyRank(left), hierarchyRank(right)
		if leftRank >= 0 && rightRank >= 0 && leftRank != rightRank {
			if leftRank < rightRank {
				return left
			}
			return right
		}
	}
	if rules&8 != 0 && isOppositePair(left, right) {
		return '|'
	}
	if rules&16 != 0 {
		switch string([]rune{left, right}) {
		case "/\\":
			return '|'
		case "\\/":
			return 'Y'
		case "><":
			return 'X'
		}
	}
	if rules&32 != 0 && left == hardblank && right == hardblank {
		return hardblank
	}
	return right
}

func hierarchyRank(char rune) int {
	switch char {
	case '|':
		return 0
	case '/', '\\':
		return 1
	case '[', ']':
		return 2
	case '{', '}':
		return 3
	case '(', ')':
		return 4
	case '<', '>':
		return 5
	default:
		return -1
	}
}

func isOppositePair(left, right rune) bool {
	return left == '[' && right == ']' || left == ']' && right == '[' ||
		left == '{' && right == '}' || left == '}' && right == '{' ||
		left == '(' && right == ')' || left == ')' && right == '('
}

func isBigXPair(left, right rune) bool {
	return left == '/' && right == '\\' || left == '\\' && right == '/' || left == '>' && right == '<'
}
