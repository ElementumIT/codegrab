package dependencies

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestJSResolver_Resolve(t *testing.T) {
	files := map[string]string{
		"app.ts":                   "import { utilFunc } from './utils/util';\nimport Component from '../components/Component.tsx';\nconst other = require('./other');\nexport { helper } from './helpers/index.js';\n// import commented from './commented';\nimport Default from './default_export';\n\nasync function dynamic() {\n  const dyn = await import('./dynamic');\n}\n\nexport * from './reexport';",
		"utils/util.ts":            "export function utilFunc() {}",
		"components/Component.tsx": "import React from 'react';\n\nexport default function Component() { return <div>Hi</div>; }",
		"other.js":                 "module.exports = {};",
		"helpers/index.js":         "export const helper = 1;",
		"helpers/another.js":       "export const anotherHelper = 2;",
		"commented.ts":             "export const commented = true;",
		"default_export.ts":        "export default class Def {}",
		"dynamic.ts":               "export const dyn = 1;",
		"reexport.ts":              "export const re = 1;",
		"empty.js":                 "",
		"style.css":                "body { color: blue; }",
		"invalid.ts":               "import { oops } from ",
	}

	tempDir, cleanup := setupTestEnv(t, files)
	defer cleanup()

	resolver := JSResolver{}

	tests := []struct {
		name         string
		filePath     string
		expectedDeps []string
		expectError  bool
	}{
		{
			name:     "App file with various imports",
			filePath: "app.ts",
			expectedDeps: []string{
				"default_export.ts",
				"dynamic.ts",
				"helpers/index.js",
				"other.js",
				"reexport.ts",
				"utils/util.ts",
			},
			expectError: false,
		},
		{
			name:         "Util file with no imports",
			filePath:     "utils/util.ts",
			expectedDeps: []string{},
			expectError:  false,
		},
		{
			name:         "Component file with external import",
			filePath:     "components/Component.tsx",
			expectedDeps: []string{},
			expectError:  false,
		},
		{
			name:         "Empty JS file",
			filePath:     "empty.js",
			expectedDeps: []string{},
			expectError:  false,
		},
		{
			name:         "Non-existent file path",
			filePath:     "nonexistent/file.js",
			expectedDeps: nil,
			expectError:  true,
		},
		{
			name:         "Invalid TS syntax",
			filePath:     "invalid.ts",
			expectedDeps: nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePathAbs := filepath.Join(tempDir, filepath.FromSlash(tt.filePath))
			fileContent, readErr := os.ReadFile(filePathAbs)

			if tt.filePath == "nonexistent/file.js" {
				if !errors.Is(readErr, fs.ErrNotExist) {
					t.Fatalf("Expected os.ReadFile to fail with ErrNotExist for %q, but got: %v", tt.filePath, readErr)
				}
				if tt.expectError {
					return
				}
			}

			if readErr != nil && !tt.expectError {
				t.Fatalf("Failed to read test file %s: %v", tt.filePath, readErr)
			}

			deps, err := resolver.Resolve(fileContent, tt.filePath, tempDir, "")

			if tt.expectError {
				if err == nil && !errors.Is(readErr, fs.ErrNotExist) {
					t.Errorf("Resolve(%q) error = nil, want error (readErr was: %v)", tt.filePath, readErr)
				}
				if tt.filePath == "invalid.ts" && err != nil && !strings.Contains(err.Error(), "parsing error detected") {
					t.Errorf("Resolve(%q) expected parsing error, but got: %v", tt.filePath, err)
				}
			} else {
				if err != nil {
					t.Errorf("Resolve(%q) unexpected error: %v", tt.filePath, err)
				}
				if readErr != nil {
					t.Errorf("Resolve(%q) had unexpected file read error: %v", tt.filePath, readErr)
				}

				sort.Strings(deps)
				sort.Strings(tt.expectedDeps)

				if !reflect.DeepEqual(deps, tt.expectedDeps) {
					t.Errorf("Resolve(%q) deps = %v, want %v", tt.filePath, deps, tt.expectedDeps)
				}
			}
		})
	}
}

func TestResolveJSPath(t *testing.T) {
	files := map[string]string{
		"src/app.ts":                "",
		"src/utils/index.ts":        "",
		"src/utils/helper.js":       "",
		"src/components/Button.tsx": "",
		"src/components/Card.jsx":   "",
		"src/data.json":             "",
		"src/legacy/script.js":      "",
		"src/folder/index.js":       "",
		"src/exact_dir/index.ts":    "",
		"src/exact_file.ts":         "",
		"src/needs_norm/../norm.js": "",
	}

	tempDir, cleanup := setupTestEnv(t, files)
	defer cleanup()

	containingDir := filepath.Join(tempDir, "src")

	tests := []struct {
		name         string
		importPath   string
		expectedPath string
	}{
		{name: "Direct .ts import", importPath: "./app.ts", expectedPath: "src/app.ts"},
		{name: "Import .ts without extension", importPath: "./app", expectedPath: "src/app.ts"},
		{name: "Import .js without extension", importPath: "./utils/helper", expectedPath: "src/utils/helper.js"},
		{name: "Import .tsx without extension", importPath: "./components/Button", expectedPath: "src/components/Button.tsx"},
		{name: "Import .jsx without extension", importPath: "./components/Card", expectedPath: "src/components/Card.jsx"},
		{name: "Import directory resolves index.ts first", importPath: "./utils", expectedPath: "src/utils/index.ts"},
		{name: "Import directory resolves index.js if index.ts not found", importPath: "./folder", expectedPath: "src/folder/index.js"},
		{name: "Relative path up", importPath: "../src/legacy/script", expectedPath: "src/legacy/script.js"},
		{name: "Import non-existent file", importPath: "./nonexistent", expectedPath: ""},
		{name: "Import non-JS file", importPath: "./data.json", expectedPath: ""},
		{name: "Import path outside root (simulated)", importPath: "../../external/file", expectedPath: ""},
		{name: "Import exact directory name (resolves index)", importPath: "./exact_dir", expectedPath: "src/exact_dir/index.ts"},
		{name: "Import exact file name without extension", importPath: "./exact_file", expectedPath: "src/exact_file.ts"},
		{name: "Import path needing normalization", importPath: "./needs_norm/../norm", expectedPath: "src/norm.js"},
		{name: "Import only dots and slashes", importPath: "././.", expectedPath: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved := resolveJSPath(tt.importPath, containingDir, tempDir)
			if resolved != tt.expectedPath {
				t.Errorf("resolveJSPath(%q, %q, %q) = %q, want %q", tt.importPath, containingDir, tempDir, resolved, tt.expectedPath)
			}
		})
	}
}
