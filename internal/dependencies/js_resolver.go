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
	dependencies := make(map[string]struct{})

	parser := sitter.NewParser()
	parser.SetLanguage(tsx.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, fileContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JS/TS file %s: %w", filePath, err)
	}

	// Query to find import/require/export sources
	queryStr := `
        [
          (import_statement source: (string (string_fragment) @import_path))

          (call_expression
            function: [
              (identifier) @id         ; Capture direct identifier
              (member_expression       ; Or capture property identifier in member expr
                 property: (property_identifier) @id)
            ]
            arguments: (arguments (string (string_fragment) @import_path))
            (#eq? @id "require")       ; Check if the captured @id is "require"
          )

          (export_statement source: (string (string_fragment) @import_path))
        ] @import
    `
	query, err := sitter.NewQuery([]byte(queryStr), tsx.GetLanguage())
	if err != nil {
		return nil, fmt.Errorf("failed to create JS/TS query: %w", err)
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
			nodeName := query.CaptureNameForId(capture.Index)

			if nodeName == "import_path" {
				importPath := node.Content(fileContent)
				if importPath != "" {
					importPaths = append(importPaths, importPath)
				}
			}
		}
	}

	containingDir := filepath.Dir(filepath.Join(projectRoot, filePath))

	for _, importPath := range importPaths {
		// Skip non-relative paths (likely node_modules or aliases)
		// TODO: Add support for tsconfig path aliases
		if !strings.HasPrefix(importPath, "./") && !strings.HasPrefix(importPath, "../") {
			continue
		}

		resolvedRelPath := resolveJSPath(importPath, containingDir, projectRoot)
		if resolvedRelPath != "" {
			if filepath.Join(projectRoot, resolvedRelPath) != filepath.Join(projectRoot, filePath) {
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

// resolveJSPath attempts to resolve a relative TS/JS import path to an existing file.
func resolveJSPath(importPath, containingDir, projectRoot string) string {
	extensions := []string{".ts", ".tsx", ".js", ".jsx", ""}
	basePath := filepath.Join(containingDir, importPath)

	for _, ext := range extensions {
		potentialPath := basePath + ext
		if fileExists(potentialPath) {
			if relPath, ok := normalizePath(potentialPath, "", projectRoot); ok {
				return relPath
			}
		}
	}

	if info, err := os.Stat(basePath); err == nil && info.IsDir() {
		for _, ext := range extensions {
			potentialPath := filepath.Join(basePath, "index"+ext)
			if fileExists(potentialPath) {
				if relPath, ok := normalizePath(potentialPath, "", projectRoot); ok {
					return relPath
				}
			}
		}
	}

	if fileExists(basePath) {
		if relPath, ok := normalizePath(basePath, "", projectRoot); ok {
			return relPath
		}
	}

	return ""
}
