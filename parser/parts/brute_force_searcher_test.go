package parts_test

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/parser/filestructure"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const testPdfData = "%PDF-1.5\n" +
	"%test binary header\n" +
	"2 17 obj\n" +
	"<< /Linearized 1 /L 26082 /H [ 722 130 ] /O 6 /E 25807 /N 1 /T 25806 >>\n" +
	"endobj\n" +
	"\n" +
	"3 0 obj\n" +
	"<< /Type /XRef /Length 58 /Filter /FlateDecode /DecodeParms << /Columns 4 /Predictor 12 >> /W [ 1 2 1 ] /Index [ 2 20 ] /Info 13 0 R /Root 4 0 R /Size 22 /Prev 25807                 /ID [<2ee88f3ee8a59b4041754ecb5960c518><2ee88f3ee8a59b4041754ecb5960c518>] >>\n" +
	"stream\n" +
	"xcompressedbinarydatahere\n" +
	"endstream\n" +
	"endobj\n" +
	"4 0 obj\n" +
	"<< /Pages 14 0 R /Type /Catalog >>\n" +
	"endobj\n" +
	"5 0 obj\n" +
	"<< /Filter /FlateDecode /S 36 /Length 53 >>\n" +
	"stream\n" +
	"xcompressedbinarydatahere2\n" +
	"endstream\n" +
	"endobj\n" +
	"\n" +
	"startxref\n" +
	"216\n" +
	"%%EOF"

var testDataOffsets = func() []int64 {
	dataBytes := core.StringAsLatin1Bytes(testPdfData)
	return []int64{
		int64(bytes.Index(dataBytes, []byte("2 17 obj"))),
		int64(bytes.Index(dataBytes, []byte("3 0 obj"))),
		int64(bytes.Index(dataBytes, []byte("4 0 obj"))),
		int64(bytes.Index(dataBytes, []byte("5 0 obj"))),
	}
}()

func TestGetObjectLocationsNilBytes(t *testing.T) {
	_, err := filestructure.GetObjectLocations(nil)
	if err == nil {
		t.Error("expected error for nil bytes")
	}
}

func TestSearcherFindsCorrectObjects(t *testing.T) {
	input := core.NewMemoryInputBytes(core.StringAsLatin1Bytes(testPdfData))

	locations, err := filestructure.GetObjectLocations(input)
	if err != nil {
		t.Fatalf("GetObjectLocations returned error: %v", err)
	}

	if len(locations) != 4 {
		t.Errorf("expected 4 locations, got %d", len(locations))
	}

	actualOffsets := make([]int64, 0, len(locations))
	for _, loc := range locations {
		actualOffsets = append(actualOffsets, loc.Value1)
	}
	sort.Slice(actualOffsets, func(i, j int) bool { return actualOffsets[i] < actualOffsets[j] })

	expectedSorted := make([]int64, len(testDataOffsets))
	copy(expectedSorted, testDataOffsets)
	sort.Slice(expectedSorted, func(i, j int) bool { return expectedSorted[i] < expectedSorted[j] })

	if !equalInt64Slices(actualOffsets, expectedSorted) {
		t.Errorf("expected offsets %v, got %v", expectedSorted, actualOffsets)
	}
}

func TestReaderOnlyCallsOnce(t *testing.T) {
	reader := testutil.ConvertStringToBytes(testPdfData, false)

	locations, err := filestructure.GetObjectLocations(reader.Bytes)
	if err != nil {
		t.Fatalf("GetObjectLocations returned error: %v", err)
	}

	if len(locations) != 4 {
		t.Errorf("expected 4 locations, got %d", len(locations))
	}

	newLocations, err := filestructure.GetObjectLocations(reader.Bytes)
	if err != nil {
		t.Fatalf("second GetObjectLocations call returned error: %v", err)
	}

	if len(newLocations) != 4 {
		t.Errorf("expected 4 locations on second call, got %d", len(newLocations))
	}

	for ref := range locations {
		if _, ok := newLocations[ref]; !ok {
			t.Errorf("second call missing key %s", ref.String())
		}
	}
}

func TestReaderEscapesUnexpectedObject(t *testing.T) {
	s := `%PDF-1.7
abcd

1 0 obj
<< /Type /Any >>

endobj

%AZ 0 obj
11 0 obj
769
endobj

%%EOF`

	bytes := core.NewMemoryInputBytes(core.StringAsLatin1Bytes(s))

	locations, err := filestructure.GetObjectLocations(bytes)
	if err != nil {
		t.Fatalf("GetObjectLocations returned error: %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(locations))
	}

	expectedLocations := []int64{
		int64(strings.Index(s, "1 0 obj")),
		int64(strings.Index(s, "11 0 obj")),
	}

	actualOffsets := make([]int64, 0, len(locations))
	for _, loc := range locations {
		actualOffsets = append(actualOffsets, loc.Value1)
	}
	sort.Slice(actualOffsets, func(i, j int) bool { return actualOffsets[i] < actualOffsets[j] })

	expectedSorted := make([]int64, len(expectedLocations))
	copy(expectedSorted, expectedLocations)
	sort.Slice(expectedSorted, func(i, j int) bool { return expectedSorted[i] < expectedSorted[j] })

	if !equalInt64Slices(actualOffsets, expectedSorted) {
		t.Errorf("expected offsets %v, got %v", expectedSorted, actualOffsets)
	}
}

func TestReaderEscapesUnexpectedGenerationNumber(t *testing.T) {
	s := `%PDF-2.0
abcdefghijklmnop

1 0 obj
256
endobj

16-0 obj

5 0 obj
<< /IsEmpty false >>
endobj`

	bytes := core.NewMemoryInputBytes(core.StringAsLatin1Bytes(s))

	locations, err := filestructure.GetObjectLocations(bytes)
	if err != nil {
		t.Fatalf("GetObjectLocations returned error: %v", err)
	}

	if len(locations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(locations))
	}

	expectedLocations := []int64{
		int64(strings.Index(s, "1 0 obj")),
		int64(strings.Index(s, "5 0 obj")),
	}

	actualOffsets := make([]int64, 0, len(locations))
	for _, loc := range locations {
		actualOffsets = append(actualOffsets, loc.Value1)
	}
	sort.Slice(actualOffsets, func(i, j int) bool { return actualOffsets[i] < actualOffsets[j] })

	expectedSorted := make([]int64, len(expectedLocations))
	copy(expectedSorted, expectedLocations)
	sort.Slice(expectedSorted, func(i, j int) bool { return expectedSorted[i] < expectedSorted[j] })

	if !equalInt64Slices(actualOffsets, expectedSorted) {
		t.Errorf("expected offsets %v, got %v", expectedSorted, actualOffsets)
	}
}

func TestBruteForceSearcherCorrectlyFindsAllObjectsWhenOffset(t *testing.T) {
	input := core.NewMemoryInputBytes(core.StringAsLatin1Bytes(testPdfData))

	input.Seek(593, 0) // io.SeekStart

	locations, err := filestructure.GetObjectLocations(input)
	if err != nil {
		t.Fatalf("GetObjectLocations returned error: %v", err)
	}

	actualOffsets := make([]int64, 0, len(locations))
	for _, loc := range locations {
		actualOffsets = append(actualOffsets, loc.Value1)
	}
	sort.Slice(actualOffsets, func(i, j int) bool { return actualOffsets[i] < actualOffsets[j] })

	expectedSorted := make([]int64, len(testDataOffsets))
	copy(expectedSorted, testDataOffsets)
	sort.Slice(expectedSorted, func(i, j int) bool { return expectedSorted[i] < expectedSorted[j] })

	if !equalInt64Slices(actualOffsets, expectedSorted) {
		t.Errorf("expected offsets %v, got %v", expectedSorted, actualOffsets)
	}
}

func TestBruteForceSearcherFindsAllObjectsWhenMissingEndObj(t *testing.T) {
	s := `%PDF-1.7
abcd

1 0 obj
<< /Type /Any >>

2 0 obj
<< /Type /Any >>

%AZ 0 obj
11 0 obj
769
endobj

%%EOF`

	bytes := core.NewMemoryInputBytes(core.StringAsLatin1Bytes(s))

	locations, err := filestructure.GetObjectLocations(bytes)
	if err != nil {
		t.Fatalf("GetObjectLocations returned error: %v", err)
	}

	if len(locations) != 3 {
		t.Errorf("expected 3 locations, got %d", len(locations))
	}

	expectedLocations := []int64{
		int64(strings.Index(s, "1 0 obj")),
		int64(strings.Index(s, "2 0 obj")),
		int64(strings.Index(s, "11 0 obj")),
	}

	actualOffsets := make([]int64, 0, len(locations))
	for _, loc := range locations {
		actualOffsets = append(actualOffsets, loc.Value1)
	}
	sort.Slice(actualOffsets, func(i, j int) bool { return actualOffsets[i] < actualOffsets[j] })

	expectedSorted := make([]int64, len(expectedLocations))
	copy(expectedSorted, expectedLocations)
	sort.Slice(expectedSorted, func(i, j int) bool { return expectedSorted[i] < expectedSorted[j] })

	if !equalInt64Slices(actualOffsets, expectedSorted) {
		t.Errorf("expected offsets %v, got %v", expectedSorted, actualOffsets)
	}
}

func getStringAt(bytes core.InputBytes, location int64) string {
	bytes.Seek(location, 0) // io.SeekStart
	txt := make([]byte, 10)
	bytes.Read(txt)

	return core.BytesAsLatin1String(txt)
}

func equalInt64Slices(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
