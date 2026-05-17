package logging

// Log logs internal messages from the PDF parsing process. Consumers can provide their own implementation
// in the ParsingOptions to intercept log messages.
type Log interface {
	// Debug records an informational debug message.
	Debug(message string)

	// DebugWithException records an informational debug message with error.
	DebugWithException(message string, err error)

	// Warn records a warning message due to a non-error issue encountered in parsing.
	Warn(message string)

	// Error logs an error message due to an issue encountered in parsing.
	Error(message string)

	// ErrorWithException logs an error message due to an issue encountered in parsing with error.
	ErrorWithException(message string, err error)
}

// noopLog is a no-operation implementation of Log that discards all messages.
type noopLog struct{}

// Debug implements Log.Debug as a no-op.
func (noopLog) Debug(_ string) {}

// DebugWithException implements Log.DebugWithException as a no-op.
func (noopLog) DebugWithException(_ string, _ error) {}

// Warn implements Log.Warn as a no-op.
func (noopLog) Warn(_ string) {}

// Error implements Log.Error as a no-op.
func (noopLog) Error(_ string) {}

// ErrorWithException implements Log.ErrorWithException as a no-op.
func (noopLog) ErrorWithException(_ string, _ error) {}

// NoopLog is the package-level singleton instance of a no-operation logger.
var NoopLog Log = noopLog{}
