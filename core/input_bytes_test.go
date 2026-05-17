package core

import (
	"bytes"
	"testing"
)

func TestArrayAndStreamBehaveTheSame(t *testing.T) {
	testData := "123456789"
	bytesData := StringAsLatin1Bytes(testData)

	array := NewMemoryInputBytes(bytesData)
	streamReader := bytes.NewReader(bytesData)
	stream, err := NewStreamInputBytes(streamReader, false)
	if err != nil {
		t.Fatalf("NewStreamInputBytes error: %v", err)
	}

	if array.Length() != int64(len(bytesData)) {
		t.Errorf("array.Length = %d, want %d", array.Length(), len(bytesData))
	}
	if stream.Length() != int64(len(bytesData)) {
		t.Errorf("stream.Length = %d, want %d", stream.Length(), len(bytesData))
	}

	if array.CurrentOffset() != 0 {
		t.Errorf("array.CurrentOffset = %d, want 0", array.CurrentOffset())
	}
	if stream.CurrentOffset() != 0 {
		t.Errorf("stream.CurrentOffset = %d, want 0", stream.CurrentOffset())
	}

	array.Seek(5, 0) // io.SeekStart
	stream.Seek(5, 0)

	if array.CurrentOffset() != stream.CurrentOffset() {
		t.Errorf("after Seek(5): array.CurrentOffset=%d, stream.CurrentOffset=%d",
			array.CurrentOffset(), stream.CurrentOffset())
	}

	if byte('5') != array.CurrentByte() {
		t.Errorf("array.CurrentByte = %c, want '5'", array.CurrentByte())
	}
	if array.CurrentByte() != stream.CurrentByte() {
		t.Errorf("array.CurrentByte=%c, stream.CurrentByte=%c",
			array.CurrentByte(), stream.CurrentByte())
	}

	arrayPeek, arrayHas := array.Peek()
	streamPeek, streamHas := stream.Peek()
	if arrayHas != streamHas || arrayPeek != streamPeek {
		t.Errorf("array.Peek=(%c,%v), stream.Peek=(%c,%v)",
			arrayPeek, arrayHas, streamPeek, streamHas)
	}

	array.Seek(0, 0)
	stream.Seek(0, 0)

	if byte(0) != array.CurrentByte() {
		t.Errorf("after Seek(0): array.CurrentByte = %c (0x%02x), want 0",
			array.CurrentByte(), array.CurrentByte())
	}
	if array.CurrentByte() != stream.CurrentByte() {
		t.Errorf("after Seek(0): array.CurrentByte=%c, stream.CurrentByte=%c",
			array.CurrentByte(), stream.CurrentByte())
	}

	array.Seek(7, 0)
	stream.Seek(7, 0)

	var arrayString string
	var streamString string

	for array.MoveNext() {
		arrayString += string(array.CurrentByte())
	}

	for stream.MoveNext() {
		streamString += string(stream.CurrentByte())
	}

	if streamString != "89" {
		t.Errorf("streamString = %q, want %q", streamString, "89")
	}

	if arrayString != streamString {
		t.Errorf("arrayString=%q, streamString=%q", arrayString, streamString)
	}

	if !stream.IsAtEnd() {
		t.Error("stream.IsAtEnd = false, want true")
	}
	if !array.IsAtEnd() {
		t.Error("array.IsAtEnd = false, want true")
	}

	stream.Seek(0, 0)
	array.Seek(0, 0)

	if stream.IsAtEnd() {
		t.Error("stream.IsAtEnd after Seek(0) = true, want false")
	}
	if array.IsAtEnd() {
		t.Error("array.IsAtEnd after Seek(0) = true, want false")
	}
}

func TestReadFromBeginningIsCorrect(t *testing.T) {
	inputBytes := NewMemoryInputBytes(StringAsLatin1Bytes("endstream and then <</go[]>>"))

	buffer := make([]byte, len("endstream"))
	result, err := inputBytes.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if result != len(buffer) {
		t.Errorf("Read returned %d, want %d", result, len(buffer))
	}
	if BytesAsLatin1String(buffer) != "endstream" {
		t.Errorf("buffer = %q, want %q", BytesAsLatin1String(buffer), "endstream")
	}

	if inputBytes.CurrentByte() != byte('m') {
		t.Errorf("CurrentByte = %c, want 'm'", inputBytes.CurrentByte())
	}
	if !inputBytes.MoveNext() {
		t.Error("MoveNext returned false, want true")
	}
	if !inputBytes.MoveNext() {
		t.Error("second MoveNext returned false, want true")
	}
	if inputBytes.CurrentByte() != byte('a') {
		t.Errorf("CurrentByte = %c, want 'a'", inputBytes.CurrentByte())
	}
}

func TestReadMatchesMoveBehaviour(t *testing.T) {
	bytesRead := NewMemoryInputBytes(StringAsLatin1Bytes("cows in the south"))
	bytesMove := NewMemoryInputBytes(StringAsLatin1Bytes("cows in the north"))

	readLength := 3
	buffer := make([]byte, readLength)

	readResult, err := bytesRead.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	for i := 0; i < readLength; i++ {
		bytesMove.MoveNext()
	}

	if readResult != readLength {
		t.Errorf("Read returned %d, want %d", readResult, readLength)
	}

	if bytesRead.CurrentOffset() != bytesMove.CurrentOffset() {
		t.Errorf("bytesRead.CurrentOffset=%d, bytesMove.CurrentOffset=%d",
			bytesRead.CurrentOffset(), bytesMove.CurrentOffset())
	}
	if bytesRead.CurrentByte() != bytesMove.CurrentByte() {
		t.Errorf("bytesRead.CurrentByte=%c, bytesMove.CurrentByte=%c",
			bytesRead.CurrentByte(), bytesMove.CurrentByte())
	}
	if bytesRead.MoveNext() != bytesMove.MoveNext() {
		t.Errorf("bytesRead.MoveNext=%v, bytesMove.MoveNext=%v",
			bytesRead.MoveNext(), bytesMove.MoveNext())
	}
	if bytesRead.CurrentOffset() != bytesMove.CurrentOffset() {
		t.Errorf("after MoveNext: bytesRead.CurrentOffset=%d, bytesMove.CurrentOffset=%d",
			bytesRead.CurrentOffset(), bytesMove.CurrentOffset())
	}
	if bytesRead.CurrentByte() != bytesMove.CurrentByte() {
		t.Errorf("after MoveNext: bytesRead.CurrentByte=%c, bytesMove.CurrentByte=%c",
			bytesRead.CurrentByte(), bytesMove.CurrentByte())
	}
}

func TestReadFromMiddleIsCorrect(t *testing.T) {
	inputBytes := NewMemoryInputBytes(StringAsLatin1Bytes("aa stream <<>>"))

	if !inputBytes.MoveNext() {
		t.Error("first MoveNext returned false")
	}
	if !inputBytes.MoveNext() {
		t.Error("second MoveNext returned false")
	}
	if !inputBytes.MoveNext() {
		t.Error("third MoveNext returned false")
	}

	if inputBytes.CurrentByte() != byte(' ') {
		t.Errorf("CurrentByte = %c, want ' '", inputBytes.CurrentByte())
	}

	buffer := make([]byte, len("stream"))
	result, err := inputBytes.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if result != len(buffer) {
		t.Errorf("Read returned %d, want %d", result, len(buffer))
	}
	if BytesAsLatin1String(buffer) != "stream" {
		t.Errorf("buffer = %q, want %q", BytesAsLatin1String(buffer), "stream")
	}

	if inputBytes.CurrentByte() != byte('m') {
		t.Errorf("CurrentByte = %c, want 'm'", inputBytes.CurrentByte())
	}
	if !inputBytes.MoveNext() {
		t.Error("MoveNext returned false")
	}
	if !inputBytes.MoveNext() {
		t.Error("second MoveNext returned false")
	}
	if inputBytes.CurrentByte() != byte('<') {
		t.Errorf("CurrentByte = %c, want '<'", inputBytes.CurrentByte())
	}
}

func TestReadPastEndIsCorrect(t *testing.T) {
	inputBytes := NewMemoryInputBytes(StringAsLatin1Bytes("stream"))

	if !inputBytes.MoveNext() {
		t.Error("first MoveNext returned false")
	}
	if !inputBytes.MoveNext() {
		t.Error("second MoveNext returned false")
	}

	buffer := make([]byte, len("stream"))
	result, err := inputBytes.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	expectedResult := len(buffer) - 2
	if result != expectedResult {
		t.Errorf("Read returned %d, want %d", result, expectedResult)
	}
	readContent := BytesAsLatin1String(buffer[:result])
	if readContent != "ream" {
		t.Errorf("read content = %q, want %q", readContent, "ream")
	}

	if inputBytes.CurrentByte() != byte('m') {
		t.Errorf("CurrentByte = %c, want 'm'", inputBytes.CurrentByte())
	}
	if !inputBytes.IsAtEnd() {
		t.Error("IsAtEnd = false, want true")
	}
	if inputBytes.MoveNext() {
		t.Error("MoveNext after end returned true, want false")
	}
}

func TestReadFromStreamBeginningIsCorrect(t *testing.T) {
	streamReader := bytes.NewReader(StringAsLatin1Bytes("endstream and then <</go[]>>"))
	inputStream, err := NewStreamInputBytes(streamReader, false)
	if err != nil {
		t.Fatalf("NewStreamInputBytes error: %v", err)
	}

	buffer := make([]byte, len("endstream"))
	result, err := inputStream.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if result != len(buffer) {
		t.Errorf("Read returned %d, want %d", result, len(buffer))
	}
	if BytesAsLatin1String(buffer) != "endstream" {
		t.Errorf("buffer = %q, want %q", BytesAsLatin1String(buffer), "endstream")
	}

	if inputStream.CurrentByte() != byte('m') {
		t.Errorf("CurrentByte = %c, want 'm'", inputStream.CurrentByte())
	}
	if !inputStream.MoveNext() {
		t.Error("MoveNext returned false, want true")
	}
	if !inputStream.MoveNext() {
		t.Error("second MoveNext returned false, want true")
	}
	if inputStream.CurrentByte() != byte('a') {
		t.Errorf("CurrentByte = %c, want 'a'", inputStream.CurrentByte())
	}
}

func TestReadFromStreamMiddleIsCorrect(t *testing.T) {
	streamReader := bytes.NewReader(StringAsLatin1Bytes("aa stream <<>>"))
	inputStream, err := NewStreamInputBytes(streamReader, false)
	if err != nil {
		t.Fatalf("NewStreamInputBytes error: %v", err)
	}

	if !inputStream.MoveNext() {
		t.Error("first MoveNext returned false")
	}
	if !inputStream.MoveNext() {
		t.Error("second MoveNext returned false")
	}
	if !inputStream.MoveNext() {
		t.Error("third MoveNext returned false")
	}

	if inputStream.CurrentByte() != byte(' ') {
		t.Errorf("CurrentByte = %c, want ' '", inputStream.CurrentByte())
	}

	buffer := make([]byte, len("stream"))
	result, err := inputStream.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if result != len(buffer) {
		t.Errorf("Read returned %d, want %d", result, len(buffer))
	}
	if BytesAsLatin1String(buffer) != "stream" {
		t.Errorf("buffer = %q, want %q", BytesAsLatin1String(buffer), "stream")
	}

	if inputStream.CurrentByte() != byte('m') {
		t.Errorf("CurrentByte = %c, want 'm'", inputStream.CurrentByte())
	}
	if !inputStream.MoveNext() {
		t.Error("MoveNext returned false")
	}
	if !inputStream.MoveNext() {
		t.Error("second MoveNext returned false")
	}
	if inputStream.CurrentByte() != byte('<') {
		t.Errorf("CurrentByte = %c, want '<'", inputStream.CurrentByte())
	}
}

func TestReadPastStreamEndIsCorrect(t *testing.T) {
	streamReader := bytes.NewReader(StringAsLatin1Bytes("stream"))
	inputStream, err := NewStreamInputBytes(streamReader, false)
	if err != nil {
		t.Fatalf("NewStreamInputBytes error: %v", err)
	}

	if !inputStream.MoveNext() {
		t.Error("first MoveNext returned false")
	}
	if !inputStream.MoveNext() {
		t.Error("second MoveNext returned false")
	}

	buffer := make([]byte, len("stream"))
	result, err := inputStream.Read(buffer)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	expectedResult := len(buffer) - 2
	if result != expectedResult {
		t.Errorf("Read returned %d, want %d", result, expectedResult)
	}
	readContent := BytesAsLatin1String(buffer[:result])
	if readContent != "ream" {
		t.Errorf("read content = %q, want %q", readContent, "ream")
	}

	if inputStream.CurrentByte() != byte('m') {
		t.Errorf("CurrentByte = %c, want 'm'", inputStream.CurrentByte())
	}
	if !inputStream.IsAtEnd() {
		t.Error("IsAtEnd = false, want true")
	}
	if inputStream.MoveNext() {
		t.Error("MoveNext after end returned true, want false")
	}
}
