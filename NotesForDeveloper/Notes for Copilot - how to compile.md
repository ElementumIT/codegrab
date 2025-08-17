# Notes for Copilot - How to Compile CodeGrab

This document serves as a reference for Copilot on how to compile CodeGrab for different platforms, including the Windows compilation solution implemented.

## Overview

CodeGrab is a Go project that has dependencies on tree-sitter parsers for dependency resolution. The main challenge for Windows compilation was the CGO-dependent tree-sitter libraries that don't compile properly on Windows.

## Solution Implemented

**Windows Stub Resolvers**: Created Windows-specific stub implementations for dependency resolvers that bypass tree-sitter entirely. This allows Windows compilation while disabling advanced dependency resolution features.

### Files Created for Windows Support

1. `internal/dependencies/go_resolver_windows.go` - Windows stub for Go dependency resolution
2. `internal/dependencies/js_resolver_windows.go` - Windows stub for JS/TS dependency resolution  
3. `internal/dependencies/py_resolver_windows.go` - Windows stub for Python dependency resolution

Each stub resolver:
- Uses `//go:build windows` build constraint
- Returns empty dependency lists instead of parsing with tree-sitter
- Allows the main application to compile and run on Windows

The original resolvers use `//go:build !windows` to exclude them from Windows builds.

## Compilation Instructions

### Linux/WSL Compilation

**Standard build:**
```bash
go build ./cmd/grab
```

**With verbose output:**
```bash
go build -v ./cmd/grab
```

### Windows Compilation (Native)

**From Windows with Go installed:**
```bash
go build ./cmd/grab
```

The Windows stub resolvers will automatically be used, bypassing tree-sitter dependencies.

### Cross-Compilation for Windows (from Linux/WSL)

**Prerequisites:**
- MinGW-w64 cross-compiler: `sudo apt install mingw-w64`

**Build Windows executable:**
```bash
export CGO_ENABLED=1
export GOOS=windows  
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
go build -v -o grab.exe ./cmd/grab
```

**Alternative single command:**
```bash
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -v -o grab.exe ./cmd/grab
```

### Verification

**Check Windows executable:**
```bash
file grab.exe  # Should show "PE32+ executable"
```

**Test basic functionality:**
```bash
./grab.exe --version  # Should display version information
```

## Key Technical Details

### Build Constraints Used

- `//go:build !windows` - Excludes files from Windows builds
- `//go:build windows` - Includes files only in Windows builds

### Tree-sitter Dependency Impact

**On Windows:**
- Dependency resolution is disabled (returns empty lists)
- Main functionality (file selection, output generation) works normally
- No tree-sitter parsing capabilities

**On Linux/Unix:**
- Full dependency resolution available
- Tree-sitter parsing works for Go, JS/TS, Python
- Complete feature set available

### Output Location

Windows executables are typically placed in:
- `CompiledForWindows/grab.exe` - For distribution
- `grab.exe` - Local builds

## Troubleshooting

### Common Windows Compilation Issues

1. **CGO not enabled**: Ensure `CGO_ENABLED=1` is set
2. **Missing cross-compiler**: Install `mingw-w64` package
3. **Wrong CC variable**: Use `x86_64-w64-mingw32-gcc` for 64-bit Windows
4. **Path issues**: Use forward slashes in paths even on Windows

### Build Verification Commands

```bash
# Check Go environment
go env GOOS GOARCH CGO_ENABLED

# Verify cross-compiler
x86_64-w64-mingw32-gcc --version

# Clean build cache if needed
go clean -cache
go clean -modcache
```

## Integration Notes

This solution maintains:
- Cross-platform compatibility
- All non-dependency-resolution features on Windows
- Full functionality on Unix-like systems
- Clean separation of platform-specific code

The Windows build sacrifices advanced dependency parsing for basic compilation compatibility, which is acceptable for the core use case of file selection and bundling.