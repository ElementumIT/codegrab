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
	if len(fileContent) == 0 {
		return nil, nil
	}

	dependencies := make(map[string]struct{})

	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Go file %s: %w", filePath, err)
	}
	defer tree.Close()

	if tree.RootNode().HasError() {
		return nil, fmt.Errorf("parsing error detected in Go file %s", filePath)
	}

	// Query to find import paths within import specs
	queryStr := `(import_spec path: (interpreted_string_literal) @import_path)`
	query, err := sitter.NewQuery([]byte(queryStr), golang.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to create Go query: %w", err)
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	qc.Exec(query, tree.RootNode())
	defer qc.Close()

	importPaths := []string{}
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, capture := range match.Captures {
			node := capture.Node
			pathLiteral := node.Content(fileContent)
			importPath := strings.TrimSpace(strings.Trim(pathLiteral, `"`))
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
			// External package
			continue
		}

		// Find all relevant .go files in the resolved directory
		absResolvedDir := filepath.Join(projectRoot, resolvedDir)
		filesInDir, err := os.ReadDir(absResolvedDir)
		if err != nil {
			continue
		}

		for _, entry := range filesInDir {
			entryName := entry.Name()
			if !entry.IsDir() && strings.HasSuffix(entryName, ".go") && !strings.HasSuffix(entryName, "_test.go") {
				goFilePath := filepath.Join(resolvedDir, entryName)
				goFilePath = filepath.ToSlash(goFilePath)
				dependencies[goFilePath] = struct{}{}
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
	if path == "" || path == "C" || strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/") {
		return false
	}

	firstPart := path
	if idx := strings.Index(path, "/"); idx != -1 {
		firstPart = path[:idx]
	}
	if strings.Contains(firstPart, ".") {
		return false
	}

	pkg, err := build.Import(path, "", build.FindOnly)
	return err == nil && pkg.Goroot
}
