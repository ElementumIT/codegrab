package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epilande/codegrab/internal/filesystem"
	"github.com/epilande/codegrab/internal/utils"
)

// Generator organizes how we generate the output in different formats
type Generator struct {
	format          Format
	SelectedFiles   map[string]bool
	DeselectedFiles map[string]bool
	GitIgnoreMgr    *filesystem.GitIgnoreManager
	FilterMgr       *filesystem.FilterManager
	RootPath        string
	OutputPath      string
	UseTempFile     bool
	UseGitIgnore    bool
	ShowHidden      bool
}

// NewGenerator constructs a generator with default settings
func NewGenerator(rootPath string, gitIgnoreMgr *filesystem.GitIgnoreManager, filterMgr *filesystem.FilterManager, outputPath string, useTempFile bool) *Generator {
	return &Generator{
		RootPath:        rootPath,
		OutputPath:      outputPath,
		UseTempFile:     useTempFile,
		SelectedFiles:   make(map[string]bool),
		DeselectedFiles: make(map[string]bool),
		GitIgnoreMgr:    gitIgnoreMgr,
		FilterMgr:       filterMgr,
		UseGitIgnore:    true,
		ShowHidden:      false,
		// The format will be set by SetFormat
	}
}

// SetFormat changes the output format
func (g *Generator) SetFormat(format Format) {
	g.format = format
}

// GetFormat returns the current format
func (g *Generator) GetFormat() Format {
	return g.format
}

// GetFormatName returns the name of the current format
func (g *Generator) GetFormatName() string {
	if g.format == nil {
		return "unknown"
	}
	return g.format.Name()
}

// Generate creates an output file in the specified format
func (g *Generator) Generate() (string, int, error) {
	if len(g.SelectedFiles) == 0 {
		return "", 0, fmt.Errorf("no files selected, skipping generation")
	}

	if g.format == nil {
		return "", 0, fmt.Errorf("no format set, cannot generate output")
	}

	data, err := g.PrepareTemplateData()
	if err != nil {
		return "", 0, fmt.Errorf("failed to prepare template data: %w", err)
	}

	content, tokenCount, err := g.format.Render(data)
	if err != nil {
		return "", 0, fmt.Errorf("failed to render %s: %w", g.format.Name(), err)
	}

	var outputPath string
	var displayPath string

	if g.UseTempFile {
		tmpFile, err := os.CreateTemp("", fmt.Sprintf("codegrab-*%s", g.format.Extension()))
		if err != nil {
			return "", 0, fmt.Errorf("failed to create temporary file: %w", err)
		}
		defer tmpFile.Close()

		if _, err := tmpFile.Write([]byte(content)); err != nil {
			return tmpFile.Name(), tokenCount, fmt.Errorf("failed to write to temporary file: %w", err)
		}
		outputPath = tmpFile.Name()
		displayPath = outputPath
	} else {
		if g.OutputPath != "" {
			// If output path doesn't have the correct extension, add it
			if !strings.HasSuffix(g.OutputPath, g.format.Extension()) {
				g.OutputPath += g.format.Extension()
			}
			outputPath = g.OutputPath
		} else {
			outputPath = fmt.Sprintf("./codegrab-output%s", g.format.Extension())
		}

		displayPath = outputPath
		absPath, err := filepath.Abs(outputPath)
		if err != nil {
			return "", tokenCount, fmt.Errorf("failed to get absolute path: %w", err)
		}

		if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
			return "", tokenCount, fmt.Errorf("failed to write to output file: %w", err)
		}

		outputPath = absPath
	}

	if err := utils.CopyFileObject(outputPath); err != nil {
		return displayPath, tokenCount, fmt.Errorf("clipboard copy failed: %w", err)
	}

	return displayPath, tokenCount, nil
}

// GenerateString returns the rendered content as a string along with its estimated token count
func (g *Generator) GenerateString() (string, int, error) {
	if len(g.SelectedFiles) == 0 {
		return "", 0, fmt.Errorf("no files selected, skipping generation")
	}

	if g.format == nil {
		return "", 0, fmt.Errorf("no format set, cannot generate output")
	}

	data, err := g.PrepareTemplateData()
	if err != nil {
		return "", 0, fmt.Errorf("failed to prepare template data: %w", err)
	}

	return g.format.Render(data)
}

// PrepareTemplateData finalizes the selection and builds TemplateData for the template
func (g *Generator) PrepareTemplateData() (TemplateData, error) {
	expandedSelection := make(map[string]bool)
	for path := range g.SelectedFiles {
		fullPath := filepath.Join(g.RootPath, path)
		info, err := os.Stat(fullPath)
		if err != nil {
			// Skip files that no longer exist
			continue
		}
		if !g.ShowHidden && strings.HasPrefix(info.Name(), ".") {
			continue
		}
		if !g.FilterMgr.ShouldInclude(path) {
			continue
		}
		if info.IsDir() {
			if err := filepath.Walk(fullPath, func(walkPath string, info os.FileInfo, err error) error {
				if err != nil {
					return fmt.Errorf("failed to walk path %s: %w", walkPath, err)
				}
				if !g.ShowHidden && strings.HasPrefix(info.Name(), ".") {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				relPath, err := filepath.Rel(g.RootPath, walkPath)
				if err != nil {
					return fmt.Errorf("failed to get relative path for %s: %w", walkPath, err)
				}
				if g.UseGitIgnore && g.GitIgnoreMgr.IsIgnored(walkPath) {
					return nil
				}
				if !g.FilterMgr.ShouldInclude(relPath) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if !g.DeselectedFiles[relPath] {
					expandedSelection[relPath] = true
				}
				return nil
			}); err != nil {
				return TemplateData{}, fmt.Errorf("failed to walk directory %s: %w", fullPath, err)
			}
		} else {
			expandedSelection[path] = true
		}
	}
	g.SelectedFiles = expandedSelection

	rootNode := g.buildTree()
	var structureBuilder strings.Builder
	fmt.Fprintf(&structureBuilder, "%s/\n", filepath.Base(g.RootPath))
	for i, child := range rootNode.Children {
		isLast := i == len(rootNode.Children)-1
		renderTree(child, "", isLast, &structureBuilder, filepath.Base(g.RootPath), g.DeselectedFiles)
	}

	var filesData []FileData
	collectFiles(rootNode, &filesData)

	return TemplateData{
		Structure: structureBuilder.String(),
		Files:     filesData,
	}, nil
}
