package dependencies

import (
	"testing"
)

func TestReadGoModFile(t *testing.T) {
	tests := []struct {
		name           string
		files          map[string]string
		setupFunc      func(string)
		cleanupFunc    func(string)
		expectedModule string
		expectError    bool
	}{
		{
			name: "Valid go.mod",
			files: map[string]string{
				"go.mod": "module example.com/testmod\n\ngo 1.21\n",
			},
			expectedModule: "example.com/testmod",
		},
		{
			name: "go.mod with comments and extra lines",
			files: map[string]string{
				"go.mod": "// comment\nrequire (\n\tgithub.com/gin-gonic/gin v1.9.1\n)\n\nmodule my/project/module\n\ngo 1.20\n",
			},
			expectedModule: "my/project/module",
		},
		{
			name: "go.mod with extra spaces around module",
			files: map[string]string{
				"go.mod": "   module    spaced/out/module   \n",
			},
			expectedModule: "spaced/out/module",
		},
		{
			name: "go.mod with multiple module lines (first should win)",
			files: map[string]string{
				"go.mod": "module first/module\ngo 1.19\nmodule second/module\n",
			},
			expectedModule: "first/module",
		},
		{
			name:           "No go.mod file",
			files:          map[string]string{},
			expectedModule: "",
		},
		{
			name: "Empty go.mod file",
			files: map[string]string{
				"go.mod": "",
			},
			expectedModule: "",
		},
		{
			name: "go.mod without module line",
			files: map[string]string{
				"go.mod": "go 1.19\n",
			},
			expectedModule: "",
		},
		{
			name: "Malformed module line",
			files: map[string]string{
				"go.mod": "module\n",
			},
			expectedModule: "",
		},
		{
			name: "Path is a directory",
			files: map[string]string{
				"go.mod/some_file": "content",
			},
			expectedModule: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, cleanup := setupTestEnv(t, tt.files)
			defer cleanup()

			if tt.setupFunc != nil {
				tt.setupFunc(tempDir)
			}
			if tt.cleanupFunc != nil {
				defer tt.cleanupFunc(tempDir)
			}

			moduleName := ReadGoModFile(tempDir)

			if tt.expectError && moduleName != "" {
				t.Errorf("ReadGoModFile() = %q, want empty string because an error was expected opening/reading go.mod", moduleName)
			}

			if !tt.expectError && moduleName != tt.expectedModule {
				t.Errorf("ReadGoModFile() = %q, want %q", moduleName, tt.expectedModule)
			}
		})
	}
}
