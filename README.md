<div align="center">
  <h1>CodeGrab ✋</h1>
</div>

<p align="center">
  <strong>An interactive CLI tool for selecting and bundling code into a single, LLM-ready output file.</strong>
</p>

![codegrab-demo](https://github.com/user-attachments/assets/77b8984e-913f-4646-a1f7-e4d16aa8f7b5)

## ❓ Why?

When working with LLMs, sharing code context is essential for getting accurate responses. However,
manually copying files or creating code snippets is tedious. CodeGrab streamlines this process by
providing a clean Terminal UI (TUI), alongside a versatile command-line interface (CLI). This allows
you to easily select files from your project, generate well-formatted output, and copy it directly
to your clipboard, ready for LLM processing.

## ✨ Features

- 🎮 **Interactive Mode**: Navigate your project structure with vim-like keybindings in a
  TUI environment
- 🧹 **Filtering Options**: Respect `.gitignore` rules, handle hidden files, and apply customizable glob
  patterns
- 🔍 **Fuzzy Search**: Quickly find files across your project
- ✅ **File Selection**: Toggle files or entire directories (with child items) for inclusion or exclusion
- 📄 **Multiple Output Formats**: Generate Markdown, Plain Text, or XML output
- ⏳ **Temp File**: Generate the output file in your system's temporary directory
- 📋 **Clipboard Integration**: Copy content or output file directly to your clipboard
- 🌲 **Directory Tree View**: Display a tree-style view of your project structure
- 🧮 **Token Estimation**: Get estimated token count for LLM context windows

## 📦 Installation

```sh
go install github.com/epilande/codegrab/cmd/grab@latest
```

Or build from source:

```sh
git clone https://github.com/epilande/codegrab
cd codegrab
go build ./cmd/grab
```

Then move the binary to your `PATH`

## 🚀 Quick Start

1. Go to your project directory and run:

   ```sh
   grab
   ```

2. Navigate with <kbd>h</kbd>/<kbd>j</kbd>/<kbd>k</kbd>/<kbd>l</kbd>
3. Select files using the <kbd>Space</kbd> or <kbd>Tab</kbd> key
4. Press <kbd>g</kbd> to generate output file or <kbd>y</kbd> to copy contents to clipboard

CodeGrab will generate `codegrab-output.md` in your current working directory (on macOS this file is automatically copied to your clipboard), which you can immediately send to an AI assistant for better context-aware coding assistance.

https://github.com/user-attachments/assets/48f245f4-695d-4cea-8fc0-4b0158bb46a5

## 🎮 Usage

```sh
grab [options] [directory]
```

### Arguments

| Argument    | Description                                                        |
| :---------- | :----------------------------------------------------------------- |
| `directory` | Path to the project directory (default: current working directory) |

### Options

| Option                  | Description                                                                             |
| :---------------------- | :-------------------------------------------------------------------------------------- |
| `-h, --help`            | Display help information                                                                |
| `-n, --non-interactive` | Run in non-interactive mode (grabs all files)                                           |
| `-o, --output file`     | Output file path (default: `./codegrab-output.<format>`)                                |
| `-t, --temp`            | Use system temporary directory for output file                                          |
| `-g, --glob pattern`    | Include/exclude files and directories (e.g., `--glob="*.{ts,tsx}" --glob="!*.spec.ts"`) |
| `-f, --format format`   | Output format (available: markdown, text, xml)                                          |
| `--theme`               | Set the UI theme                                                                        |

### 📖 Examples

1. Run in interactive mode (default):

   ```bash
   grab
   ```

2. Grab all files in current directory (non-interactive):

   ```bash
   grab -n
   ```

3. Grab a specific directory:

   ```bash
   grab /path/to/project
   ```

4. Specify custom output file:

   ```bash
   grab -o output.md /path/to/project
   ```

5. Generate XML output:

   ```bash
   grab -f xml -o output.xml /path/to/project
   ```

6. Filter files using glob pattern:

   ```bash
   grab -g="*.go" /path/to/project
   ```

7. Use multiple glob patterns for include/exclude:
   ```bash
   grab -g="*.{ts,tsx}" -g="!*.spec.{ts,tsx}"
   ```

## ⌨️ Keyboard Controls

### Navigation

| Action                     | Key                             | Description                                             |
| :------------------------- | :------------------------------ | :------------------------------------------------------ |
| Move cursor down           | <kbd>j</kbd> or <kbd>↓</kbd>    | Move the cursor to the next item in the list            |
| Move cursor up             | <kbd>k</kbd> or <kbd>↑</kbd>    | Move the cursor to the previous item in the list        |
| Collapse directory         | <kbd>h</kbd> or <kbd>←</kbd>    | Collapse the currently selected directory               |
| Expand directory           | <kbd>l</kbd> or <kbd>→</kbd>    | Expand the currently selected directory                 |
| Go to top                  | <kbd>H</kbd> or <kbd>home</kbd> | Jump to the first item in the list                      |
| Go to bottom               | <kbd>L</kbd> or <kbd>end</kbd>  | Jump to the last item in the list                       |
| Toggle expand/collapse all | <kbd>e</kbd>                    | Toggle between expanding and collapsing all directories |

### Search

| Action                 | Key                               | Description                                            |
| :--------------------- | :-------------------------------- | :----------------------------------------------------- |
| Start search           | <kbd>/</kbd>                      | Begin searching for files                              |
| Next search result     | <kbd>ctrl+n</kbd> or <kbd>↓</kbd> | Navigate to the next search result                     |
| Previous search result | <kbd>ctrl+p</kbd> or <kbd>↑</kbd> | Navigate to the previous search result                 |
| Select/deselect file   | <kbd>tab</kbd>                    | Toggle selection of the current file in search results |
| Exit search            | <kbd>esc</kbd>                    | Exit search mode and return to normal navigation       |

### Selection & Output

| Action               | Key                                | Description                                                  |
| :------------------- | :--------------------------------- | :----------------------------------------------------------- |
| Select/deselect item | <kbd>tab</kbd> or <kbd>space</kbd> | Toggle selection of the current file or directory            |
| Copy to clipboard    | <kbd>y</kbd>                       | Copy the generated output to clipboard                       |
| Generate output file | <kbd>g</kbd>                       | Generate the output file with selected content               |
| Cycle output formats | <kbd>F</kbd>                       | Cycle through available output formats (markdown, text, xml) |

### View Options

| Action                     | Key          | Description                                       |
| :------------------------- | :----------- | :------------------------------------------------ |
| Toggle `.gitignore` filter | <kbd>i</kbd> | Toggle whether to respect `.gitignore` rules      |
| Toggle hidden files        | <kbd>.</kbd> | Toggle visibility of hidden files and directories |
| Toggle help screen         | <kbd>?</kbd> | Show or hide the help screen                      |
| Quit                       | <kbd>q</kbd> | Exit the application                              |

## 🎨 Themes

CodeGrab comes with several built-in themes:

- Catppuccin (Latte, Frappe, Macchiato, Mocha)
- Dracula
- Nord

Select a theme using the `--theme` flag:

```sh
grab --theme=dracula
```

## 📄 Output Formats

### Markdown (Default)

```sh
grab --format markdown
```

#### Example Output

````md
# Project Structure

```
./
└── internal/
    ├── filesystem/
    │   ├── filter.go
    │   ├── gitignore.go
    │   └── walker.go
    └── generator/
        └── formats/
            ├── markdown.go
            └── registry.go
```

# Project Files

## File: `internal/filesystem/filter.go`

```go
package filesystem

import (
    "path/filepath"
    "strings"
)

// ... rest of the file content
```
````

### Plain Text

```sh
grab --format text
```

#### Example Output

```
=============================================================
PROJECT STRUCTURE
=============================================================

./
└── internal/
    ├── filesystem/
    │   ├── filter.go
    │   ├── gitignore.go
    │   └── walker.go
    └── generator/
        └── formats/
            ├── markdown.go
            └── registry.go


=============================================================
PROJECT FILES
=============================================================
=============================================================
FILE: internal/filesystem/filter.go
=============================================================

package filesystem

import (
    "path/filepath"
    "strings"
)

// ... rest of the file content
```

### XML

```sh
grab --format xml
```

#### Example Output

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project>
  <filesystem>
    <directory name=".">
      <directory name="internal">
        <directory name="filesystem">
          <file name="filter.go"/>
          <file name="gitignore.go"/>
          <file name="walker.go"/>
        </directory>
        <directory name="generator">
          <directory name="formats">
            <file name="markdown.go"/>
            <file name="registry.go"/>
          </directory>
        </directory>
      </directory>
    </directory>
  </filesystem>
  <files>
    <file path="internal/filesystem/filter.go" language="go"><![CDATA[
package filesystem

import (
    "path/filepath"
    "strings"
)

// ... rest of the file content
]]></file>
    <!-- Additional files -->
  </files>
</project>
```
