package tokens

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/core"
)

// DictEntry holds a single key-value pair for ordered iteration.
type DictEntry struct {
	Key   string
	Value Token
}

// internalDictEntry holds a single key-value pair for ordered iteration (internal).
type internalDictEntry struct {
	key   string
	value Token
}

// DictionaryToken represents a PDF dictionary object, an associative table
// containing pairs of objects known as the dictionary's entries. The key must
// be a name and the value may be any kind of token.
type DictionaryToken struct {
	data    map[string]Token
	ordered []internalDictEntry
}

var _ Token = (*DictionaryToken)(nil)
var _ DataToken[map[string]Token] = (*DictionaryToken)(nil)

// NewDictionary creates a new DictionaryToken from a map keyed by NameToken.
func NewDictionary(data map[*NameToken]Token) (*DictionaryToken, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	result := make(map[string]Token, len(data))
	ordered := make([]internalDictEntry, 0, len(data))

	for key, value := range data {
		k := key.Data()
		result[k] = value
		ordered = append(ordered, internalDictEntry{key: k, value: value})
	}

	return &DictionaryToken{data: result, ordered: ordered}, nil
}

// WithMap creates a new DictionaryToken from a string-keyed map.
func WithMap(data map[string]Token) (*DictionaryToken, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	result := make(map[string]Token, len(data))
	ordered := make([]internalDictEntry, 0, len(data))
	for k, v := range data {
		result[k] = v
		ordered = append(ordered, internalDictEntry{key: k, value: v})
	}

	return &DictionaryToken{data: result, ordered: ordered}, nil
}

// NewOrderedDictionary creates a new DictionaryToken preserving insertion order.
func NewOrderedDictionary(entries []DictEntry) (*DictionaryToken, error) {
	if entries == nil {
		entries = make([]DictEntry, 0)
	}

	result := make(map[string]Token, len(entries))
	ordered := make([]internalDictEntry, 0, len(entries))
	for _, e := range entries {
		result[e.Key] = e.Value
		ordered = append(ordered, internalDictEntry{key: e.Key, value: e.Value})
	}

	return &DictionaryToken{data: result, ordered: ordered}, nil
}

// OrderedEntries returns the dictionary entries in insertion order.
func (d *DictionaryToken) OrderedEntries() []DictEntry {
	result := make([]DictEntry, 0, len(d.ordered))
	for _, e := range d.ordered {
		result = append(result, DictEntry{Key: e.key, Value: e.value})
	}
	return result
}

// Data returns the key-value pairs in this dictionary.
func (d *DictionaryToken) Data() map[string]Token {
	return d.data
}

// TryGet retrieves the entry with the given name.
// Returns true if the token is found, false otherwise.
// Panics if name is nil, matching C# ArgumentNullException behavior.
func (d *DictionaryToken) TryGet(name *NameToken) (Token, bool) {
	if name == nil {
		panic(fmt.Errorf("name cannot be null"))
	}

	token, ok := d.data[name.Data()]
	return token, ok
}

// TryGetTyped retrieves the entry with the given name and casts it to the
// expected type T. Returns true only if the token exists and matches the
// requested type, false otherwise. This corresponds to C# TryGet[T].
// Note: Go does not support generic methods on structs, so this is a
// standalone function rather than a receiver method.
func TryGetTyped[T any](d *DictionaryToken, name *NameToken) (T, bool) {
	var zero T
	if d == nil || name == nil {
		return zero, false
	}

	token, ok := d.data[name.Data()]
	if !ok {
		return zero, false
	}

	typed, ok := token.(T)
	return typed, ok
}

// GetTyped retrieves the entry with the given name and casts it to type T.
// Returns an error if the key doesn't exist or the value isn't of type T.
// This corresponds to C# Get[T] which throws PdfDocumentFormatException on mismatch.
func GetTyped[T any](d *DictionaryToken, name *NameToken) (T, error) {
	var zero T
	if d == nil {
		return zero, fmt.Errorf("dictionary cannot be null")
	}
	if name == nil {
		return zero, fmt.Errorf("name cannot be null")
	}

	token, ok := d.data[name.Data()]
	if !ok {
		return zero, core.NewPdfDocumentFormatException(
			fmt.Sprintf("dictionary does not contain token with name %s", name.Data()))
	}

	typed, ok := token.(T)
	if !ok {
		return zero, core.NewPdfDocumentFormatException(
			fmt.Sprintf("dictionary does not contain token with name %s of type %T, found %T instead", name.Data(), zero, token))
	}

	return typed, nil
}

// ContainsKey reports whether the dictionary contains an entry with the given name.
func (d *DictionaryToken) ContainsKey(name *NameToken) bool {
	if name == nil {
		return false
	}

	_, ok := d.data[name.Data()]
	return ok
}

// With creates a copy of this dictionary with the additional entry added or
// the existing entry overridden.
func (d *DictionaryToken) With(key *NameToken, value Token) *DictionaryToken {
	if key == nil {
		return d
	}

	return d.WithString(key.Data(), value)
}

// WithString creates a copy of this dictionary with the additional entry added
// or the existing entry overridden using a string key.
func (d *DictionaryToken) WithString(key string, value Token) *DictionaryToken {
	if key == "" || value == nil {
		return d
	}

	result := make(map[string]Token, len(d.data)+1)
	ordered := make([]internalDictEntry, 0, len(d.ordered)+1)
	exists := false

	for _, e := range d.ordered {
		if e.key == key {
			result[e.key] = value
			ordered = append(ordered, internalDictEntry{key: e.key, value: value})
			exists = true
		} else {
			result[e.key] = e.value
			ordered = append(ordered, e)
		}
	}

	if !exists {
		result[key] = value
		ordered = append(ordered, internalDictEntry{key: key, value: value})
	}

	return &DictionaryToken{data: result, ordered: ordered}
}

// Without creates a copy of this dictionary with the entry having the specified
// key removed, if it exists.
func (d *DictionaryToken) Without(key *NameToken) *DictionaryToken {
	if key == nil {
		return d
	}

	return d.WithoutString(key.Data())
}

// WithoutString creates a copy of this dictionary with the entry having the
// specified string key removed, if it exists.
func (d *DictionaryToken) WithoutString(key string) *DictionaryToken {
	if key == "" {
		return d
	}

	count := len(d.data)
	if _, ok := d.data[key]; ok {
		count--
	}

	result := make(map[string]Token, count)
	ordered := make([]internalDictEntry, 0, count)

	for _, e := range d.ordered {
		if e.key != key {
			result[e.key] = e.value
			ordered = append(ordered, e)
		}
	}

	return &DictionaryToken{data: result, ordered: ordered}
}

// Equals reports whether other is a DictionaryToken with equivalent entries.
func (d *DictionaryToken) Equals(other Token) bool {
	o, ok := other.(*DictionaryToken)
	if !ok || o == nil {
		return false
	}

	if d == o {
		return true
	}

	if len(d.data) != len(o.data) {
		return false
	}

	for k, v := range o.data {
		val, ok := d.data[k]
		if !ok {
			return false
		}

		eq, ok := val.(interface{ Equals(Token) bool })
		if !ok || !eq.Equals(v) {
			return false
		}
	}

	return true
}

// String returns the string representation of the dictionary token.
func (d *DictionaryToken) String() string {
	if d == nil {
		return ""
	}

	parts := make([]string, 0, len(d.data))
	for k, v := range d.data {
		parts = append(parts, fmt.Sprintf("<%s, %v>", k, v))
	}

	return strings.Join(parts, ", ")
}
