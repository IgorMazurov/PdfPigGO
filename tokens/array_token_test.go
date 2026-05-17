package tokens

import "testing"

func TestArrayTokenSetsData(t *testing.T) {
	token := NewArrayToken([]Token{StartStream, EndStream})

	if got := len(token.Data()); got != 2 {
		t.Errorf("expected Data length 2, got %d", got)
	}

	if token.Data()[0] != StartStream {
		t.Error("expected Data[0] to be StartStream")
	}

	if token.Data()[1] != EndStream {
		t.Error("expected Data[1] to be EndStream")
	}
}

func TestArrayTokenSetsDataEmpty(t *testing.T) {
	token := NewArrayToken([]Token{})

	if len(token.Data()) != 0 {
		t.Errorf("expected empty Data, got length %d", len(token.Data()))
	}
}

func TestArrayTokenNilInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected NewArrayToken(nil) to panic")
		}
	}()

	NewArrayToken(nil)
}

func TestArrayTokenString(t *testing.T) {
	token := NewArrayToken([]Token{
		NewStringToken("hedgehog"),
		NewNumericToken(7),
		StartObject,
	})

	expected := "[ (hedgehog), 7, obj ]"
	if got := token.String(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
