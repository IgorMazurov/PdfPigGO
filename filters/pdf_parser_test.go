//go:build integration

package filters

import (
	"bytes"
	"compress/flate"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

const integrationDocRoot = "testdata/integration/Documents"

func TestPdfParser_CanDecompressNormalObjectStream(t *testing.T) {
	fileName := "Single Page Simple - from google drive.pdf"
	filePath := filepath.Join(integrationDocRoot, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", filePath, err)
	}

	obj7Offset := getOffset(data, []byte("7 0 obj"), 0)
	streamPos := getOffset(data, []byte("stream"), obj7Offset.end)
	endStreamPos := getOffset(data, []byte("endstream"), streamPos.end)

	streamBytes := bytesBetween(streamPos.end+1, endStreamPos.start-1, data)

	reader := bytes.NewReader(streamBytes[2:])

	deflateReader := flate.NewReader(reader)
	defer deflateReader.Close()

	_, err = io.ReadAll(deflateReader)
	if err != nil {
		t.Fatalf("flate decompress: %v", err)
	}
}

func TestPdfParser_CanDecompressPngEncodedFlateStream(t *testing.T) {
	fileName := "Single Page Simple - from google drive.pdf"
	filePath := filepath.Join(integrationDocRoot, fileName)

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", filePath, err)
	}

	streamPos := getOffset(data, []byte("stream"), 0)
	endStreamPos := getOffset(data, []byte("endstream"), streamPos.end)

	streamBytes := bytesBetween(streamPos.end+1, endStreamPos.start-1, data)

	paramsDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Predictor: tokens.NewNumericToken(12),
		tokens.Columns:   tokens.NewNumericToken(4),
	})

	dictionary, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Filter:     tokens.FlateDecode,
		tokens.DecodeParms: paramsDict,
	})

	filter := NewFlateFilter()
	filtered, err := filter.Decode(streamBytes, dictionary, Instance, 0)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	expectedStr := "1 0 15 0 1 0 216 0 1 2 160 0 1 2 210 0 1 3 84 0 1 4 46 0 1 7 165 0 1 70 229 0 1 72 84 0 1 96 235 0 1 98 18 0 2 0 12 0 2 0 12 1 2 0 12 2 2 0 12 3 2 0 12 4 2 0 12 5 2 0 12 6 2 0 12 7 2 0 12 8"
	expected := parseExpectedBytes(expectedStr)

	if !bytes.Equal(filtered, expected) {
		t.Errorf("Decoded bytes mismatch.\nexpected: %v\n  actual: %v", expected, filtered)
	}
}

type byteOffset struct {
	start int
	end   int
}

func getOffset(input []byte, term []byte, offsetStart int) byteOffset {
	for i := offsetStart; i < len(input); i++ {
		if input[i] == term[0] {
			found := true
			for j := 1; j < len(term); j++ {
				if i+j >= len(input) || input[i+j] != term[j] {
					found = false
					break
				}
			}
			if found {
				return byteOffset{start: i, end: i + len(term)}
			}
		}
	}
	return byteOffset{}
}

func bytesBetween(start int, endExclusive int, data []byte) []byte {
	result := make([]byte, 0, endExclusive-start)
	for i := start; i < endExclusive; i++ {
		result = append(result, data[i])
	}
	return result
}

func parseExpectedBytes(s string) []byte {
	parts := strings.Split(s, " ")
	result := make([]byte, len(parts))
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		result[i] = byte(v)
	}
	return result
}
