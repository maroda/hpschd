package main

import (
	"encoding/json"
	"io"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/rs/zerolog/log"
)

/*
	This is the refactored v2 API
	All files in the pattern `*sys[_test].go` belong to this version.
*/

// Mesostic contains all data and their transformation for a single poem
type Mesostic struct {
	MU         sync.Mutex
	Date       string      `json:"date"`
	Source     io.Reader   `json:"source"`     // Source JSON
	SourceData interface{} `json:"source_raw"` // Source data, decoded from JSON
	SourceTxt  string      `json:"source_txt"` // Raw text to transform
	Title      string      `json:"title"`      // The poem title
	Spine      []string    `json:"spine"`      // The poem spine
	SpineIdx   int         `json:"spine_idx"`  // Current spine char
	Width      int         `json:"width"`      // Longest line length
	WWidth     int         `json:"west_width"` // Longest westline length
	MLines     []string    `json:"lines"`      // Mesostic lines
	MLinesIdx  int         `json:"lines_idx"`  // Current line
	MLineCt    int         `json:"line_count"` // Total lines
	LineWest   []string    `json:"line_west"`  // Used for alignment
	LineEast   []string    `json:"line_east"`  // Used for alignment
	EmptyLine  []int       `json:"empty_line"` // Line address for empty space
	Poem       string      `json:"poem"`       // Final multi-line poem
}

func NewMesostic(title, source string, data interface{}) *Mesostic {
	m := &Mesostic{
		MU:         sync.Mutex{},
		Date:       "",
		Source:     strings.NewReader(source),
		SourceData: data,
		SourceTxt:  "",
		Title:      title,
		Spine:      make([]string, 0),
		SpineIdx:   0,
		Width:      0,
		MLines:     make([]string, 0),
		MLinesIdx:  0,
		MLineCt:    0,
		LineWest:   make([]string, 0),
		LineEast:   make([]string, 0),
		EmptyLine:  make([]int, 0),
		Poem:       "",
	}
	m.ParseSpine()
	m.ParseSourceJSON(data)
	return m
}

// BuildMeso takes the populated struct and builds the final poem
// replaces mesoMain
func (m *Mesostic) BuildMeso() string {
	m.MU.Lock()
	defer m.MU.Unlock()

	var mesostic string

	// Split the source text up into (hopefully) usable blocks
	re := regexp.MustCompile(`[,.;:]`)
	sourceLines := re.Split(m.SourceTxt, -1)

	for _, sl := range sourceLines {
		if m.FormatLine(sl) {
			m.SpineIdx++
		} else {
			m.EmptyLine = append(m.EmptyLine, m.MLinesIdx)
		}

		// Advance line every time (non-matches get line-returns)
		m.MLinesIdx++
	}

	// Pull all elements together into final mesostic lines
	m.FormatFullLines()

	// Build and return the full test
	for _, ml := range m.MLines {
		mesostic += "\n" + ml
	}
	m.Poem = mesostic
	return mesostic
}

// FormatFullLines builds the final line entries
// Caller holds the lock
func (m *Mesostic) FormatFullLines() bool {
	var line string

	for i, lw := range m.LineWest {
		east := m.LineEast[i]
		west := strings.Repeat(" ", m.WWidth-len(lw)) + lw
		line = west + east

		// Is this line listed as Empty?
		// TODO: This logic doesn't work?
		/*
			for _, el := range m.EmptyLine {
				if i == el {
					line = ""
				}
			}
		*/

		// This should be the only place m.MLines is modified
		m.MLines = append(m.MLines, line)
	}
	return true
}

func isStruct(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Kind() == reflect.Struct
}

// FormatLine creates the mesostic line
// Caller holds the lock
func (m *Mesostic) FormatLine(line string) bool {
	ssChar := m.Spine[m.SpineIdx]
	nxChar := m.Spine[(m.SpineIdx+1)%len(m.Spine)]
	chars := make(map[string][]string)
	lowerline := strings.ToLower(line)

	// Step through each rune in the line
	mode := "west"
	for _, c := range lowerline {
		char := string(c)

		if char != ssChar { // Not the Spine String, either side of it
			// If the next char in the SS is found,
			// drop from here and start the next line
			if mode == "east" && char == nxChar {
				break
			}
			chars[mode] = append(chars[mode], char)
		} else if char == ssChar { // This is the Spine String
			// If this is a repeat of the SS char,
			// drop from here and start the next line
			if mode == "east" {
				break
			}
			char = strings.ToUpper(char)
			chars[mode] = append(chars[mode], char)
			mode = "east"
		}
	}

	// Any line that makes it through with mode=west isn't used.
	if mode == "west" {
		return false
	}

	// Record each side of the Spine String as west|east lines
	westline := strings.TrimSpace(strings.Join(chars["west"], ""))
	eastline := strings.TrimSpace(strings.Join(chars["east"], ""))
	m.LineWest = append(m.LineWest, westline)
	m.LineEast = append(m.LineEast, eastline)

	// Record the widest line
	mline := westline + eastline
	m.Width = wider(len(mline), m.Width)
	m.WWidth = wider(len(westline), m.WWidth)
	return true
}

func wider(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ParseSourceJSON validates and transforms the raw text into usable entry text (???)
// It takes a pointer to the struct for decoding.
func (m *Mesostic) ParseSourceJSON(ps interface{}) bool {
	m.MU.Lock()
	defer m.MU.Unlock()

	if !isStruct(ps) {
		log.Error().Msg("Not a recognized decode target for JSON source")
		return false
	}

	decoder := json.NewDecoder(m.Source)
	if err := decoder.Decode(&ps); err != nil {
		log.Error().Msg("Mesostic.JSON.Decoder Error")
		return false
	}

	return true
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
