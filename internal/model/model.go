package model

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"

	"github.com/epilande/codegrab/internal/dependencies"
	"github.com/epilande/codegrab/internal/filesystem"
	"github.com/epilande/codegrab/internal/generator"
	"github.com/epilande/codegrab/internal/generator/formats"
	"github.com/epilande/codegrab/internal/ui"
)

// FileNode represents a file/folder in the display list.
type FileNode struct {
	Path         string
	Name         string
	Level        int
	IsDir        bool
	IsLast       bool
	Selected     bool
	IsDeselected bool
	IsDependency bool
}

type Model struct {
	err               error
	selected          map[string]bool
	deselected        map[string]bool
	collapsed         map[string]bool
	isDependency      map[string]bool
	gitIgnoreMgr      *filesystem.GitIgnoreManager
	filterMgr         *filesystem.FilterManager
	generator         *generator.Generator
	rootPath          string
	projectModuleName string
	successMsg        string
	warningMsg        string
	viewport          viewport.Model
	files             []filesystem.FileItem
	displayNodes      []FileNode
	searchResults     []FileNode
	searchInput       textinput.Model
	cursor            int
	width             int
	height            int
	showHelp          bool
	useGitIgnore      bool
	showHidden        bool
	isSearching       bool
	isGrabbing        bool
	redactSecrets     bool
	resolveDeps       bool
}

type Config struct {
	FilterMgr     *filesystem.FilterManager
	RootPath      string
	OutputPath    string
	Format        string
	UseTempFile   bool
	SkipRedaction bool
	ResolveDeps   bool
}

func NewModel(config Config) Model {
	gitIgnoreMgr, err := filesystem.NewGitIgnoreManager(config.RootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading .gitignore: %v\n", err)
		os.Exit(1)
	}

	gen := generator.NewGenerator(config.RootPath, gitIgnoreMgr, config.FilterMgr, config.OutputPath, config.UseTempFile)
	format := formats.GetFormat(config.Format)
	gen.SetFormat(format)
	gen.SetRedactionMode(!config.SkipRedaction)

	moduleName := dependencies.ReadGoModFile(config.RootPath)

	return Model{
		rootPath:          config.RootPath,
		selected:          make(map[string]bool),
		deselected:        make(map[string]bool),
		collapsed:         make(map[string]bool),
		isDependency:      make(map[string]bool),
		useGitIgnore:      true,
		gitIgnoreMgr:      gitIgnoreMgr,
		filterMgr:         config.FilterMgr,
		generator:         gen,
		redactSecrets:     !config.SkipRedaction,
		resolveDeps:       config.ResolveDeps,
		projectModuleName: moduleName,
		showHidden:        false,
		searchInput:       ui.NewSearchInput(),
		viewport: viewport.Model{
			Width:  80,
			Height: 10,
		},
	}
}
