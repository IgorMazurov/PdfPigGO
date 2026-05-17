// Package fonts provides types for PDF font handling.
package fonts

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/uglytoad/pdfpig/go/fonts/util"
)

//go:embed resources/glyphlist/*
var glyphListFS embed.FS

// GlyphListFactory loads and parses glyph list files.
type GlyphListFactory struct{}

// Get loads one or more embedded glyph lists by name and returns a combined GlyphList.
// Known names are "glyphlist", "additional", and "zapfdingbats".
func (GlyphListFactory) Get(listNames ...string) (*GlyphList, error) {
	if len(listNames) == 0 {
		return nil, fmt.Errorf("at least one list name is required")
	}

	capacity := 0
	for _, n := range listNames {
		if strings.EqualFold(n, "glyphlist") {
			capacity = 4300
			break
		}
	}

	result := make(map[string]string, capacity)

	for _, listName := range listNames {
		data, err := glyphListFS.ReadFile("resources/glyphlist/" + listName)
		if err != nil {
			return nil, fmt.Errorf("no embedded glyph list resource was found with the name %s: %w", listName, err)
		}

		reader := strings.NewReader(string(data))
		if parseErr := readInternal(reader, result); parseErr != nil {
			return nil, parseErr
		}
	}

	return NewGlyphList(result), nil
}

// Read parses a glyph list from the given reader and returns a GlyphList.
func (GlyphListFactory) Read(stream io.Reader) (*GlyphList, error) {
	if stream == nil {
		return nil, fmt.Errorf("stream must not be nil")
	}

	result := make(map[string]string)
	if err := readInternal(stream, result); err != nil {
		return nil, err
	}

	return NewGlyphList(result), nil
}

func readInternal(reader io.Reader, result map[string]string) error {
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		if line == "" {
			continue
		}

		if line[0] == '#' {
			continue
		}

		parts := strings.Split(line, ";")
		var nonEmptyParts []string
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				nonEmptyParts = append(nonEmptyParts, trimmed)
			}
		}

		if len(nonEmptyParts) != 2 {
			return fmt.Errorf("the line in the glyph list did not match the expected format at line %d: %s", lineNum, line)
		}

		key := nonEmptyParts[0]
		valueStr := nonEmptyParts[1]

		splitter := util.NewStringSplitter(valueStr, ' ')
		var value strings.Builder

		for {
			token, ok := splitter.TryRead()
			if !ok {
				break
			}

			code, parseErr := strconv.ParseUint(token, 16, 32)
			if parseErr != nil {
				return fmt.Errorf("failed to parse hex code %q at line %d: %w", token, lineNum, parseErr)
			}

			value.WriteRune(rune(code))
		}

		result[key] = value.String()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading glyph list: %w", err)
	}

	return nil
}

var glyphListFactoryInstance GlyphListFactory

// GetGlyphListFactory returns the singleton GlyphListFactory instance.
func GetGlyphListFactory() *GlyphListFactory {
	return &glyphListFactoryInstance
}
