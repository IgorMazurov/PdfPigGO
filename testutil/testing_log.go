package testutil

import "github.com/uglytoad/pdfpig/go/logging"

// TestingLog is a no-operation implementation of logging.Log used in tests
// to suppress log output. All methods are empty stubs that discard messages.
type TestingLog struct{}

var _ logging.Log = TestingLog{}

// Debug implements logging.Log.Debug as a no-op.
func (TestingLog) Debug(_ string) {}

// DebugWithException implements logging.Log.DebugWithException as a no-op.
func (TestingLog) DebugWithException(_ string, _ error) {}

// Warn implements logging.Log.Warn as a no-op.
func (TestingLog) Warn(_ string) {}

// Error implements logging.Log.Error as a no-op.
func (TestingLog) Error(_ string) {}

// ErrorWithException implements logging.Log.ErrorWithException as a no-op.
func (TestingLog) ErrorWithException(_ string, _ error) {}
