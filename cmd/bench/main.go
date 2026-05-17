// Package main provides a convenience entry point for running PdfPigGO benchmarks.
//
// Go does not use a dedicated benchmark runner program like BenchmarkDotNet in C#.
// Instead, benchmarks are defined as _test.go files using testing.B and executed
// via the `go test` command with the -bench flag.
//
// To run all benchmarks:
//
//	go test -bench=. ./benchmarks/...
//
// To run a specific benchmark (e.g., BenchmarkReadPdf):
//
//	go test -bench=BenchmarkReadPdf ./benchmarks/...
//
// To run with custom iterations:
//
//	go test -bench=. -benchtime=10s ./benchmarks/...
//
// For memory allocation reports:
//
//	go test -bench=. -benchmem ./benchmarks/...
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, `PdfPigGO Benchmark Runner
========================

Go benchmarks are run via the "go test" command, not a standalone program.

Usage:
  go test -bench=. ./benchmarks/...

Examples:
  go test -bench=.                    ./benchmarks/...   # Run all benchmarks
  go test -bench=BenchmarkReadPdf     ./benchmarks/...   # Run specific benchmark
  go test -bench=. -benchtime=10s     ./benchmarks/...   # Custom duration
  go test -bench=. -benchmem          ./benchmarks/...   # Include memory stats`)
	os.Exit(0)
}
