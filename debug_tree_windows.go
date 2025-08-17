// Debug utility to diagnose Windows tree issues
// Compile for Windows: GOOS=windows GOARCH=amd64 go build -o debug_tree.exe debug_tree_windows.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debug_tree <directory_path>")
		fmt.Println("Example: debug_tree.exe C:\\Users\\YourName\\project")
		os.Exit(1)
	}
	
	rootPath := os.Args[1]
	
	fmt.Printf("=== WINDOWS TREE DEBUG ===\n")
	fmt.Printf("Root path: %s\n", rootPath)
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("ARCH: %s\n", runtime.GOARCH)
	fmt.Printf("Separator: %q (%c)\n", string(filepath.Separator), filepath.Separator)
	
	// Check if root path exists
	rootInfo, err := os.Stat(rootPath)
	if err != nil {
		fmt.Printf("ERROR: Cannot access root path: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Root exists: %v, IsDir: %v\n", rootInfo != nil, rootInfo.IsDir())
	
	// Walk directory manually to see what we find
	fmt.Printf("\n=== FILESYSTEM WALK ===\n")
	fileCount := 0
	selectedFiles := make(map[string]bool)
	
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("ERROR walking %s: %v\n", path, err)
			return nil
		}
		
		if path == rootPath {
			return nil // Skip root
		}
		
		relPath, relErr := filepath.Rel(rootPath, path)
		if relErr != nil {
			fmt.Printf("ERROR getting relative path for %s: %v\n", path, relErr)
			return nil
		}
		
		normalizedPath := filepath.ToSlash(relPath)
		
		fmt.Printf("Found: %s\n", path)
		fmt.Printf("  Relative: %s\n", relPath)
		fmt.Printf("  Normalized: %s\n", normalizedPath)
		fmt.Printf("  IsDir: %v\n", info.IsDir())
		
		// Add non-directories to selected files
		if !info.IsDir() {
			selectedFiles[normalizedPath] = true
			fileCount++
		}
		
		fmt.Println()
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("ERROR during walk: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("=== SUMMARY ===\n")
	fmt.Printf("Total files found: %d\n", fileCount)
	
	fmt.Printf("\n=== SELECTED FILES ===\n")
	for path := range selectedFiles {
		fmt.Printf("Selected: %s\n", path)
	}
	
	fmt.Printf("\n=== PATH PROCESSING TEST ===\n")
	// Test the path processing logic that buildTree uses
	for origPath := range selectedFiles {
		fmt.Printf("\n--- Processing file: %s ---\n", origPath)
		
		// Show the path normalization process (from buildTree)
		normPath := strings.ReplaceAll(origPath, "\\", "/")
		normPath = filepath.ToSlash(filepath.Clean(normPath))
		fmt.Printf("  After backslash replace: %s\n", normPath)
		
		osSpecificPath := filepath.FromSlash(normPath)
		fmt.Printf("  OS-specific path: %s\n", osSpecificPath)
		
		fullPath := filepath.Join(rootPath, osSpecificPath)
		fmt.Printf("  Full path: %s\n", fullPath)
		
		info, err := os.Stat(fullPath)
		if err != nil {
			fmt.Printf("  ERROR stat: %v\n", err)
			// Try fallback
			fallbackPath := filepath.Join(rootPath, origPath)
			fmt.Printf("  Trying fallback: %s\n", fallbackPath)
			info, err = os.Stat(fallbackPath)
			if err != nil {
				fmt.Printf("  ERROR fallback stat: %v\n", err)
				continue
			}
			fmt.Printf("  Fallback SUCCESS\n")
		} else {
			fmt.Printf("  Stat SUCCESS\n")
		}
		
		fmt.Printf("  IsDir: %v\n", info.IsDir())
		
		// Show path parts for tree building
		parts := strings.Split(normPath, "/")
		fmt.Printf("  Path parts: %v\n", parts)
		fmt.Printf("  Number of parts: %d\n", len(parts))
		
		// Show what tree depth this would create
		if len(parts) > 1 {
			fmt.Printf("  Tree depth: %d (nested)\n", len(parts)-1)
		} else {
			fmt.Printf("  Tree depth: 0 (root level)\n")
		}
	}
	
	fmt.Printf("\n=== END DEBUG ===\n")
}