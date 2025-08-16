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


## Compile for Windows (Recommended: Using xgo)

The easiest and most reliable way to cross-compile this project for Windows (especially with CGO dependencies like `go-tree-sitter`) is to use [xgo](https://github.com/techknowlogick/xgo), a Docker-based Go cross-compiler.


### Prerequisites

- Docker Desktop must be installed and running.
- This repo must be cloned locally.

### Step-by-step Instructions (Recommended: Using Docker Compose)


1. **Install Docker Desktop**
	- Download and install Docker Desktop from [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/).

2. **Build your project for Windows using Docker Compose**
	- Open a terminal in your project root (where `docker-compose.yml` is located).
	- Run:

	  ```bash
	  docker-compose up xgo-build
	  ```

	- This will use the xgo Docker image to build the Windows executable.
	- When finished, look for a file like `grab-windows-4.0-amd64.exe` in your project directory.

	- You can edit the `docker-compose.yml` to change build targets (e.g., Linux, ARM).
	- For advanced builds, see the [xgo documentation](https://github.com/techknowlogick/xgo).

3. **(Optional) Advanced: Use xgo directly**
	 - You can still use xgo directly if you prefer, following the steps below:
		 - Pull the xgo Docker image:
			 ```bash
			 docker pull techknowlogick/xgo:latest
			 ```
		 - Install the xgo wrapper:
			 ```bash
			 go install src.techknowlogick.com/xgo@latest
			 ```
		 - Make sure `$GOPATH/bin` is in your PATH, or use the full path to `xgo`.
		 - Build your project for Windows:
			 ```bash
			 xgo --targets=windows/amd64 -out grab ./cmd/grab
			 ```
		 - This will produce a file like `grab-windows-4.0-amd64.exe` in your current directory.

1. **Install Docker Desktop**

