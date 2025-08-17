//go:build windows

package dependencies

// GoResolver implements Resolver for Go files on Windows.
// Windows builds disable tree-sitter dependency resolution due to build issues.
type GoResolver struct{}

// Resolve returns an empty slice for Windows builds.
func (r *GoResolver) Resolve(fileContent []byte, filePath string, projectRoot string, projectModuleName string) ([]string, error) {
	// Tree-sitter dependency resolution is disabled on Windows due to build constraints
	return []string{}, nil
}