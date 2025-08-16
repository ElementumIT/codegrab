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

## Compile for Windows (Cross-Compile from Linux)

> **Note:** As of now, this project cannot be compiled for Windows due to build constraints in the `go-tree-sitter` dependency. Some language bindings (Go, Python, TypeScript/TSX) are not supported on Windows and will cause build errors. See the error message below:
>
> ```
> github.com/smacker/go-tree-sitter/golang: build constraints exclude all Go files ...
> # github.com/smacker/go-tree-sitter
> ... undefined: Node
> ```
>
> If you need Windows support, you will need to remove or replace these dependencies, or use WSL/Docker on Windows.



## Compile for Windows (Recommended: Using WSL)

To build for Windows, follow the instructions in [NotesForDeveloper/WslBuildInstructions.md](./WslBuildInstructions.md).

This document provides step-by-step guidance for compiling the project for Windows using WSL (Windows Subsystem for Linux), including all required toolchain setup and troubleshooting notes.

