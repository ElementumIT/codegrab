// Enhanced Windows debug utility to diagnose tree issues
// Compile for Windows: GOOS=windows GOARCH=amd64 go build -o debug_windows_enhanced.exe debug_windows_enhanced.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	fmt.Printf("=== ENHANCED WINDOWS DEBUG ===\n")
	fmt.Printf("Number of command line arguments: %d\n", len(os.Args))
	fmt.Printf("Arguments received:\n")
	for i, arg := range os.Args {
		fmt.Printf("  Args[%d]: %q\n", i, arg)
	}
	fmt.Printf("OS: %s\n", runtime.GOOS)
	fmt.Printf("ARCH: %s\n", runtime.GOARCH)
	fmt.Printf("Separator: %q (%c)\n", string(filepath.Separator), filepath.Separator)
	
	if len(os.Args) < 2 {
		fmt.Println("\n❌ ERROR: No directory path provided")
		fmt.Println("\nUsage: debug_windows_enhanced.exe <directory_path>")
		fmt.Println("\nFor paths with spaces, make sure to use quotes:")
		fmt.Println(`  debug_windows_enhanced.exe "C:\Path With Spaces\project"`)
		fmt.Println(`  debug_windows_enhanced.exe "D:\Users\bolges\Documents\Elementum\Elementum Code\zzz moved to wsl - codegrab"`)
		os.Exit(1)
	}
	
	rootPath := os.Args[1]
	
	fmt.Printf("\n=== PATH ANALYSIS ===\n")
	fmt.Printf("Root path received: %q\n", rootPath)
	fmt.Printf("Root path length: %d characters\n", len(rootPath))
	fmt.Printf("Contains spaces: %v\n", strings.Contains(rootPath, " "))
	fmt.Printf("Contains backslashes: %v\n", strings.Contains(rootPath, "\\"))
	fmt.Printf("Contains forward slashes: %v\n", strings.Contains(rootPath, "/"))
	fmt.Printf("Is absolute path: %v\n", filepath.IsAbs(rootPath))
	
	// Try to get absolute path
	absPath, absErr := filepath.Abs(rootPath)
	if absErr != nil {
		fmt.Printf("Error getting absolute path: %v\n", absErr)
	} else {
		fmt.Printf("Absolute path: %q\n", absPath)
		if absPath != rootPath {
			fmt.Printf("NOTE: Path was converted from relative to absolute\n")
			rootPath = absPath // Use absolute path for rest of analysis
		}
	}
	
	// Check if root path exists
	fmt.Printf("\n=== ROOT PATH CHECK ===\n")
	rootInfo, err := os.Stat(rootPath)
	if err != nil {
		fmt.Printf("❌ ERROR: Cannot access root path: %v\n", err)
		fmt.Printf("\nTroubleshooting tips:\n")
		fmt.Printf("1. Make sure the path exists\n")
		fmt.Printf("2. Check that you have permission to access the directory\n")
		fmt.Printf("3. If using quotes, make sure they're around the entire path\n")
		fmt.Printf("4. Try running as administrator if you get permission errors\n")
		os.Exit(1)
	}
	fmt.Printf("✅ Root path accessible\n")
	fmt.Printf("   Is directory: %v\n", rootInfo.IsDir())
	fmt.Printf("   Size: %d bytes\n", rootInfo.Size())
	fmt.Printf("   Mode: %s\n", rootInfo.Mode())
	
	if !rootInfo.IsDir() {
		fmt.Printf("❌ ERROR: Path is not a directory\n")
		os.Exit(1)
	}
	
	// Test basic directory reading
	fmt.Printf("\n=== BASIC DIRECTORY READ ===\n")
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Printf("❌ ERROR: Cannot read directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Directory readable, found %d entries\n", len(entries))
	
	// Show first few entries
	showLimit := 10
	if len(entries) < showLimit {
		showLimit = len(entries)
	}
	fmt.Printf("First %d entries:\n", showLimit)
	for i := 0; i < showLimit; i++ {
		entry := entries[i]
		entryType := "file"
		if entry.IsDir() {
			entryType = "dir"
		}
		fmt.Printf("  %s: %s\n", entryType, entry.Name())
	}
	if len(entries) > showLimit {
		fmt.Printf("  ... and %d more entries\n", len(entries)-showLimit)
	}
	
	// Walk directory manually to see what we find (like filesystem walker would do)
	fmt.Printf("\n=== DETAILED FILESYSTEM WALK ===\n")
	fileCount := 0
	dirCount := 0
	selectedFiles := make([]string, 0)
	walkErrors := make([]string, 0)
	
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errMsg := fmt.Sprintf("ERROR walking %s: %v", path, err)
			walkErrors = append(walkErrors, errMsg)
			fmt.Printf("❌ %s\n", errMsg)
			return nil
		}
		
		if path == rootPath {
			return nil // Skip root
		}
		
		// Get relative path (this is what codegrab uses internally)
		relPath, relErr := filepath.Rel(rootPath, path)
		if relErr != nil {
			errMsg := fmt.Sprintf("ERROR getting relative path for %s: %v", path, relErr)
			walkErrors = append(walkErrors, errMsg)
			fmt.Printf("❌ %s\n", errMsg)
			return nil
		}
		
		// Show every entry we find
		normalizedPath := filepath.ToSlash(relPath)
		entryType := "file"
		if info.IsDir() {
			entryType = "dir"
			dirCount++
		} else {
			fileCount++
			selectedFiles = append(selectedFiles, normalizedPath)
		}
		
		fmt.Printf("Found %s: %q\n", entryType, normalizedPath)
		fmt.Printf("  Raw path: %q\n", path)
		fmt.Printf("  Rel path: %q\n", relPath)
		fmt.Printf("  Size: %d bytes\n", info.Size())
		fmt.Println()
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("❌ ERROR during walk: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("=== WALK SUMMARY ===\n")
	fmt.Printf("✅ Walk completed successfully\n")
	fmt.Printf("Total directories found: %d\n", dirCount)
	fmt.Printf("Total files found: %d\n", fileCount)
	fmt.Printf("Walk errors: %d\n", len(walkErrors))
	
	if len(walkErrors) > 0 {
		fmt.Printf("Walk error details:\n")
		for _, errMsg := range walkErrors {
			fmt.Printf("  %s\n", errMsg)
		}
	}
	
	// Analyze what should create a tree structure
	fmt.Printf("\n=== TREE STRUCTURE ANALYSIS ===\n")
	nestedFiles := 0
	rootFiles := 0
	maxDepth := 0
	
	if fileCount == 0 {
		fmt.Printf("⚠️  WARNING: No files found! This will result in an empty tree.\n")
		fmt.Printf("Make sure your directory contains files, not just subdirectories.\n")
	} else {
		for _, filePath := range selectedFiles {
			parts := strings.Split(filePath, "/")
			depth := len(parts) - 1
			if depth > maxDepth {
				maxDepth = depth
			}
			
			if depth > 0 {
				nestedFiles++
			} else {
				rootFiles++
			}
			
			fmt.Printf("File: %q (depth: %d)\n", filePath, depth)
		}
		
		fmt.Printf("\nTree structure prediction:\n")
		fmt.Printf("  Root level files: %d\n", rootFiles)
		fmt.Printf("  Nested files: %d\n", nestedFiles)
		fmt.Printf("  Maximum depth: %d\n", maxDepth)
		
		if nestedFiles == 0 {
			fmt.Printf("⚠️  WARNING: All files are at root level - tree will appear 'flat'\n")
			fmt.Printf("This means your directory structure is:\n")
			fmt.Printf("  YourDirectory/\n")
			showCount := len(selectedFiles)
			if showCount > 5 {
				showCount = 5
			}
			for i := 0; i < showCount; i++ {
				fmt.Printf("  ├── %s\n", selectedFiles[i])
			}
			if len(selectedFiles) > 5 {
				fmt.Printf("  └── ... (%d more files)\n", len(selectedFiles)-5)
			}
		} else {
			fmt.Printf("✅ Expected tree structure with %d levels of nesting\n", maxDepth)
		}
	}
	
	// Test the exact path processing logic that buildTree uses
	fmt.Printf("\n=== PATH PROCESSING TEST (buildTree logic) ===\n")
	for i, origPath := range selectedFiles {
		if i >= 5 { // Only test first 5 files for brevity
			fmt.Printf("... (testing first 5 files only)\n")
			break
		}
		
		fmt.Printf("\n--- Testing file: %q ---\n", origPath)
		
		// Mimic the exact buildTree path processing
		normPath := strings.ReplaceAll(origPath, "\\", "/")
		normPath = filepath.ToSlash(filepath.Clean(normPath))
		fmt.Printf("  1. After normalization: %q\n", normPath)
		
		osSpecificPath := filepath.FromSlash(normPath)
		fmt.Printf("  2. OS-specific path: %q\n", osSpecificPath)
		
		fullPath := filepath.Join(rootPath, osSpecificPath)
		fmt.Printf("  3. Full path: %q\n", fullPath)
		
		info, err := os.Stat(fullPath)
		if err != nil {
			fmt.Printf("  4. ❌ Stat error: %v\n", err)
			// Try fallback (exact buildTree logic)
			fallbackPath := filepath.Join(rootPath, origPath)
			fmt.Printf("  5. Trying fallback: %q\n", fallbackPath)
			info, err = os.Stat(fallbackPath)
			if err != nil {
				fmt.Printf("  6. ❌ Fallback failed: %v\n", err)
				fmt.Printf("  🔍 This file will NOT appear in the tree!\n")
				continue
			}
			fmt.Printf("  6. ✅ Fallback succeeded\n")
		} else {
			fmt.Printf("  4. ✅ Stat succeeded\n")
		}
		
		fmt.Printf("  7. File info: size=%d, isDir=%v\n", info.Size(), info.IsDir())
		
		// Show path parts for tree building
		parts := strings.Split(normPath, "/")
		fmt.Printf("  8. Path parts: %v\n", parts)
		fmt.Printf("  9. Tree depth: %d\n", len(parts)-1)
		fmt.Printf("  🔍 This file WILL appear in the tree\n")
	}
	
	fmt.Printf("\n=== RECOMMENDATIONS ===\n")
	if fileCount == 0 {
		fmt.Printf("🔧 Your directory has no files, only subdirectories.\n")
		fmt.Printf("   Add some files to subdirectories to see a tree structure.\n")
	} else if nestedFiles == 0 {
		fmt.Printf("🔧 All your files are in the root directory.\n")
		fmt.Printf("   Move some files to subdirectories to see a tree structure.\n")
	} else {
		fmt.Printf("✅ Your directory structure should display as a tree.\n")
		fmt.Printf("   If codegrab still shows a flat tree, the issue is elsewhere.\n")
	}
	
	fmt.Printf("\n=== NEXT STEPS ===\n")
	fmt.Printf("1. Run grab_debug.exe on this same directory\n")
	fmt.Printf("2. Look for 'DEBUG buildTree:' messages in the output\n")
	fmt.Printf("3. Share both outputs for analysis\n")
	
	fmt.Printf("\n=== END ENHANCED DEBUG ===\n")
}