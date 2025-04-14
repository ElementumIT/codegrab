package dependencies

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ReadGoModFile reads the go.mod file at the project root to find the module name.
// Returns the module name or an empty string if not found or error occurred.
func ReadGoModFile(rootPath string) string {
	absRootPath, err := filepath.Abs(rootPath)
	goModPath := filepath.Join(absRootPath, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return "" // No go.mod or cannot read it
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				return parts[1]
			}
		}
	}
	return ""
}
