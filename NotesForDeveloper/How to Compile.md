# See Also

- [README-Bretts.md](../README-Bretts.md)

# How to Compile This Project

This project is written in Go. You can compile it for different operating systems.

## Compile for Linux

To build the executable for Linux, run:

```bash
go build ./cmd/grab
```

This will create an executable named `grab` in the current directory.

## Compile for Windows 

**Windows compilation is now supported!** The project includes Windows-specific stub implementations that bypass tree-sitter dependencies while maintaining core functionality.

### Native Windows Compilation

If you have Go installed on Windows:

```bash
go build ./cmd/grab
```

This will create `grab.exe` in the current directory. The Windows build automatically uses stub dependency resolvers, so advanced dependency resolution features are disabled, but all core functionality works.

### Cross-Compile for Windows (from Linux)

To build a Windows executable from Linux, you need the MinGW cross-compiler:

```bash
# Install cross-compiler (Ubuntu/Debian)
sudo apt update
sudo apt install mingw-w64

# Build Windows executable
export CGO_ENABLED=1
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
go build -v -o grab.exe ./cmd/grab
```

### Windows Build Limitations

- **Dependency resolution disabled**: The `--deps` flag and <kbd>D</kbd> key in interactive mode won't parse dependencies using tree-sitter
- **Core features work**: File selection, output generation, filtering, etc. all function normally



## Compile for Windows (Recommended: Using WSL)

For development and testing, you may still want to use WSL for access to full dependency resolution features. Follow the instructions in [NotesForDeveloper/WslBuildInstructions.md](./WslBuildInstructions.md).

This document provides step-by-step guidance for compiling the project using WSL (Windows Subsystem for Linux), including all required toolchain setup and troubleshooting notes.

