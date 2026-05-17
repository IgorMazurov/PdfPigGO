package functions

import "strings"

// ParserState indicates the current parsing state of the Type 4 function tokenizer.
type ParserState int

const (
	// NewLineParserState indicates a newline character sequence was encountered.
	NewLineParserState ParserState = iota
	// WhitespaceParserState indicates whitespace characters were encountered.
	WhitespaceParserState
	// CommentParserState indicates a comment was encountered.
	CommentParserState
	// TokenParserState indicates a token (operator or value) was encountered.
	TokenParserState
)

const (
	parserNUL   = '\u0000'
	parserEOT   = '\u0004'
	parserTAB   = '\u0009'
	parserFF    = '\u000C'
	parserCR    = '\r'
	parserLF    = '\n'
	parserSPACE = '\u0020'
)

// SyntaxHandler defines all possible syntactic elements of a Type 4 function.
// It is called by the parser as the function stream is tokenized.
type SyntaxHandler interface {
	// NewLine indicates that a new line starts.
	// The text parameter contains the newline character sequence (CR, LF, CRLF, or FF).
	NewLine(text string)

	// Whitespace is called when whitespace characters are encountered.
	Whitespace(text string)

	// Token is called when a token is encountered. No distinction between operators and values is made here.
	Token(text string)

	// Comment is called for a comment (text starting with %).
	Comment(text string)
}

// DefaultSyntaxHandler provides no-op implementations for NewLine, Whitespace, and Comment,
// requiring only Token to be implemented by the embedding type.
type DefaultSyntaxHandler struct{}

func (DefaultSyntaxHandler) Comment(string)  {}
func (DefaultSyntaxHandler) NewLine(string)   {}
func (DefaultSyntaxHandler) Whitespace(string) {}

// Parse parses a Type 4 function stream and sends the syntactic elements to the given handler.
// This implements a small subset of the PostScript language but is not a full PostScript interpreter.
func Parse(input string, handler SyntaxHandler) {
	tz := newTokenizer(input, handler)
	tz.tokenize()
}

type tokenizer struct {
	input   string
	index   int
	handler SyntaxHandler
	state   ParserState
	buffer  strings.Builder
}

func newTokenizer(text string, syntaxHandler SyntaxHandler) *tokenizer {
	return &tokenizer{
		input:   text,
		handler: syntaxHandler,
		state:   WhitespaceParserState,
	}
}

func (t *tokenizer) hasMore() bool {
	return t.index < len(t.input)
}

func (t *tokenizer) currentChar() rune {
	return rune(t.input[t.index])
}

func (t *tokenizer) nextChar() rune {
	t.index++
	if !t.hasMore() {
		return parserEOT
	}
	return t.currentChar()
}

func (t *tokenizer) peek() rune {
	if t.index < len(t.input)-1 {
		return rune(t.input[t.index+1])
	}
	return parserEOT
}

func (t *tokenizer) nextState() ParserState {
	ch := t.currentChar()
	switch ch {
	case parserCR, parserLF, parserFF:
		t.state = NewLineParserState
	case parserNUL, parserTAB, parserSPACE:
		t.state = WhitespaceParserState
	case '%':
		t.state = CommentParserState
	default:
		t.state = TokenParserState
	}
	return t.state
}

func (t *tokenizer) tokenize() {
	for t.hasMore() {
		t.buffer.Reset()
		t.nextState()
		switch t.state {
		case NewLineParserState:
			t.scanNewLine()
		case WhitespaceParserState:
			t.scanWhitespace()
		case CommentParserState:
			t.scanComment()
		default:
			t.scanToken()
		}
	}
}

func (t *tokenizer) scanNewLine() {
	ch := t.currentChar()
	t.buffer.WriteRune(ch)
	if ch == parserCR && t.peek() == parserLF {
		t.buffer.WriteRune(t.nextChar())
	}
	t.handler.NewLine(t.buffer.String())
	t.nextChar()
}

func (t *tokenizer) scanWhitespace() {
	t.buffer.WriteRune(t.currentChar())

	for t.hasMore() {
		ch := t.nextChar()
		switch ch {
		case parserNUL, parserTAB, parserSPACE:
			t.buffer.WriteRune(ch)
		default:
			return
		}
	}
	t.handler.Whitespace(t.buffer.String())
}

func (t *tokenizer) scanComment() {
	t.buffer.WriteRune(t.currentChar())

	for t.hasMore() {
		ch := t.nextChar()
		switch ch {
		case parserCR, parserLF, parserFF:
			t.handler.Comment(t.buffer.String())
			return
		default:
			t.buffer.WriteRune(ch)
		}
	}
	t.handler.Comment(t.buffer.String())
}

func (t *tokenizer) scanToken() {
	ch := t.currentChar()
	t.buffer.WriteRune(ch)
	if ch == '{' || ch == '}' {
		t.handler.Token(t.buffer.String())
		t.nextChar()
		return
	}

	for t.hasMore() {
		ch = t.nextChar()
		switch ch {
		case parserNUL, parserTAB, parserSPACE, parserCR, parserLF, parserFF, '{', '}':
			t.handler.Token(t.buffer.String())
			return
		case parserEOT:
			t.handler.Token(t.buffer.String())
			return
		default:
			t.buffer.WriteRune(ch)
		}
	}
	t.handler.Token(t.buffer.String())
}
