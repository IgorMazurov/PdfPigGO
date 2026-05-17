//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"strings"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/outline"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/writer"
)

// testBookmarkData holds a single test case for bookmark character tests.
type testBookmarkData struct {
	name  string
	words string
}

// TestCanGetBookmarksNonAscii verifies that bookmarks with non-ASCII characters
// are correctly written to and read from a PDF document. This matches C#
// NonAsciiCharactersBookmarksTests.CanGetBookmarks.
func TestCanGetBookmarksNonAscii(t *testing.T) {
	testCases := []testBookmarkData{
		// --- TestData_Pass (originally passing in C#) ---
		{name: "French_Diacritics_Uppercase", words: "É À È Ù Â Ê Î Ô Û Ë Ï Ü Ç Œ Æ"},
		{name: "French_Diacritics_Lowercase", words: "é à è ù â ê î ô û ë ï ü ç œ æ"},

		// Greek Alphabet
		{name: "Greek_Uppercase", words: "Α Β Γ Δ Ε Ζ Η Θ Ι Κ Λ Μ Ν Ξ Ο Π Ρ Σ Τ Υ Φ Χ Ψ Ω"},
		{name: "Greek_Lowercase", words: "α β γ δ ε ζ η θ ι κ λ μ ν ξ ο π ρ σ τ υ φ χ ψ ω"},

		// Cyrillic Small Letter
		{name: "Cyrillic_Small", words: "а б в г д е ж з и к л м н о п р с т у ф х ц ч ш щ ы э ю я"},

		// Hangul Choseong
		{name: "Hangul_Choseong", words: "ㄱ ㄴ ㄷ ㄹ ㅁ ㅂ ㅅ ㅇ ㅈ ㅊ ㅋ ㅌ ㅍ ㅎ"},

		// Fullwidth Latin Small Letter
		{name: "Fullwidth_Latin_Small", words: "ａ ｂ ｃ ｄ ｅ ｆ ｇ ｈ ｉ ｊ ｋ ｌ ｍ ｎ ｏ ｐ ｑ ｒ ｓ ｔ ｕ ｖ ｗ ｘ ｙ ｚ"},

		// Halfwidth Katakana
		{name: "Halfwidth_Katakana", words: "ｱ ｲ ｳ ｴ ｵ ｶ ｷ ｸ ｹ ｺ ｻ ｼ ｽ ｾ ｿ ﾀ ﾁ ﾂ ﾃ ﾄ ﾅ ﾆ ﾇ ﾈ ﾉ ﾊ ﾋ ﾌ ﾍ ﾎ ﾏ ﾐ ﾑ ﾒ ﾓ ﾔ ﾕ ﾖ ﾗ ﾘ ﾙ ﾚ ﾛ ﾜ ｦ ﾝ"},

		// Fullwidth Katakana
		{name: "Fullwidth_Katakana", words: "ア イ ウ エ オ カ キ ク ケ コ サ シ ス セ ソ タ チ ツ テ ト ナ ニ ヌ ネ ノ ハ ヒ フ ヘ ホ マ ミ ム メ モ ヤ ユ ヨ ラ リ ル レ ロ ワ ヲ ン"},

		// Fullwidth Hiragana
		{name: "Fullwidth_Hiragana", words: "あ い う え お か き く け こ さ し す せ そ た ち つ て と な に ぬ ね の は ひ ふ へ ほ ま み む め も や ゆ よ ら り る れ ろ わ を ん"},

		// Kanji (Surrogate Pair)
		{name: "Kanji_SurrogatePair", words: "𩸽 𩹉 𡵅"},

		// Emoji
		{name: "Emoji", words: "🏠 🚗 📝"},

		// Emoji (with ZWJ Sequences)
		{name: "Emoji_ZWJ", words: "👨‍💻 👁‍🗨 😶‍🌫️"},

		// --- TestData_Failed (originally failing in C#) ---
		{name: "Mixed_Cyrillic_Korean_Japanese", words: "ШЩＨＩ차岸岩還館小少尚"},
		{name: "Cyrillic_Space_Separated_1", words: "A Ш Z"},
		{name: "Cyrillic_NoSpace_1", words: "AШZ"},
		{name: "Cyrillic_Space_Separated_2", words: "A Щ Z"},
		{name: "Cyrillic_NoSpace_2", words: "AЩZ"},
		{name: "Fullwidth_Latin_Space", words: "Ｈ Ｉ"},
		{name: "Fullwidth_Latin_NoSpace", words: "ＨＩ"},
		{name: "Hangul_Single", words: "차"},
		{name: "Kanji_Space_1", words: "岸 岩"},
		{name: "Kanji_NoSpace_1", words: "岸岩"},
		{name: "Kanji_Space_2", words: "還 館"},
		{name: "Kanji_NoSpace_2", words: "還館"},
		{name: "Kanji_Space_3", words: "小 少 尚"},
		{name: "Kanji_NoSpace_3", words: "小少尚"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Split input words by space and filter empty strings.
			inputs := splitAndFilter(tc.words)

			// Build a PDF with bookmarks for each word.
			bytes := buildPDFWithBookmarks(t, inputs)

			// Open the built PDF and read bookmark data.
			doc, err := pdfpig.Open(bytes, &content.ParsingOptions{})
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			defer doc.Close()

			resultAny, isSuccess, err := doc.TryGetBookmarks(false)
			if err != nil {
				t.Fatalf("TryGetBookmarks: %v", err)
			}

			if !isSuccess {
				t.Fatal("expected bookmarks to be present")
			}

			bookmarks, ok := resultAny.(*outline.Bookmarks)
			if !ok {
				t.Fatalf("expected *outline.Bookmarks, got %T", resultAny)
			}

			nodes := bookmarks.GetNodes()
			results := make([]string, len(nodes))
			for i, node := range nodes {
				results[i] = node.Title
			}

			if !stringSlicesEqual(inputs, results) {
				t.Errorf("bookmark titles mismatch:\n  expected: %q\n  got:      %q", inputs, results)
			}
		})
	}
}

// splitAndFilter splits a string by spaces and removes empty entries.
func splitAndFilter(s string) []string {
	parts := strings.Split(s, " ")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if len(p) > 0 {
			result = append(result, p)
		}
	}
	return result
}

// buildPDFWithBookmarks creates a minimal PDF containing one A4 page and
// bookmarks for each title in the given slice.
func buildPDFWithBookmarks(t *testing.T, titles []string) []byte {
	t.Helper()

	builder := writer.NewPdfDocumentBuilder()

	_, err := builder.AddPageWithSize(595, 842)
	if err != nil {
		t.Fatalf("AddPageWithSize: %v", err)
	}

	// Create bookmark nodes for each title.
	nodes := make([]outline.BookmarkNodeTyped, len(titles))
	for i, title := range titles {
		dest := destinations.NewExplicitDestination(
			1,
			destinations.XyzCoordinates,
			destinations.Empty,
		)
		children := make([]*outline.BookmarkNode, 0)
		docNode, err := outline.NewDocumentBookmarkNode(title, 0, dest, children)
		if err != nil {
			t.Fatalf("NewDocumentBookmarkNode(%q): %v", title, err)
		}
		nodes[i] = docNode
	}

	builder.SetBookmarks(outline.NewBookmarksTyped(nodes))

	bytes, err := builder.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	return bytes
}

// stringSlicesEqual reports whether two string slices are equal element-wise.
func stringSlicesEqual(a, b []string) bool {
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
