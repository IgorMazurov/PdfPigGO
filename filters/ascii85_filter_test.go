package filters

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/uglytoad/pdfpig/go/tokens"
)

func TestAscii85Filter_DecodesWikipediaExample(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})

	input := []byte("9jqo^BlbD-BleB1DJ+*+F(f,q/0JhKF<GL>Cj@.4Gp$d7F!,L7@<6@)/0JDEF<G%<+EV:2F!,\nO<DJ+*.@<*K0@<6L(Df-\\0Ec5e;DffZ(EZee.Bl.9pF\"AGXBPCsi + DGm >@3BB / F * &OCAfu2 / AKY\n            i(DIb: @FD, *) + C]U =@3BN#EcYf8ATD3s@q?d$AftVqCh[NqF<G:8+EV:.+Cf>-FD5W8ARlolDIa\n            l(DId<j@<? 3r@:F % a + D58'ATD4$Bl@l3De:,-DJs`8ARoFb/0JMK@qB4^F!,R<AKZ&-DfTqBG%G\n                > uD.RTpAKYo'+CT/5+Cei#DII?(E,9)oF*2M7/c~>")

	result, err := filter.Decode(input, dict, Instance, 0)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	expected := "Man is distinguished, not only by his reason, but by this singular passion from other animals, which is a lust of the mind, " +
		"that by a perseverance of delight in the continued and indefatigable generation of knowledge, " +
		"exceeds the short vehemence of any carnal pleasure."

	got := string(result)
	if got != expected {
		t.Errorf("Decoded text mismatch.\nexpected: %q\n  actual: %q", expected, got)
	}
}

func TestAscii85Filter_DecodesHelloWorld(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	tests := []struct {
		name     string
		encoded  string
		expected string
	}{
		{name: "BE", encoded: "BE", expected: "h"},
		{name: "BOq", encoded: "BOq", expected: "he"},
		{name: "BOtu", encoded: "BOtu", expected: "hel"},
		{name: "BOu!r", encoded: "BOu!r", expected: "hell"},
		{name: "BOu!rDZ", encoded: "BOu!rDZ", expected: "hello"},
		{name: "BOu!rD]f", encoded: "BOu!rD]f", expected: "hello "},
		{name: "BOu!rD]j6", encoded: "BOu!rD]j6", expected: "hello w"},
		{name: "BOu!rD]j7B", encoded: "BOu!rD]j7B", expected: "hello wo"},
		{name: "BOu!rD]j7BEW", encoded: "BOu!rD]j7BEW", expected: "hello wor"},
		{name: "BOu!rD]j7BEbk", encoded: "BOu!rD]j7BEbk", expected: "hello worl"},
		{name: "BOu!rD]j7BEbo7", encoded: "BOu!rD]j7BEbo7", expected: "hello world"},
		{name: "BOu!rD]j7BEbo80", encoded: "BOu!rD]j7BEbo80", expected: "hello world!"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := filter.Decode([]byte(tc.encoded), dict, Instance, 0)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			got := string(result)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestAscii85Filter_ReplacesZWithEmptyBytes(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	tests := []struct {
		name     string
		encoded  string
		expected string
	}{
		{name: "wikipedia_z", encoded: "9jqo^zBlbD-", expected: "Man \x00\x00\x00\x00is d" },
		{name: "empty", encoded: "", expected: "" },
		{name: "single_z", encoded: "z", expected: "\x00\x00\x00\x00" },
		{name: "double_z", encoded: "zz", expected: "\x00\x00\x00\x00\x00\x00\x00\x00" },
		{name: "triple_z", encoded: "zzz", expected: "\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00" },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := filter.Decode([]byte(tc.encoded), dict, Instance, 1)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			got := string(result)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestAscii85Filter_ZInMiddleOf5CharacterSequenceThrows(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	input := []byte("qjzqo^")
	_, err := filter.Decode(input, dict, Instance, 0)
	if err == nil {
		t.Fatal("expected error when z appears in middle of sequence, got nil")
	}
	if !strings.Contains(err.Error(), "unexpected empty block marker") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAscii85Filter_SingleCharacterLastIgnores(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	tests := []struct {
		name     string
		encoded  string
		expected string
	}{
		{name: "cool", encoded: "@rH:%B", expected: "cool"},
		{name: "end_marker_only", encoded: "A~>", expected: ""},
		{name: "cool_with_end", encoded: "@rH:%A~>", expected: "cool"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := filter.Decode([]byte(tc.encoded), dict, Instance, 1)
			if err != nil {
				t.Fatalf("Decode error: %v", err)
			}
			got := string(result)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

const pdfContent = `1 0 obj
<< /Length 568 >>
stream
2 J
BT
/F1 12 Tf
0 Tc
0 Tw
72.5 712 TD
[(Unencoded streams can be read easily) 65 (, )] TJ
0 -14 TD
[(b) 20 (ut generally tak) 10 (e more space than \311)] TJ
T* (encoded streams.) Tj
0 -28 TD
[(Se) 25 (v) 15 (eral encoding methods are a) 20 (v) 25 (ailable in PDF) 80 (.)] TJ
0 -14 TD
(Some are used for compression and others simply) Tj
T* [(to represent binary data in an ) 55 (ASCII format.)] TJ
T* (Some of the compression encoding methods are \
suitable ) Tj
T* (for both data and images, while others are \
suitable only ) Tj
T* (for continuous-tone images.) Tj
ET
endstream
endobj`

// pdfContentEncoded is the ASCII85-encoded version with ~> end marker.
const pdfContentEncoded = "0d&.mDdmGg4?O`>9P&*SFD)dS2E2gC4pl@QEb/Zr$8N_r$:7]!01IZ=0eskNAdU47<+?7h+B3Ol2_m!C+?)#1+B1\n`9>:<KhASu!rA7]9oF*)G6@;U'.@ps6t@V$[&ART*lARTXoCj@HP2DlU*/0HBI+B1r?0H_r%1a#ac$<nof.3LB\"+=MAS+D58'ATD3qCj@.\n           F@;@;70ea^uAKYi.Eb-A7E+*6f+EV:*DBN1?0ek+_+B1r?<%9\"=ASu!rA7]9oF*)G6@;U'<.3MT)$8<SS1,pCU6jd-H;e7\nC#1,U1&Ft\"Og2'=;YEa`c,ASu!rA8,po+Dk\\3BQ%F&+CT;%+CQ]A1,'h!Ft\"Oh2'=;UBl%3eCh4`'DBMbD7O]H>0H_br.:\"&q8d[6p/M\nT()<(%'A;f?Ma+CT;%+E_a:A0>K&EZek1D/aN,F)u&6DBNA*A0>f4BOu4*+EM76E,9eK+B3(_<%9\"p.!0AMEb031ATMF#F<G%,DIIR2+Cno\n&@3B9%+CT.1.3LK*+=KNS6V0ilAoD^,@<=+N>p**=$</Jt-rY&$AKYo'+EV:.+Cf>,E,oN2F(oQ1+D#G#De*R\"B-;&&FD,T'F!+n3AKY4b\nF*22=@:F%a+=SF4C'moi+=Li?EZeh0FD)e-@<>p#@;]TuBl.9kATKCFGA(],AKYo5BOu4*+CT;%+C#7pF_Pr+@VfTuDf0B:+=SF4C'moi+=\nLi?EZek1DKKT1F`2DD/TboKAKY](@:s.m/h%oBC'mC/$>\"*cF*)G6@;Q?_DIdZpC&~>"

// pdfContentEncodedNoEnd is the same without the ~> end marker.
const pdfContentEncodedNoEnd = "0d&.mDdmGg4?O`>9P&*SFD)dS2E2gC4pl@QEb/Zr$8N_r$:7]!01IZ=0eskNAdU47<+?7h+B3Ol2_m!C+?)#1+B1\n`9>:<KhASu!rA7]9oF*)G6@;U'.@ps6t@V$[&ART*lARTXoCj@HP2DlU*/0HBI+B1r?0H_r%1a#ac$<nof.3LB\"+=MAS+D58'ATD3qCj@.\n           F@;@;70ea^uAKYi.Eb-A7E+*6f+EV:*DBN1?0ek+_+B1r?<%9\"=ASu!rA7]9oF*)G6@;U'<.3MT)$8<SS1,pCU6jd-H;e7\nC#1,U1&Ft\"Og2'=;YEa`c,ASu!rA8,po+Dk\\3BQ%F&+CT;%+CQ]A1,'h!Ft\"Oh2'=;UBl%3eCh4`'DBMbD7O]H>0H_br.:\"&q8d[6p/M\nT()<(%'A;f?Ma+CT;%+E_a:A0>K&EZek1D/aN,F)u&6DBNA*A0>f4BOu4*+EM76E,9eK+B3(_<%9\"p.!0AMEb031ATMF#F<G%,DIIR2+Cno\n&@3B9%+CT.1.3LK*+=KNS6V0ilAoD^,@<=+N>p**=$</Jt-rY&$AKYo'+EV:.+Cf>,E,oN2F(oQ1+D#G#De*R\"B-;&&FD,T'F!+n3AKY4b\nF*22=@:F%a+=SF4C'moi+=Li?EZeh0FD)e-@<>p#@;]TuBl.9kATKCFGA(],AKYo5BOu4*+CT;%+C#7pF_Pr+@VfTuDf0B:+=SF4C'moi+=\nLi?EZek1DKKT1F`2DD/TboKAKY](@:s.m/h%oBC'mC/$>\"*cF*)G6@;Q?_DIdZpC&"

func TestAscii85Filter_DecodesEncodedPdfContent(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	result, err := filter.Decode([]byte(pdfContentEncoded), dict, Instance, 0)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	got := string(result)
	expected := strings.ReplaceAll(pdfContent, "\r\n", "\n")
	if got != expected {
		t.Errorf("Decoded PDF content mismatch.\nexpected length: %d\n  actual length: %d", len(expected), len(got))
	}
}

func TestAscii85Filter_DecodesEncodedPdfContentMissingEndOfDataSymbol(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	result, err := filter.Decode([]byte(pdfContentEncodedNoEnd), dict, Instance, 0)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	got := string(result)
	expected := strings.ReplaceAll(pdfContent, "\r\n", "\n")
	if got != expected {
		t.Errorf("Decoded PDF content mismatch.\nexpected length: %d\n  actual length: %d", len(expected), len(got))
	}
}

func TestAscii85Filter_DecodeParallel(t *testing.T) {
	filter := NewAscii85Filter()
	dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{})
	var wg sync.WaitGroup
	errChan := make(chan string, 100_000)
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if idx%2 == 0 {
				result, err := filter.Decode([]byte("9jqo^zBlbD-"), dict, Instance, 1)
				if err != nil {
					errChan <- fmt.Sprintf("iteration %d: Decode error: %v", idx, err)
					return
				}
				got := string(result)
				expected := "Man \x00\x00\x00\x00is d"
				if got != expected {
					errChan <- fmt.Sprintf("iteration %d: expected %q, got %q", idx, expected, got)
					return
				}
			} else {
				result, err := filter.Decode([]byte(pdfContentEncoded), dict, Instance, 0)
				if err != nil {
					errChan <- fmt.Sprintf("iteration %d: Decode error: %v", idx, err)
					return
				}
				got := string(result)
				expected := strings.ReplaceAll(pdfContent, "\r\n", "\n")
				if got != expected {
					errChan <- fmt.Sprintf("iteration %d: PDF content mismatch", idx)
					return
				}
			}
		}(i)
	}
	wg.Wait()
	close(errChan)
	for err := range errChan {
		t.Error(err)
	}
}
