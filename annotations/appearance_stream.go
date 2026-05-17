package annotations

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// AppearanceStream describes what an annotation looks like. Each stream is a Form XObject.
// The appearance stream is either stateless (in which case IsStateless returns true)
// or stateful, in which case the states can be retrieved via States and used to retrieve
// the state-specific appearances using Get.
type AppearanceStream struct {
	appearanceStreamsByState map[string]*tokens.StreamToken
	statelessAppearanceStream *tokens.StreamToken
}

// NewAppearanceStream creates a stateless appearance stream from the given stream token.
func NewAppearanceStream(streamToken *tokens.StreamToken) (*AppearanceStream, error) {
	if streamToken == nil {
		return nil, errors.New("stream token cannot be nil")
	}

	return &AppearanceStream{
		statelessAppearanceStream: streamToken,
	}, nil
}

// NewStatefulAppearanceStream creates a stateful appearance stream from the given map of states to streams.
func NewStatefulAppearanceStream(appearanceStreamsByState map[string]*tokens.StreamToken) (*AppearanceStream, error) {
	if appearanceStreamsByState == nil {
		return nil, errors.New("appearance streams by state cannot be nil")
	}

	return &AppearanceStream{
		appearanceStreamsByState: appearanceStreamsByState,
	}, nil
}

// IsStateless indicates whether this appearance stream is stateless or whether
// you can get appearances by state.
func (a *AppearanceStream) IsStateless() bool {
	return a.statelessAppearanceStream != nil
}

// States returns the list of available states. If this is a stateless appearance
// stream, an empty slice is returned.
func (a *AppearanceStream) States() []string {
	if a.appearanceStreamsByState == nil {
		return []string{}
	}

	states := make([]string, 0, len(a.appearanceStreamsByState))
	for k := range a.appearanceStreamsByState {
		states = append(states, k)
	}

	return states
}

// Get returns the appearance stream for the given state. Returns an error if this
// is a stateless appearance stream or if the state does not exist.
func (a *AppearanceStream) Get(state string) (*tokens.StreamToken, error) {
	if a.appearanceStreamsByState == nil {
		return nil, errors.New("cannot get appearance by state when this is a stateless appearance stream")
	}

	stream, ok := a.appearanceStreamsByState[state]
	if !ok {
		keys := make([]string, 0, len(a.appearanceStreamsByState))
		for k := range a.appearanceStreamsByState {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("appearance stream does not have state %q (available states: %s)", state, strings.Join(keys, ","))
	}

	return stream, nil
}
