package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/epilande/codegrab/internal/filesystem"
	"github.com/epilande/codegrab/internal/generator"
	"github.com/epilande/codegrab/internal/generator/formats"
	"github.com/epilande/codegrab/internal/model"
	"github.com/epilande/codegrab/internal/ui"
	"github.com/epilande/codegrab/internal/ui/themes"
	"github.com/epilande/codegrab/internal/utils"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	themes.Initialize()

	var globPatterns stringSliceFlag
	var showHelp bool
	var showVersion bool
	var nonInteractive bool
	var outputPath string
	var useTempFile bool
	var themeName string
	var formatName string
	var skipRedaction bool

	flag.BoolVar(&showHelp, "help", false, "Display help information")
	flag.BoolVar(&showHelp, "h", false, "Display help information (shorthand)")

	flag.BoolVar(&showVersion, "version", false, "Display version information")
	flag.BoolVar(&showVersion, "v", false, "Display version information (shorthand)")

	flag.BoolVar(&nonInteractive, "non-interactive", false, "Run in non-interactive mode")
	flag.BoolVar(&nonInteractive, "n", false, "Run in non-interactive mode (shorthand)")

	flag.Var(&globPatterns, "glob", "Include/exclude files and directories (e.g., --glob=\"*.{ts,tsx}\" --glob=\"\\!*.spec.ts\")")
	flag.Var(&globPatterns, "g", "Include/exclude files and directories (shorthand)")

	flag.StringVar(&outputPath, "output", "", "Output file path (default: current directory)")
	flag.StringVar(&outputPath, "o", "", "Output file path (shorthand)")

	flag.BoolVar(&useTempFile, "temp", false, "Use system temporary directory for output file")
	flag.BoolVar(&useTempFile, "t", false, "Use system temporary directory for output file (shorthand)")

	availableThemes := strings.Join(themes.GetThemeNames(), ", ")
	themeUsage := fmt.Sprintf("UI theme (available: %s)", availableThemes)
	flag.StringVar(&themeName, "theme", "catppuccin-mocha", themeUsage)

	availableFormats := strings.Join(formats.GetFormatNames(), ", ")
	formatUsage := fmt.Sprintf("Output format (available: %s)", availableFormats)
	flag.StringVar(&formatName, "format", "markdown", formatUsage)
	flag.StringVar(&formatName, "f", "markdown", formatUsage+" (shorthand)")

	flag.BoolVar(&skipRedaction, "skip-redaction", false, "Skip automatic secret redaction (WARNING: this may expose secrets)")
	flag.BoolVar(&skipRedaction, "S", false, "Skip automatic secret redaction (shorthand)")

	flag.Parse()

	if showHelp {
		fmt.Println(ui.UsageText)
		fmt.Println()
		fmt.Println(ui.HelpText)
		os.Exit(0)
	}

	if showVersion {
		fmt.Println(utils.VersionInfo())
		os.Exit(0)
	}

	if themeName != "" {
		if err := themes.SetTheme(themeName); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v, using default theme\n", err)
		}
	}

	// Use current directory if no argument is provided
	root := "."
	if flag.NArg() > 0 {
		root = flag.Arg(0)
	}

	// Validate the provided path
	if stat, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Error: Directory %q does not exist\n", root)
		} else {
			log.Fatalf("Error accessing %q: %v\n", root, err)
		}
	} else if !stat.IsDir() {
		log.Fatalf("Error: %q is not a directory\n", root)
	}

	filterMgr := filesystem.NewFilterManager()

	for _, pattern := range globPatterns {
		filterMgr.AddGlobPattern(pattern)
	}

	if nonInteractive {
		runNonInteractive(root, filterMgr, outputPath, useTempFile, formatName, skipRedaction)
	} else {
		config := model.Config{
			RootPath:      root,
			FilterMgr:     filterMgr,
			OutputPath:    outputPath,
			UseTempFile:   useTempFile,
			Format:        formatName,
			SkipRedaction: skipRedaction,
		}

		m := model.NewModel(config)
		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			log.Fatalf("Error running program: %v\n", err)
		}
	}
}

// runNonInteractive processes files and generates output without user interaction
func runNonInteractive(rootPath string, filterMgr *filesystem.FilterManager, outputPath string, useTempFile bool, formatName string, skipRedaction bool) {
	gitIgnoreMgr, err := filesystem.NewGitIgnoreManager(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .gitignore: %v\n", err)
		os.Exit(1)
	}

	files, err := filesystem.WalkDirectory(rootPath, gitIgnoreMgr, filterMgr, true, false)
	if err != nil {
		log.Fatalf("Error walking directory: %v\n", err)
	}

	gen := generator.NewGenerator(rootPath, gitIgnoreMgr, filterMgr, outputPath, useTempFile)
	format := formats.GetFormat(formatName)
	gen.SetFormat(format)
	gen.SetRedactionMode(!skipRedaction)

	// Automatically select all non-directory files
	selectedFiles := make(map[string]bool)
	for _, file := range files {
		if !file.IsDir {
			selectedFiles[file.Path] = true
		}
	}

	gen.SelectedFiles = selectedFiles

	outputFilePath, tokenCount, secretCount, err := gen.Generate()
	if err != nil {
		log.Fatalf("Error generating output: %v\n", err)
	}

	fmt.Printf("✅ Generated %s (%d tokens)\n", outputFilePath, tokenCount)

	if secretCount > 0 && skipRedaction {
		fmt.Fprintf(os.Stderr, "⚠️ WARNING: %d secrets detected in the output and redaction was skipped!\n", secretCount)
	} else if secretCount > 0 && !skipRedaction {
		fmt.Fprintf(os.Stderr, "ℹ️ INFO: %d secrets detected and redacted in the output.\n", secretCount)
	}
}
