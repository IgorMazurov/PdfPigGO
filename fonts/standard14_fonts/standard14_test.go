package standard14fonts

import "testing"

func TestGetNamesCount(t *testing.T) {
	names := GetNames()

	expected := 39
	if len(names) != expected {
		t.Errorf("expected %d standard 14 font names, got %d", expected, len(names))
	}
}
