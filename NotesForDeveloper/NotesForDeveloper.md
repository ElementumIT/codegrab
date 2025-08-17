# Developer Notes

## Build Commands

**Linux/WSL:**
```bash
go build ./cmd/grab
```

**Windows (native):**
```bash
go build .\cmd\grab
```

**Cross-compile for Windows (from Linux):**
```bash
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -v -o grab.exe ./cmd/grab
```

See [How to Compile.md](How%20to%20Compile.md) for detailed compilation instructions.
