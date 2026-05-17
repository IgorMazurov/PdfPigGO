package cmap

import "testing"

func TestCidRangeEndCannotBeLowerThanStart(t *testing.T) {
	_, err := NewCidRange(0, -56, 1)
	if err == nil {
		t.Error("expected error when lastCharacterCode < firstCharacterCode")
	}
}

func TestCidRangeContainsFalseForLowerNumber(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 0)

	if rangeObj.Contains(-12) {
		t.Error("expected Contains(-12) to be false")
	}
}

func TestCidRangeContainsFalseForHigherNumber(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 0)

	if rangeObj.Contains(100) {
		t.Error("expected Contains(100) to be false")
	}
}

func TestCidRangeContainsTrueForNumberInRange(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 0)

	if !rangeObj.Contains(52) {
		t.Error("expected Contains(52) to be true")
	}
}

func TestCidRangeTryMapFalseForNumberLowerThanRange(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 0)

	if _, ok := rangeObj.TryMap(-12); ok {
		t.Error("expected TryMap(-12) to return false")
	}
}

func TestCidRangeTryMapFalseForNumberHigherThanRange(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 0)

	if _, ok := rangeObj.TryMap(250); ok {
		t.Error("expected TryMap(250) to return false")
	}
}

func TestCidRangeTryMapMapsCorrectlyForNumberInRange(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 0)

	cid, ok := rangeObj.TryMap(52)
	if !ok {
		t.Fatal("expected TryMap(52) to return true")
	}

	if cid != 52 {
		t.Errorf("expected CID 52, got %d", cid)
	}
}

func TestCidRangeTryMapMapsCorrectlyForNumberInRangeWithCidOffset(t *testing.T) {
	rangeObj, _ := NewCidRange(0, 69, 9)

	cid, ok := rangeObj.TryMap(52)
	if !ok {
		t.Fatal("expected TryMap(52) to return true")
	}

	if cid != 61 {
		t.Errorf("expected CID 61, got %d", cid)
	}
}
