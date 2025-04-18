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
	if err != nil {
		return ""
	}
	goModPath := filepath.Join(absRootPath, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) == 2 && parts[0] == "module" {
				return parts[1]
			}
			return ""
		}
	}
	return ""
}
