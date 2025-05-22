# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go project called `go-flock`. The codebase is currently being initialized.

## Development Commands

Since this is a new Go project, here are the standard commands that will be commonly used:

### Building and Running
```bash
go build ./...          # Build all packages
go run .                # Run main package
go install ./...        # Install binaries
```

### Testing
```bash
go test ./...           # Run all tests
go test -v ./...        # Run tests with verbose output
go test ./pkg/...       # Run tests in specific package
go test -run TestName   # Run specific test
go test -bench=.        # Run benchmarks
```

### Code Quality
```bash
go fmt ./...            # Format all Go files
go vet ./...            # Run go vet static analysis
golangci-lint run       # Run comprehensive linting (if configured)
go mod tidy             # Clean up go.mod and go.sum
```

## Architecture Notes

As the project develops, architectural patterns and key design decisions should be documented here.