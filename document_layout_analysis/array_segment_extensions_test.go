package document_layout_analysis_test

import (
	"testing"

	dla "github.com/uglytoad/pdfpig/go/document_layout_analysis"
)

func TestArraySegmentTakeGetAt(t *testing.T) {
	arr := make([]int, 10)
	for i := range arr {
		arr[i] = i
	}
	array := dla.NewArraySegment(arr, 0, len(arr))

	if array.Count != 10 {
		t.Errorf("expected count 10, got %d", array.Count)
	}

	first5 := dla.ArraySegmentTake(array, 5)
	if first5.Count != 5 {
		t.Errorf("expected take(5) count 5, got %d", first5.Count)
	}

	for i := 0; i < 5; i++ {
		v, err := dla.ArraySegmentGetAt(first5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != i {
			t.Errorf("expected first5[%d] == %d, got %d", i, i, v)
		}
	}

	first2of5 := dla.ArraySegmentTake(first5, 2)
	if first2of5.Count != 2 {
		t.Errorf("expected take(2) count 2, got %d", first2of5.Count)
	}

	for i := 0; i < 2; i++ {
		v, err := dla.ArraySegmentGetAt(first2of5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != i {
			t.Errorf("expected first2of5[%d] == %d, got %d", i, i, v)
		}
	}
}

func TestArraySegmentSkipGetAt(t *testing.T) {
	arr := make([]int, 10)
	for i := range arr {
		arr[i] = i
	}
	array := dla.NewArraySegment(arr, 0, len(arr))

	if array.Count != 10 {
		t.Errorf("expected count 10, got %d", array.Count)
	}

	skip5 := dla.ArraySegmentSkip(array, 5)
	if skip5.Count != 5 {
		t.Errorf("expected skip(5) count 5, got %d", skip5.Count)
	}

	expectedSkip5 := []int{5, 6, 7, 8, 9}
	for i, expected := range expectedSkip5 {
		v, err := dla.ArraySegmentGetAt(skip5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected skip5[%d] == %d, got %d", i, expected, v)
		}
	}

	skip2of5 := dla.ArraySegmentSkip(skip5, 2)
	if skip2of5.Count != 3 {
		t.Errorf("expected skip(2) count 3, got %d", skip2of5.Count)
	}

	expectedSkip2of5 := []int{7, 8, 9}
	for i, expected := range expectedSkip2of5 {
		v, err := dla.ArraySegmentGetAt(skip2of5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected skip2of5[%d] == %d, got %d", i, expected, v)
		}
	}
}

func TestArraySegmentSkipTakeGetAt(t *testing.T) {
	arr := make([]int, 10)
	for i := range arr {
		arr[i] = i
	}
	array := dla.NewArraySegment(arr, 0, len(arr))

	if array.Count != 10 {
		t.Errorf("expected count 10, got %d", array.Count)
	}

	skip5 := dla.ArraySegmentSkip(array, 5)
	if skip5.Count != 5 {
		t.Errorf("expected skip(5) count 5, got %d", skip5.Count)
	}

	expectedSkip5 := []int{5, 6, 7, 8, 9}
	for i, expected := range expectedSkip5 {
		v, err := dla.ArraySegmentGetAt(skip5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected skip5[%d] == %d, got %d", i, expected, v)
		}
	}

	take2ofSkip5 := dla.ArraySegmentTake(skip5, 2)
	if take2ofSkip5.Count != 2 {
		t.Errorf("expected take(2) count 2, got %d", take2ofSkip5.Count)
	}

	v0, err := dla.ArraySegmentGetAt(take2ofSkip5, 0)
	if err != nil {
		t.Fatalf("unexpected error at index 0: %v", err)
	}
	if v0 != 5 {
		t.Errorf("expected take2ofSkip5[0] == 5, got %d", v0)
	}

	v1, err := dla.ArraySegmentGetAt(take2ofSkip5, 1)
	if err != nil {
		t.Fatalf("unexpected error at index 1: %v", err)
	}
	if v1 != 6 {
		t.Errorf("expected take2ofSkip5[1] == 6, got %d", v1)
	}

	v2, err := dla.ArraySegmentGetAt(take2ofSkip5, 2)
	if err != nil {
		t.Fatalf("unexpected error at index 2: %v", err)
	}
	if v2 != 7 {
		t.Errorf("expected take2ofSkip5[2] == 7, got %d", v2)
	}
}

func TestArraySegmentTakeSkipGetAt(t *testing.T) {
	arr := make([]int, 10)
	for i := range arr {
		arr[i] = i
	}
	array := dla.NewArraySegment(arr, 0, len(arr))

	if array.Count != 10 {
		t.Errorf("expected count 10, got %d", array.Count)
	}

	first5 := dla.ArraySegmentTake(array, 5)
	if first5.Count != 5 {
		t.Errorf("expected take(5) count 5, got %d", first5.Count)
	}

	for i := 0; i < 5; i++ {
		v, err := dla.ArraySegmentGetAt(first5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != i {
			t.Errorf("expected first5[%d] == %d, got %d", i, i, v)
		}
	}

	skip2ofTake5 := dla.ArraySegmentSkip(first5, 2)
	if skip2ofTake5.Count != 3 {
		t.Errorf("expected skip(2) count 3, got %d", skip2ofTake5.Count)
	}

	expectedSkip2 := []int{2, 3, 4}
	for i, expected := range expectedSkip2 {
		v, err := dla.ArraySegmentGetAt(skip2ofTake5, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected skip2ofTake5[%d] == %d, got %d", i, expected, v)
		}
	}
}

func TestArraySegmentSort(t *testing.T) {
	arr := make([]int, 10)
	for i := range arr {
		arr[i] = i
	}

	array := dla.NewArraySegment(arr, 0, len(arr))
	if array.Count != 10 {
		t.Fatalf("expected count 10, got %d", array.Count)
	}

	dla.ArraySegmentSort(&array, func(a, b int) bool { return a > b })

	if array.Count != 10 {
		t.Fatalf("expected count 10 after sort, got %d", array.Count)
	}

	expectedDesc := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	for i, expected := range expectedDesc {
		v, err := dla.ArraySegmentGetAt(array, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected array[%d] == %d, got %d", i, expected, v)
		}
	}

	skip1Take7 := dla.ArraySegmentTake(dla.ArraySegmentSkip(array, 1), 7)
	if skip1Take7.Count != 7 {
		t.Fatalf("expected skip(1).take(7) count 7, got %d", skip1Take7.Count)
	}

	expectedSeg := []int{8, 7, 6, 5, 4, 3, 2}
	for i, expected := range expectedSeg {
		v, err := dla.ArraySegmentGetAt(skip1Take7, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected skip1Take7[%d] == %d, got %d", i, expected, v)
		}
	}

	dla.ArraySegmentSort(&skip1Take7, func(a, b int) bool { return a < b })

	if skip1Take7.Count != 7 {
		t.Fatalf("expected count 7 after sort segment, got %d", skip1Take7.Count)
	}

	expectedAsc := []int{2, 3, 4, 5, 6, 7, 8}
	for i, expected := range expectedAsc {
		v, err := dla.ArraySegmentGetAt(skip1Take7, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected skip1Take7[%d] == %d, got %d", i, expected, v)
		}
	}

	if array.Count != 10 {
		t.Fatalf("expected parent count 10, got %d", array.Count)
	}

	expectedParent := []int{9, 2, 3, 4, 5, 6, 7, 8, 1, 0}
	for i, expected := range expectedParent {
		v, err := dla.ArraySegmentGetAt(array, i)
		if err != nil {
			t.Fatalf("unexpected error at index %d: %v", i, err)
		}
		if v != expected {
			t.Errorf("expected array[%d] == %d, got %d", i, expected, v)
		}
	}

	expectedOrig := []int{9, 2, 3, 4, 5, 6, 7, 8, 1, 0}
	for i, expected := range expectedOrig {
		if arr[i] != expected {
			t.Errorf("expected originalArray[%d] == %d, got %d", i, expected, arr[i])
		}
	}
}
