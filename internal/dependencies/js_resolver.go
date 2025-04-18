package dependencies

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"
)

// JSResolver implements Resolver for TypeScript/JavaScript files.
type JSResolver struct{}

// Resolve finds TS/JS dependencies.
func (r *JSResolver) Resolve(fileContent []byte, filePath string, projectRoot string, projectModuleName string) ([]string, error) {
	if len(fileContent) == 0 {
		return []string{}, nil
	}

	dependencies := make(map[string]struct{})

	parser := sitter.NewParser()
	parser.SetLanguage(tsx.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JS/TS file %s: %w", filePath, err)
	}
	defer tree.Close()

	if tree.RootNode().HasError() {
		return nil, fmt.Errorf("parsing error detected in JS/TS file %s", filePath)
	}

	// Query to find import/require/export sources
	queryStr := `
        [
          (import_statement source: (string (string_fragment)? @import_path))

          (call_expression
            function: [(identifier)@id (member_expression property: (property_identifier)@id)]
            arguments: (arguments (string (string_fragment)? @import_path))
            (#eq? @id "require")
          )

          (export_statement source: (string (string_fragment)? @import_path))

          (call_expression ; Dynamic import
            function: (import)
            arguments: (arguments (string (string_fragment)? @import_path))
          )
        ] @import
    `
	query, err := sitter.NewQuery([]byte(queryStr), tsx.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to create JS/TS query: %w", err)
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	qc.Exec(query, tree.RootNode())
	defer qc.Close()

	importPaths := []string{}
	processedCaptures := make(map[uintptr]struct{})

	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			node := capture.Node
			captureName := query.CaptureNameForId(capture.Index)

			nodeID := node.ID()
			if _, processed := processedCaptures[nodeID]; processed {
				continue
			}

			if captureName == "import_path" {
				importPath := node.Content(fileContent)
				if importPath != "" {
					importPaths = append(importPaths, importPath)
					processedCaptures[nodeID] = struct{}{}
				}
			}
		}
	}

	containingDir := filepath.Dir(filepath.Join(projectRoot, filePath))

	for _, importPath := range importPaths {
		if !strings.HasPrefix(importPath, "./") && !strings.HasPrefix(importPath, "../") {
			continue
		}

		resolvedRelPath := resolveJSPath(importPath, containingDir, projectRoot)
		if resolvedRelPath != "" {
			normalizedSourcePath := filepath.ToSlash(filepath.Clean(filePath))
			if resolvedRelPath != normalizedSourcePath {
				dependencies[resolvedRelPath] = struct{}{}
			}
		}
	}

	depList := make([]string, 0, len(dependencies))
	for dep := range dependencies {
		depList = append(depList, dep)
	}

	return depList, nil
}

var allowedJSExtensions = map[string]bool{
	".ts":  true,
	".tsx": true,
	".js":  true,
	".jsx": true,
}

func resolveJSPath(importPath, containingDir, projectRoot string) string {
	basePath := filepath.Join(containingDir, importPath)

	for ext := range allowedJSExtensions {
		potentialPath := basePath + ext
		if fileExists(potentialPath) {
			if relPath, ok := normalizePath(potentialPath, "", projectRoot); ok {
				return relPath
			}
		}
	}

	if info, err := os.Stat(basePath); err == nil && info.IsDir() {
		for ext := range allowedJSExtensions {
			potentialPath := filepath.Join(basePath, "index"+ext)
			if fileExists(potentialPath) {
				if relPath, ok := normalizePath(potentialPath, "", projectRoot); ok {
					return relPath
				}
			}
		}
	}

	if fileExists(basePath) {
		if ext := filepath.Ext(basePath); allowedJSExtensions[ext] {
			if relPath, ok := normalizePath(basePath, "", projectRoot); ok {
				return relPath
			}
		}
	}

	return ""
}
