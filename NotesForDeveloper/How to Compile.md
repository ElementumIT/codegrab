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

To build a Windows executable from Linux (if dependencies are fixed), use:

```bash
GOOS=windows GOARCH=amd64 go build -o grab.exe ./cmd/grab
```

> **Important:** This command will only work if all dependencies support Windows. Currently, it will fail unless you are running in WSL (Windows Subsystem for Linux) or using Docker with a Linux container, and the dependencies are compatible.

This will produce a file called `grab.exe` that can be run on Windows 64-bit systems (if the build succeeds).

## Notes

- You do **not** need to install Go on Windows to run the compiled `.exe` file.
- If you want to compile for other platforms, change the `GOOS` and `GOARCH` variables accordingly.
- For more details, see the [Go documentation on cross-compilation](https://golang.org/doc/install/source#environment).
