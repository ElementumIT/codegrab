package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/epilande/codegrab/internal/generator"
)

// Debug utility to help diagnose Windows tree issues
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: debug_tree <directory_path>")
		os.Exit(1)
	}
	
	rootPath := os.Args[1]
	
	fmt.Printf("=== DEBUG TREE ANALYSIS ===\n")
	fmt.Printf("Root path: %s\n", rootPath)
	fmt.Printf("OS: %s\n", filepath.Separator)
	fmt.Printf("Separator: %c\n", filepath.Separator)
	
	// Create generator
	gen := generator.NewGenerator(rootPath, nil, nil, "", false)
	
	// Walk directory manually to see what we find
	fmt.Printf("\n=== FILESYSTEM WALK ===\n")
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
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
		
		fmt.Printf("Found: %s -> %s -> %s (IsDir: %v)\n", 
			path, relPath, normalizedPath, info.IsDir())
		
		// Add non-directories to selected files
		if !info.IsDir() {
			gen.SelectedFiles[normalizedPath] = true
		}
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("ERROR during walk: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("\n=== SELECTED FILES ===\n")
	for path, selected := range gen.SelectedFiles {
		if selected {
			fmt.Printf("Selected: %s\n", path)
		}
	}
	
	fmt.Printf("\n=== BUILDING TREE ===\n")
	// Use reflection or copy the buildTree logic to add debug output
	debugBuildTree(gen)
}

// Debug version of buildTree with verbose output
func debugBuildTree(g *generator.Generator) {
	fmt.Printf("Building tree with %d selected files\n", len(g.SelectedFiles))
	
	for origPath, selected := range g.SelectedFiles {
		if !selected || origPath == "" {
			continue
		}
		
		fmt.Printf("\n--- Processing file: %s ---\n", origPath)
		
		// Show the path normalization process
		normPath := strings.ReplaceAll(origPath, "\\", "/")
		normPath = filepath.ToSlash(filepath.Clean(normPath))
		fmt.Printf("  Normalized: %s\n", normPath)
		
		osSpecificPath := filepath.FromSlash(normPath)
		fmt.Printf("  OS-specific: %s\n", osSpecificPath)
		
		origFull := filepath.Join(g.RootPath, osSpecificPath)
		fmt.Printf("  Full path: %s\n", origFull)
		
		info, err := os.Stat(origFull)
		if err != nil {
			fmt.Printf("  ERROR stat (trying fallback): %v\n", err)
			// Try fallback
			origFull = filepath.Join(g.RootPath, origPath)
			fmt.Printf("  Fallback path: %s\n", origFull)
			info, err = os.Stat(origFull)
			if err != nil {
				fmt.Printf("  ERROR stat fallback: %v\n", err)
				continue
			}
		}
		
		fmt.Printf("  File exists: %v, IsDir: %v\n", info != nil, info.IsDir())
		
		if info.IsDir() {
			fmt.Printf("  SKIPPED: is directory\n")
			continue
		}
		
		// Show path parts for tree building
		parts := strings.Split(normPath, "/")
		fmt.Printf("  Path parts: %v\n", parts)
		
		fmt.Printf("  SUCCESS: File will be included in tree\n")
	}
}