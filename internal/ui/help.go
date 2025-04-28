package ui

const HelpText = `Navigation:
  j / ↓                    Move cursor down
  k / ↑                    Move cursor up
  h / ←                    Collapse directory
  l / →                    Expand directory
  H / home                 Go to top
  L / end                  Go to bottom
  e                        Toggle expand/collapse all directories

Search:
  /                        Start search
  ctrl+n / ↓               Next search result
  ctrl+p / ↑               Previous search result
  tab / enter              Select/deselect file in search results
  esc                      Exit search mode

Selection & Output:
  space / tab              Select/deselect file or directory
  y                        Copy generated output to clipboard
  g                        Generate output file
  D                        Toggle automatic dependency resolution (Go, TS/JS)
  F                        Cycle through output formats (md, txt, xml)
  S                        Toggle secret redaction (Default: On)

View Options:
  i                        Toggle .gitignore filter
  .                        Toggle hidden files
  r                        Refresh file list & reset selection
  ?                        Toggle help screen
  q / ctrl+c               Quit`

const UsageText = `Usage:
  grab [options] [directory]

  Options:
    -h, --help               Display this help information.
    -v, --version            Display version information.
    -n, --non-interactive    Run in non-interactive mode (selects all valid files).
    -o, --output <file>      Output file path (default: "./codegrab-output.<format>").
    -t, --temp               Use system temporary directory for output file.
    -g, --glob <pattern>     Include/exclude files using glob patterns. Can be used multiple times.
                             Prefix with '!' to exclude (e.g., -g="*.go" -g="\!*_test.go").
                             Supports brace expansion (e.g., -g="*.{ts,tsx}").
    -f, --format <format>    Output format. Available: markdown, text, xml (default: "markdown").
    -S, --skip-redaction     Skip automatic secret redaction via gitleaks (Default: false).
                             WARNING: This may expose sensitive information!
    --deps                   Automatically include direct dependencies for selected files (Go, JS/TS).
    --max-depth <depth>      Maximum depth for dependency resolution (-1 for unlimited, default: 1).
                             Only effective when --deps is used.
    --max-file-size <size>   Maximum file size to include (e.g., "50kb", "2MB"). No limit by default.
                             Files exceeding this size will be skipped if the limit is set.
    --theme <name>           Set the UI theme. Available: catppuccin-latte, catppuccin-frappe,
                             catppuccin-macchiato, catppuccin-mocha, rose-pine, rose-pine-dawn,
	                           rose-pine-moon, dracula, nord. (default: "catppuccin-mocha").

  Examples:
    # Run interactively in the current directory
    grab

    # Grab all files in current directory (non-interactive)
    grab -n

    # Run interactively in a specific directory, resolving dependencies
    grab --deps /path/to/project

    # Specify custom output file
    grab -o output.md /path/to/project

    # Generate XML output in a temporary file
    grab --temp -f xml

    # Filter files using glob pattern, skipping files > 50kb
    grab -g="*.go" --max-file-size 50kb 

    # Multiple glob patterns
    grab -g="*.{ts,tsx}" -g="\!*.spec.{ts,tsx}"`
