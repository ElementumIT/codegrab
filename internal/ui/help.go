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
  tab                      Select/deselect file in search
  esc                      Exit search

Selection & Output:
  space / tab              Select/deselect file or directory
  y                        Copy output to clipboard
  g                        Generate output file
  D                        Toggle automatic dependency resolution (Go, TS/JS)
  F                        Cycle through output formats (md, txt, xml)
  S                        Toggle secret redaction (Default: On)

View Options:
  i                        Toggle .gitignore filter
  .                        Toggle hidden files
  ?                        Toggle help screen
  q                        Quit`

const UsageText = `Usage:
  grab [options] [directory]

  Options:
    -h, --help               Display this help information
    -n, --non-interactive    Run in non-interactive mode
    -o, --output file        Output file path (default: current directory)
    -t, --temp               Use system temporary directory for output file
    -g, --glob pattern       Include/exclude files and directories
                             (e.g., --glob="*.{ts,tsx}" --glob="\\!*.spec.ts")
    -f, --format format      Output format (available: markdown, text, xml)
    -S, --skip-redaction     Skip automatic secret redaction (Default: false)
                             WARNING: This may expose sensitive information!
    --theme                  Set the UI theme (available: catppuccin-latte,
                             catppuccin-frappe, catppuccin-macchiato,
                             catppuccin-mocha, dracula, nord)

  Examples:
    # Grab files in current directory interactively
    grab

	  # Grab all files in current directory (non-interactive)
		grab -n

    # Grab specific directory interactively including dependencies
		grab --deps /path/to/project

    # Specify custom output file
    grab -o output.md /path/to/project

    # Generate XML output
    grab -f xml -o output.xml /path/to/project

    # Filter files using glob pattern
    grab -g="*.go" /path/to/project

    # Multiple glob patterns
    grab -g="*.{ts,tsx}" -g="\\!*.spec.{ts,tsx}"`
