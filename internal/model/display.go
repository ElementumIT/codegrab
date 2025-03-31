package model

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/epilande/codegrab/internal/filesystem"
)

func (m *Model) buildDisplayNodes() {
	m.displayNodes = nil
	var rootFiles, rootDirs []filesystem.FileItem
	// Separate files and directories for proper ordering
	for _, f := range m.files {
		if !strings.Contains(f.Path, "/") {
			if f.IsDir {
				rootDirs = append(rootDirs, f)
			} else {
				rootFiles = append(rootFiles, f)
			}
		}
	}

	// Sort directories first, then files
	sort.Slice(rootDirs, func(i, j int) bool {
		return rootDirs[i].Path < rootDirs[j].Path
	})
	sort.Slice(rootFiles, func(i, j int) bool {
		return rootFiles[i].Path < rootFiles[j].Path
	})

	// Add sorted directories and their children
	for i, d := range rootDirs {
		isLast := i == len(rootDirs)-1 && len(rootFiles) == 0
		m.addNodeAndChildren(d, 0, isLast)
	}

	for i, f := range rootFiles {
		m.displayNodes = append(m.displayNodes, FileNode{
			Path:         f.Path,
			Name:         f.Path,
			IsDir:        false,
			Level:        0,
			IsLast:       i == len(rootFiles)-1,
			Selected:     m.selected[f.Path],
			IsDeselected: m.deselected[f.Path],
		})
	}
}

func (m *Model) addNodeAndChildren(item filesystem.FileItem, level int, isLast bool) {
	node := FileNode{
		Path:         item.Path,
		Name:         filepath.Base(item.Path),
		IsDir:        item.IsDir,
		Level:        level,
		IsLast:       isLast,
		Selected:     m.selected[item.Path],
		IsDeselected: m.deselected[item.Path],
	}
	m.displayNodes = append(m.displayNodes, node)

	if !item.IsDir || m.collapsed[item.Path] {
		return
	}

	prefix := item.Path + string(os.PathSeparator)
	var children []filesystem.FileItem
	for _, f := range m.files {
		if strings.HasPrefix(f.Path, prefix) {
			sub := f.Path[len(prefix):]
			if !strings.Contains(sub, string(os.PathSeparator)) {
				children = append(children, f)
			}
		}
	}

	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return children[i].Path < children[j].Path
	})

	for i, c := range children {
		childLast := i == len(children)-1
		m.addNodeAndChildren(c, level+1, childLast)
	}
}
