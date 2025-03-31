package utils

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CopyFileObject attempts to place a "file object" onto the clipboard.
// - On macOS, it calls AppleScript to set the clipboard to a POSIX file reference.
// - On non-macOS, it does nothing.
func CopyFileObject(filePath string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`set the clipboard to (POSIX file "%s")`, filePath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy file to clipboard: %w", err)
	}
	return nil
}

// IsTextFile returns true if the file at filePath appears to be a text file.
func IsTextFile(filePath string) (bool, error) {
	const sampleSize = 512

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsPermission(err) {
			return false, fmt.Errorf("permission denied: %w", err)
		}
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	buf := make([]byte, sampleSize)
	n, err := file.Read(buf)
	if err != nil {
		return false, fmt.Errorf("failed to read file content: %w", err)
	}

	contentType := http.DetectContentType(buf[:n])
	// Allow common text types, including JSON
	if !strings.HasPrefix(contentType, "text/") && contentType != "application/json" {
		return false, nil
	}
	return true, nil
}

func IsHiddenPath(path string) bool {
	for _, part := range strings.Split(path, string(os.PathSeparator)) {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}
