# Building codegrab in WSL (Windows Subsystem for Linux)

## Prerequisites
- WSL2 installed and running (Ubuntu recommended)
- Go installed in WSL (`sudo apt install golang`)
- Project files accessible from WSL (see notes below)

## Step-by-step Instructions

1. **Open your WSL terminal**

2. **Install build tools**
   ```bash
   sudo apt update
   sudo apt install build-essential gcc
   ```

3. **Enable CGO**
   ```bash
   export CGO_ENABLED=1
   ```

4. **Check Go installation**
   ```bash
   which go
   go version
   # Should show /usr/bin/go and a recent version
   ```

5. **Clean and prepare Go modules**
   ```bash
   go clean -modcache
   go mod tidy
   ```

6. **Build the project**
   ```bash
   go build -v ./cmd/grab
   ```

7. **If you see errors like 'undefined: Node' or 'build constraints exclude all Go files':**
   - This means go-tree-sitter grammars are not built for Linux/WSL. You may need to vendor grammars or use a compatible fork. See project README for advanced fixes.

## Notes on File System Location
- **Best practice:** Place your project in the WSL file system (e.g., `/home/<user>/codegrab`).
- **Avoid:** Building in `/mnt/c/...` or other Windows-mounted drives. This can cause file permission, symlink, and performance issues with Go and CGO.
- If you must use Windows drives, expect slower builds and possible errors with native dependencies.


## Cross-Compiling for Windows from WSL

To build a Windows executable from WSL, follow these steps:

1. **Install MinGW-w64 (Windows cross-compiler toolchain)**
   ```bash
   sudo apt update
   sudo apt install mingw-w64
   ```

2. **Set environment variables for Windows build**
   ```bash
   export CGO_ENABLED=1
   export GOOS=windows
   export GOARCH=amd64
   export CC=x86_64-w64-mingw32-gcc
   ```

3. **Build the project for Windows**
   ```bash
   go build -v -o grab.exe ./cmd/grab
   ```

4. **Troubleshooting Windows build**
   - If you see errors about missing headers, make sure `mingw-w64` is installed and `CC` is set correctly.
   - If you see CGO or go-tree-sitter errors, you may need to vendor grammars or use a compatible fork. See project README for advanced fixes.
   - If you get 'undefined: Node' or 'build constraints exclude all Go files', this is a go-tree-sitter limitation on Windows targets. You may need to remove or replace the dependency for full Windows support.

## Troubleshooting
- If you get errors about missing headers or libraries, double-check that `build-essential` and `gcc` are installed.
- If Go is not found, install it with `sudo apt install golang`.
- For advanced CGO or go-tree-sitter issues, see the official repo issues or ask for help.
