package truetypeparser_test

import (
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/fonts/truetype/tables"
)

var testdataDir = filepath.Join("..", "testdata")

func readFileBytes(t *testing.T, name string) []byte {
	t.Helper()

	if !strings.HasSuffix(name, ".ttf") && !strings.HasSuffix(name, ".txt") {
		name += ".ttf"
	}

	path := filepath.Join(testdataDir, name)

	buf, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read test file %q: %v", path, err)
	}

	return buf
}

func macTimestampToTime(ts int64) time.Time {
	base := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	return base.Add(time.Duration(ts) * time.Second)
}

func assertFloat64Approx(t *testing.T, got, want float64, epsilon float64) {
	t.Helper()
	if math.Abs(got-want) > epsilon {
		t.Errorf("float mismatch: got %g, want %g (epsilon %g)", got, want, epsilon)
	}
}

// TestParseRegularRoboto verifies parsing of the Roboto-Regular font and checks
// all header table fields against known expected values.
func TestParseRegularRoboto(t *testing.T) {
	bytes := readFileBytes(t, "Roboto-Regular")

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	assertFloat64Approx(t, float64(font.Version()), 1.0, 1e-6)

	h := font.TableRegister().HeaderTable

	assertFloat64Approx(t, float64(h.Version()), 1.0, 1e-6)
	assertFloat64Approx(t, float64(h.FontRevision()), 1.0, 1e-6)

	if h.CheckSumAdjustment() != uint32(1142661421) {
		t.Errorf("CheckSumAdjustment = %d, want %d", h.CheckSumAdjustment(), 1142661421)
	}
	if h.MagicNumber() != uint32(1594834165) {
		t.Errorf("MagicNumber = %d, want %d", h.MagicNumber(), 1594834165)
	}
	if h.Flags() != uint16(9) {
		t.Errorf("Flags = %d, want %d", h.Flags(), 9)
	}
	if h.UnitPerEm() != uint16(2048) {
		t.Errorf("UnitPerEm = %d, want %d", h.UnitPerEm(), 2048)
	}

	created := macTimestampToTime(h.Created())
	if created.Year() != 2008 {
		t.Errorf("Created.Year = %d, want %d", created.Year(), 2008)
	}
	if int(created.Month()) != 9 {
		t.Errorf("Created.Month = %d, want %d", created.Month(), 9)
	}
	if created.Day() != 12 {
		t.Errorf("Created.Day = %d, want %d", created.Day(), 12)
	}
	if created.Hour() != 12 {
		t.Errorf("Created.Hour = %d, want %d", created.Hour(), 12)
	}
	if created.Minute() != 29 {
		t.Errorf("Created.Minute = %d, want %d", created.Minute(), 29)
	}
	if created.Second() != 34 {
		t.Errorf("Created.Second = %d, want %d", created.Second(), 34)
	}

	modified := macTimestampToTime(h.Modified())
	if modified.Year() != 2011 {
		t.Errorf("Modified.Year = %d, want %d", modified.Year(), 2011)
	}
	if int(modified.Month()) != 11 {
		t.Errorf("Modified.Month = %d, want %d", modified.Month(), 11)
	}
	if modified.Day() != 30 {
		t.Errorf("Modified.Day = %d, want %d", modified.Day(), 30)
	}
	if modified.Hour() != 5 {
		t.Errorf("Modified.Hour = %d, want %d", modified.Hour(), 5)
	}
	if modified.Minute() != 13 {
		t.Errorf("Modified.Minute = %d, want %d", modified.Minute(), 13)
	}
	if modified.Second() != 10 {
		t.Errorf("Modified.Second = %d, want %d", modified.Second(), 10)
	}

	bounds := h.Bounds()
	if int(bounds.Left()) != -980 {
		t.Errorf("Bounds.Left = %g, want %g", bounds.Left(), -980.0)
	}
	if int(bounds.Bottom()) != -555 {
		t.Errorf("Bounds.Bottom = %g, want %g", bounds.Bottom(), -555.0)
	}
	if int(bounds.Right()) != 2396 {
		t.Errorf("Bounds.Right = %g, want %g", bounds.Right(), 2396.0)
	}
	if int(bounds.Top()) != 2163 {
		t.Errorf("Bounds.Top = %g, want %g", bounds.Top(), 2163.0)
	}

	if h.MacStyle() != tables.HeaderMacStyleNone {
		t.Errorf("MacStyle = %d, want %d", h.MacStyle(), tables.HeaderMacStyleNone)
	}
	if h.LowestRecommendedPpem() != uint16(9) {
		t.Errorf("LowestRecommendedPpem = %d, want %d", h.LowestRecommendedPpem(), 9)
	}

	if h.FontDirectionHint() != tables.StronglyLeftToRightWithNeutrals {
		t.Errorf("FontDirectionHint = %d, want %d", h.FontDirectionHint(), tables.StronglyLeftToRightWithNeutrals)
	}

	if h.IndexToLocFormat() != tables.IndexToLocationTableShort {
		t.Errorf("IndexToLocFormat = %d, want %d", h.IndexToLocFormat(), tables.IndexToLocationTableShort)
	}
	if h.GlyphDataFormat() != int16(0) {
		t.Errorf("GlyphDataFormat = %d, want %d", h.GlyphDataFormat(), 0)
	}
}

// TestRobotoHeaderReadCorrectly verifies that every table directory entry in the
// Roboto-Regular font has the expected offset, length, and checksum.
func TestRobotoHeaderReadCorrectly(t *testing.T) {
	type tableEntry struct {
		tag      string
		offset   uint32
		length   uint32
		checksum uint32
	}

	data := []tableEntry{
		{"DSIG", 158596, 8, 1},
		{"GDEF", 316, 72, 408950881},
		{"GPOS", 388, 35744, 355098641},
		{"GSUB", 36132, 662, 3357985284},
		{"OS/2", 36796, 96, 3097700805},
		{"cmap", 36892, 1750, 298470964},
		{"cvt ", 156132, 38, 119085513},
		{"fpgm", 156172, 2341, 2494100564},
		{"gasp", 156124, 8, 16},
		{"glyf", 38644, 88820, 3302131736},
		{"head", 127464, 54, 346075833},
		{"hhea", 127520, 36, 217516755},
		{"hmtx", 127556, 4148, 1859679943},
		{"kern", 131704, 12306, 2002873469},
		{"loca", 144012, 2076, 77421448},
		{"maxp", 146088, 32, 89459325},
		{"name", 146120, 830, 44343214},
		{"post", 146952, 9171, 3638780613},
		{"prep", 158516, 77, 251381919},
	}

	bytes := readFileBytes(t, "Roboto-Regular")

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	tableHeaders := font.TableHeaders()

	for _, entry := range data {
		match, ok := tableHeaders[entry.tag]
		if !ok {
			t.Errorf("table header %q not found", entry.tag)
			continue
		}

		if match.Offset != entry.offset {
			t.Errorf("%s Offset = %d, want %d", entry.tag, match.Offset, entry.offset)
		}
		if match.Length != entry.length {
			t.Errorf("%s Length = %d, want %d", entry.tag, match.Length, entry.length)
		}
		if match.CheckSum != entry.checksum {
			t.Errorf("%s CheckSum = %d, want %d", entry.tag, match.CheckSum, entry.checksum)
		}
	}
}

// TestParseSimpleGoogleDocssGautmi verifies that the google-simple-doc font
// parses successfully and has a valid header table.
func TestParseSimpleGoogleDocssGautmi(t *testing.T) {
	bytes := readFileBytes(t, "google-simple-doc")

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	if font.TableRegister().HeaderTable.Tag() == "" {
		t.Error("HeaderTable tag is empty")
	}
}

// TestParseAndadaRegular verifies parsing of the Andada-Regular font and checks
// name, revision, flags, unitsPerEm, and date fields.
func TestParseAndadaRegular(t *testing.T) {
	bytes := readFileBytes(t, "Andada-Regular")

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	name := font.Name()
	if name != "Andada Regular" {
		t.Errorf("Name = %q, want %q", name, "Andada Regular")
	}

	h := font.TableRegister().HeaderTable

	assertFloat64Approx(t, float64(h.FontRevision()), 1.001999, 0.00001)

	if h.Flags() != uint16(11) {
		t.Errorf("Flags = %d, want %d", h.Flags(), 11)
	}

	if h.UnitPerEm() != uint16(1000) {
		t.Errorf("UnitPerEm = %d, want %d", h.UnitPerEm(), 1000)
	}

	created := macTimestampToTime(h.Created())
	if created.Year() != 2011 {
		t.Errorf("Created.Year = %d, want %d", created.Year(), 2011)
	}
	if int(created.Month()) != 9 {
		t.Errorf("Created.Month = %d, want %d", created.Month(), 9)
	}
	if created.Day() != 30 {
		t.Errorf("Created.Day = %d, want %d", created.Day(), 30)
	}

	modified := macTimestampToTime(h.Modified())
	if modified.Year() != 2017 {
		t.Errorf("Modified.Year = %d, want %d", modified.Year(), 2017)
	}
	if int(modified.Month()) != 5 {
		t.Errorf("Modified.Month = %d, want %d", modified.Month(), 5)
	}
	if modified.Day() != 4 {
		t.Errorf("Modified.Day = %d, want %d", modified.Day(), 4)
	}
}

// TestParsePMingLiU verifies that the PMingLiU font parses successfully.
func TestParsePMingLiU(t *testing.T) {
	bytes := readFileBytes(t, "PMingLiU")

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}
}

// TestReadsRobotoGlyphSizesCorrectly reads glyph data from a reference text file
// and compares each glyph's width, height, and point count against the parsed font.
func TestReadsRobotoGlyphSizesCorrectly(t *testing.T) {
	re := regexp.MustCompile(`\?: Width (?P<width>\d+), Height: (?P<height>\d+), Points: (?P<points>\d+)`)

	bytes := readFileBytes(t, "Roboto-Regular")

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	glyphDataPath := filepath.Join(testdataDir, "Roboto-Regular.GlyphData.txt")
	raw, err := os.ReadFile(glyphDataPath)
	if err != nil {
		t.Skipf("skipping test: glyph data file not found: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	lines = filterEmpty(lines)

	glyphs := font.TableRegister().GlyphTable.Glyphs()

	for i, line := range lines {
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			t.Errorf("line %d does not match expected pattern: %q", i, line)
			continue
		}

		widthStr := re.SubexpIndex("width")
		heightStr := re.SubexpIndex("height")
		pointsStr := re.SubexpIndex("points")

		width, err := strconv.ParseFloat(matches[widthStr], 64)
		if err != nil {
			t.Errorf("line %d: failed to parse width: %v", i, err)
			continue
		}
		height, err := strconv.ParseFloat(matches[heightStr], 64)
		if err != nil {
			t.Errorf("line %d: failed to parse height: %v", i, err)
			continue
		}
		points, err := strconv.Atoi(matches[pointsStr])
		if err != nil {
			t.Errorf("line %d: failed to parse points: %v", i, err)
			continue
		}

		if i >= len(glyphs) {
			break
		}

		glyph := glyphs[i]
		bounds := glyph.Bounds()

		if width == 0 && height == 0 {
			continue
		}

		if i != 30 {
			assertFloat64Approx(t, bounds.Width, width, 1e-9)
		}

		assertFloat64Approx(t, bounds.Height, height, 1e-9)

		glyphPoints := glyph.Points()
		if len(glyphPoints) != points {
			t.Errorf("glyph[%d] Points.Length = %d, want %d", i, len(glyphPoints), points)
		}

		if points > 0 {
			lastPoint := glyphPoints[len(glyphPoints)-1]
			if !lastPoint.IsEndOfContour {
				t.Errorf("glyph[%d] last point IsEndOfContour = false, want true", i)
			}

			subpaths, ok := glyph.TryGetGlyphPath()
			if !ok {
				t.Errorf("glyph[%d] TryGetGlyphPath returned false", i)
			} else if subpaths == nil {
				t.Errorf("glyph[%d] TryGetGlyphPath returned nil subpaths", i)
			}
		}
	}
}

func filterEmpty(lines []string) []string {
	result := make([]string, 0, len(lines))
	for _, l := range lines {
		if strings.TrimSpace(l) != "" {
			result = append(result, l)
		}
	}
	return result
}

// TestParseIssue258CorruptNameTable verifies that a font with a corrupt name table
// still parses and has accessible name records. This test is skipped if the
// test data file does not exist yet.
func TestParseIssue258CorruptNameTable(t *testing.T) {
	bytes, err := os.ReadFile(filepath.Join(testdataDir, "issue-258-corrupt-name-table.ttf"))
	if err != nil {
		t.Skipf("skipping test: file not found: %v", err)
	}

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	if font.TableRegister().NameTable == nil {
		t.Fatal("NameTable is nil")
	}

	if len(font.TableRegister().NameTable.NameRecords()) == 0 {
		t.Error("NameRecords is empty, want non-empty")
	}
}

// TestParse12623CorruptFileAndGetGlyphs verifies that a corrupt font file still
// parses and allows glyph path retrieval. This test is skipped if the
// test data file does not exist yet.
func TestParse12623CorruptFileAndGetGlyphs(t *testing.T) {
	bytes, err := os.ReadFile(filepath.Join(testdataDir, "corrupt-12623.ttf"))
	if err != nil {
		t.Skipf("skipping test: file not found: %v", err)
	}

	input := truetypeparser.NewTrueTypeDataBytes(bytes)
	font := truetypeparser.Parse(input)

	if font == nil {
		t.Fatal("font is nil")
	}

	font.TryGetPath(1)
}
