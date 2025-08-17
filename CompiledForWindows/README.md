# CompiledForWindows Directory

This directory is for Windows executable builds but **binaries are not committed to Git** for security reasons.

## Why binaries are not committed

The Go application uses the [gitleaks](https://github.com/zricethezav/gitleaks) library for secret scanning, which includes example/test API keys in its embedded rules. When compiled into a binary, these test keys are embedded and can trigger GitHub's secret scanning alerts.

## How to build Windows executables

See [NotesForDeveloper/Notes for Copilot - how to compile.md](../NotesForDeveloper/Notes%20for%20Copilot%20-%20how%20to%20compile.md) for detailed compilation instructions.

Quick build commands:
```bash
# For Windows (cross-compile from Linux)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o CompiledForWindows/codegrab.exe ./cmd/codegrab

# Debug version with logging
go build -ldflags="-s -w" -tags debug -o CompiledForWindows/codegrab_debug.exe ./cmd/codegrab
```

## Important Notes

- Generated `.exe` files are automatically ignored by `.gitignore` 
- The Windows tree rendering fix has been implemented in the source code
- Users should build fresh binaries locally or download from GitHub releases