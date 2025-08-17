//go:build windows

package dependencies

// JSResolver implements Resolver for TypeScript/JavaScript files on Windows.
// Windows builds disable tree-sitter dependency resolution due to build issues.
type JSResolver struct{}

// Resolve returns an empty slice for Windows builds.
func (r *JSResolver) Resolve(fileContent []byte, filePath string, projectRoot string, projectModuleName string) ([]string, error) {
	// Tree-sitter dependency resolution is disabled on Windows due to build constraints
	return []string{}, nil
}