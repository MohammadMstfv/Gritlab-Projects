package main

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// Matches a marker with any internal spacing, e.g. "(cap, 6)" or "( up )".
	markerSpacingRe = regexp.MustCompile(`\(\s*(hex|bin|up|low|cap)\s*(,\s*(\d+)\s*)?\)`)

	// Matches a normalized (whitespace-free) marker token.
	hexBinMarkerRe = regexp.MustCompile(`^\((hex|bin)\)$`)
	caseMarkerRe   = regexp.MustCompile(`^\((up|low|cap)(?:,(\d+))?\)$`)

	spaceBeforePunctRe = regexp.MustCompile(`\s+([.,!?:;]+)`)
	spaceAfterPunctRe  = regexp.MustCompile(`([.,!?:;]+)([^\s.,!?:;])`)

	quoteRe = regexp.MustCompile(`'\s*([^']*?)\s*'`)

	articleRe     = regexp.MustCompile(`\b([Aa])\b(\s+)(\S+)`)
	firstLetterRe = regexp.MustCompile(`[A-Za-z]`)
)

// Transform applies every text modification described in the spec and returns
// the result while preserving the original line structure.
func Transform(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = applyTransformSteps(line)
	}
	return strings.Join(lines, "\n")
}

func applyTransformSteps(line string) string {
	steps := []func(string) string{
		normalizeMarkerSpacing,
		resolveMarkers,
		normalizePunctuation,
		tightenQuotes,
		fixArticles,
	}

	for _, step := range steps {
		line = step(line)
	}
	return line
}

// normalizeMarkerSpacing rewrites a marker like "(cap, 6)" (with arbitrary
// internal spacing) into the whitespace-free token "(cap,6)" so a later
// split on whitespace always keeps a whole marker as a single token.
func normalizeMarkerSpacing(line string) string {
	return markerSpacingRe.ReplaceAllStringFunc(line, func(marker string) string {
		parts := markerSpacingRe.FindStringSubmatch(marker)
		kind, count := parts[1], parts[3]
		if count == "" {
			return "(" + kind + ")"
		}
		return "(" + kind + "," + count + ")"
	})
}

// resolveMarkers walks the whitespace-separated tokens of a line left to right.
// On a (hex)/(bin) marker it converts the previous word from that base to
// decimal. On an (up)/(low)/(cap) marker (optionally with a count) it applies
// the case change to the preceding N words. Markers are dropped from output.
func resolveMarkers(line string) string {
	tokens := strings.Fields(line)
	out := make([]string, 0, len(tokens))

	for _, token := range tokens {
		switch {
		case isHexOrBinMarker(token):
			if len(out) == 0 {
				continue
			}
			if converted, ok := convertBase(out[len(out)-1], hexOrBinKind(token)); ok {
				out[len(out)-1] = converted
			}
		case isCaseMarker(token):
			kind, count := caseMarkerValues(token)
			n := clampMarkerCount(count, len(out))
			for i := len(out) - n; i < len(out); i++ {
				out[i] = applyCase(out[i], kind)
			}
		default:
			out = append(out, token)
		}
	}

	return strings.Join(out, " ")
}

func isHexOrBinMarker(token string) bool {
	return hexBinMarkerRe.MatchString(token)
}

func isCaseMarker(token string) bool {
	return caseMarkerRe.MatchString(token)
}

func hexOrBinKind(token string) string {
	match := hexBinMarkerRe.FindStringSubmatch(token)
	if match == nil {
		return ""
	}
	return match[1]
}

func caseMarkerValues(token string) (string, int) {
	match := caseMarkerRe.FindStringSubmatch(token)
	if match == nil {
		return "", 1
	}

	kind := match[1]
	count := 1
	if match[2] != "" {
		if value, err := strconv.Atoi(match[2]); err == nil {
			count = value
		}
	}
	return kind, count
}

func clampMarkerCount(count, available int) int {
	if count > available {
		return available
	}
	if count < 0 {
		return 0
	}
	return count
}

// convertBase parses a word in the given base ("hex" or "bin") and returns
// its base-10 string form.
func convertBase(word, kind string) (string, bool) {
	base := 16
	if kind == "bin" {
		base = 2
	}
	value, err := strconv.ParseInt(word, base, 64)
	if err != nil {
		return word, false
	}
	return strconv.FormatInt(value, 10), true
}

func applyCase(word, kind string) string {
	switch kind {
	case "up":
		return strings.ToUpper(word)
	case "low":
		return strings.ToLower(word)
	case "cap":
		return capitalize(word)
	default:
		return word
	}
}

// capitalize upper-cases only the first letter of a word, leaving every other
// character (including any punctuation glued to it) untouched.
func capitalize(word string) string {
	runes := []rune(word)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(runes)
}

// normalizePunctuation removes any space before a run of . , ! ? : ; and
// ensures exactly one space after it, so groups like "..." or "!?" stay
// glued together while still being separated from the next word.
func normalizePunctuation(line string) string {
	line = spaceBeforePunctRe.ReplaceAllString(line, "$1")
	line = spaceAfterPunctRe.ReplaceAllString(line, "$1 $2")
	return line
}

// tightenQuotes removes the padding space(s) just inside a pair of single
// quotes, e.g. "' awesome '" -> "'awesome'", while preserving the spacing
// between words enclosed by the quotes.
func tightenQuotes(line string) string {
	return quoteRe.ReplaceAllString(line, "'$1'")
}

// fixArticles turns a standalone "a"/"A" into "an"/"An" when the next word
// begins with a vowel or "h", preserving the article's original case.
func fixArticles(line string) string {
	return articleRe.ReplaceAllStringFunc(line, func(match string) string {
		parts := articleRe.FindStringSubmatch(match)
		article, spacing, nextWord := parts[1], parts[2], parts[3]

		letter := firstLetterRe.FindString(nextWord)
		if letter == "" {
			return match
		}
		if shouldUseAn(letter) {
			if article == "A" {
				article = "An"
			} else {
				article = "an"
			}
		}
		return article + spacing + nextWord
	})
}

func shouldUseAn(letter string) bool {
	switch strings.ToLower(letter) {
	case "a", "e", "i", "o", "u", "h":
		return true
	default:
		return false
	}
}
