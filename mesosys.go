package main

import (
	"strings"
	"sync"
	"unicode"
)

// Mesostic is here to begin refactoring the poem itself to struct-based data
// The purpose is to get rid of the need for global variables.
type Mesostic struct {
	MU        sync.Mutex
	Date      string   `json:"date"`
	Title     string   `json:"title"`
	Spine     []string `json:"spine"`
	SpineIdx  int      `json:"spine_idx"`
	TextEntry string   `json:"text_entry"`
	MLines    []string `json:"lines"`
	MLinesIdx int      `json:"lines_idx"`
	MLineCt   int      `json:"line_count"`
}

// ParseSpine changes the Title into a lowercase slice without whitespace
func (m *Mesostic) ParseSpine() {
	m.MU.Lock()
	defer m.MU.Unlock()
	for _, c := range m.Title {
		if !unicode.IsSpace(c) {
			m.Spine = append(m.Spine, strings.ToLower(string(c)))
		}
	}
}

// FormatLine would be called by a new mesoMain
func (m *Mesostic) FormatLine(line string) bool {
	ssChar := m.Spine[m.SpineIdx]
	chars := make(map[string][]string)

	// Step through each rune in the line
	mode := "west"
	for _, c := range line {
		char := string(c)
		if char != ssChar {
			chars[mode] = append(chars[mode], char)
		} else if char == ssChar { // This is the Spine String
			char = strings.ToUpper(char)
			chars[mode] = append(chars[mode], char)
			mode = "east"
		}
	}

	mline := strings.Join(chars["west"], "") + strings.Join(chars["east"], "")
	m.MLines = append(m.MLines, mline)
	m.SpineIdx++
	return true
}
