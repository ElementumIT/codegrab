package dependencies

import (
	"context"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// GoResolver implements Resolver for Go files.
type GoResolver struct{}

// Resolve finds Go dependencies.
func (r *GoResolver) Resolve(fileContent []byte, filePath string, projectRoot string, projectModuleName string) ([]string, error) {
	dependencies := make(map[string]struct{})

	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file %s: %w", filePath, err)
	}

	// Query to find import paths within import specs
	queryStr := `(import_spec path: (interpreted_string_literal) @import_path)`
	query, err := sitter.NewQuery([]byte(queryStr), golang.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to create Go query: %w", err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(query, tree.RootNode())

	importPaths := []string{}
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			pathLiteral := node.Content(fileContent)
			importPath := strings.Trim(pathLiteral, `"`)
			if importPath != "" {
				importPaths = append(importPaths, importPath)
			}
		}
	}

	containingDir := filepath.Dir(filepath.Join(projectRoot, filePath))

	for _, importPath := range importPaths {
		if isGoStandardLibrary(importPath) {
			continue
		}

		var resolvedDir string
		var isLocal bool

		if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
			// Relative import
			absImportDir := filepath.Join(containingDir, importPath)
			resolvedDir, isLocal = normalizePath(absImportDir, "", projectRoot)
			if !isLocal {
				continue
			}
		} else if projectModuleName != "" && (importPath == projectModuleName || strings.HasPrefix(importPath, projectModuleName+"/")) {
			// Project-local module import
			relImportPath := "."
			if importPath != projectModuleName {
				relImportPath = strings.TrimPrefix(importPath, projectModuleName+"/")
			}
			absImportDir := filepath.Join(projectRoot, relImportPath)
			resolvedDir, isLocal = normalizePath(absImportDir, "", projectRoot)
			if !isLocal {
				continue
			}
		} else {
			// Likely external package (not relative, not part of current module), skip
			continue
		}

		// Find all relevant .go files in the resolved directory
		absResolvedDir := filepath.Join(projectRoot, resolvedDir)
		filesInDir, err := os.ReadDir(absResolvedDir)
		if err != nil {
			continue
		}

		for _, entry := range filesInDir {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
				goFilePath := filepath.Join(resolvedDir, entry.Name())
				goFilePath = filepath.ToSlash(goFilePath)

				if filepath.Join(projectRoot, goFilePath) != filepath.Join(projectRoot, filePath) {
					dependencies[goFilePath] = struct{}{}
				}
			}
		}
	}

	depList := make([]string, 0, len(dependencies))
	for dep := range dependencies {
		depList = append(depList, dep)
	}

	return depList, nil
}

// isGoStandardLibrary checks if an import path belongs to the Go standard library.
func isGoStandardLibrary(path string) bool {
	pkg, err := build.Import(path, "", build.FindOnly)
	firstPart := strings.Split(path, "/")[0]
	return !strings.Contains(firstPart, ".") || (err == nil && pkg.Goroot)
}
