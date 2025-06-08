package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/epilande/codegrab/internal/cache"
	"github.com/epilande/codegrab/internal/secrets"
	"github.com/epilande/codegrab/internal/utils"
)

// Node represents a node in the in-memory file tree.
type Node struct {
	Name     string
	Path     string
	Content  string
	Language string
	Children []*Node
	IsDir    bool
	Findings []secrets.Finding
}

func (g *Generator) buildTree() *Node {
	root := &Node{
		Name:     filepath.Base(g.RootPath),
		IsDir:    true,
		Children: []*Node{},
	}
	var paths []string
	for path, selected := range g.SelectedFiles {
		if selected {
			fullPath := filepath.Join(g.RootPath, path)
			info, err := os.Stat(fullPath)
			if err != nil {
				continue
			}
			if !info.IsDir() {
				fileCache := cache.GetGlobalFileCache()
				if ok, err := fileCache.GetTextFileStatus(fullPath, utils.IsTextFile); err != nil || !ok {
					continue
				}
			}
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	dirSet := make(map[string]bool)
	for _, path := range paths {
		parts := strings.Split(path, string(os.PathSeparator))
		current := root
		fullPath := ""
		for i, part := range parts {
			fullPath = filepath.Join(fullPath, part)
			isLast := i == len(parts)-1
			found := false
			isDir := !isLast || dirSet[fullPath]
			if isLast {
				if !dirSet[fullPath] {
					if info, err := os.Stat(filepath.Join(g.RootPath, fullPath)); err == nil && info.IsDir() {
						isDir = true
						dirSet[fullPath] = true
					}
				}
			} else {
				isDir = true
				dirSet[fullPath] = true
			}
			for _, child := range current.Children {
				if child.Name == part {
					current = child
					found = true
					break
				}
			}
			if !found {
				newNode := &Node{
					Name:     part,
					IsDir:    isDir,
					Children: []*Node{},
					Path:     fullPath,
				}
				current.Children = append(current.Children, newNode)
				current = newNode
				if !isDir {
					fileCache := cache.GetGlobalFileCache()
					absolutePath := filepath.Join(g.RootPath, fullPath)
					if content, err := fileCache.Get(absolutePath); err == nil {
						newNode.Content = content
						newNode.Language = determineLanguage(part)
						if g.SecretScanner != nil {
							findings, scanErr := g.SecretScanner.Scan(content)
							if scanErr != nil {
								fmt.Fprintf(os.Stderr, "Warning: failed to scan %s for secrets: %v\n", fullPath, scanErr)
							} else if len(findings) > 0 {
								newNode.Findings = findings
							}
						}
					} else {
						fmt.Fprintf(os.Stderr, "Warning: failed to read file %s: %v\n", fullPath, err)
					}
				}
			}
		}
	}
	pruneEmptyDirectories(root)
	sortTree(root)
	return root
}

func pruneEmptyDirectories(node *Node) bool {
	var newChildren []*Node
	for _, child := range node.Children {
		if pruneEmptyDirectories(child) {
			newChildren = append(newChildren, child)
		}
	}
	node.Children = newChildren
	return !node.IsDir || len(node.Children) > 0
}

func sortTree(node *Node) {
	sort.Slice(node.Children, func(i, j int) bool {
		if node.Children[i].IsDir != node.Children[j].IsDir {
			return node.Children[i].IsDir
		}
		return node.Children[i].Name < node.Children[j].Name
	})
	for _, child := range node.Children {
		sortTree(child)
	}
}

func renderTree(node *Node, prefix string, isLast bool, builder *strings.Builder, rootName string, deselectedFiles map[string]bool) {
	if deselectedFiles[node.Path] {
		return
	}
	if node.Name != rootName {
		branch := "├── "
		if isLast {
			branch = "└── "
		}
		name := node.Name
		if node.IsDir {
			name += "/"
		}
		fmt.Fprintf(builder, "%s%s%s\n", prefix, branch, name)
	}
	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}
	for i, child := range node.Children {
		renderTree(child, newPrefix, i == len(node.Children)-1, builder, rootName, deselectedFiles)
	}
}

func collectFiles(node *Node, files *[]FileData) {
	if !node.IsDir && node.Content != "" {
		*files = append(*files, FileData{
			Path:     node.Path,
			Content:  node.Content,
			Language: node.Language,
			Findings: node.Findings,
		})
	}
	for _, child := range node.Children {
		collectFiles(child, files)
	}
}

func determineLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		// Check for special filenames without extensions
		name := strings.ToLower(filepath.Base(filename))
		switch name {
		case "makefile", "dockerfile", "jenkinsfile", "vagrantfile":
			return name
		case "gemfile", "rakefile":
			return "ruby"
		case ".gitignore", ".dockerignore", ".gitattributes":
			return "gitignore"
		case ".env", ".env.local", ".env.development", ".env.production":
			return "dotenv"
		case "readme", "license", "contributing", "changelog":
			return "markdown"
		default:
			return "text"
		}
	}

	ext = ext[1:]

	extensionMap := map[string]string{
		"html":       "html",
		"htm":        "html",
		"xhtml":      "html",
		"css":        "css",
		"scss":       "scss",
		"sass":       "sass",
		"less":       "less",
		"js":         "javascript",
		"jsx":        "jsx",
		"ts":         "typescript",
		"tsx":        "tsx",
		"vue":        "vue",
		"svelte":     "svelte",
		"go":         "go",
		"py":         "python",
		"rb":         "ruby",
		"php":        "php",
		"java":       "java",
		"c":          "c",
		"h":          "c",
		"cpp":        "cpp",
		"cc":         "cpp",
		"cxx":        "cpp",
		"hpp":        "cpp",
		"hxx":        "cpp",
		"cs":         "csharp",
		"fs":         "fsharp",
		"fsx":        "fsharp",
		"rs":         "rust",
		"swift":      "swift",
		"kt":         "kotlin",
		"kts":        "kotlin",
		"scala":      "scala",
		"clj":        "clojure",
		"cljs":       "clojure",
		"cljc":       "clojure",
		"edn":        "clojure",
		"hs":         "haskell",
		"lhs":        "haskell",
		"elm":        "elm",
		"ex":         "elixir",
		"exs":        "elixir",
		"erl":        "erlang",
		"hrl":        "erlang",
		"dart":       "dart",
		"pl":         "perl",
		"pm":         "perl",
		"r":          "r",
		"lua":        "lua",
		"groovy":     "groovy",
		"tcl":        "tcl",
		"m":          "objectivec",
		"mm":         "objectivec",
		"d":          "d",
		"jl":         "julia",
		"cr":         "crystal",
		"nim":        "nim",
		"zig":        "zig",
		"v":          "v",
		"sh":         "bash",
		"bash":       "bash",
		"zsh":        "bash",
		"fish":       "fish",
		"ps1":        "powershell",
		"psm1":       "powershell",
		"bat":        "batch",
		"cmd":        "batch",
		"awk":        "awk",
		"json":       "json",
		"yaml":       "yaml",
		"yml":        "yaml",
		"toml":       "toml",
		"xml":        "xml",
		"csv":        "csv",
		"tsv":        "tsv",
		"ini":        "ini",
		"conf":       "conf",
		"cfg":        "conf",
		"plist":      "xml",
		"md":         "markdown",
		"markdown":   "markdown",
		"rst":        "restructuredtext",
		"tex":        "latex",
		"latex":      "latex",
		"txt":        "text",
		"adoc":       "asciidoc",
		"asciidoc":   "asciidoc",
		"org":        "org",
		"sql":        "sql",
		"graphql":    "graphql",
		"gql":        "graphql",
		"proto":      "protobuf",
		"tf":         "terraform",
		"tfvars":     "terraform",
		"hcl":        "hcl",
		"dockerfile": "dockerfile",
		"lock":       "text",
		"gradle":     "gradle",
		"properties": "properties",
		"diff":       "diff",
		"patch":      "diff",
		"svg":        "svg",
		"log":        "log",
	}

	if language, ok := extensionMap[ext]; ok {
		return language
	}

	// If not found, return the extension itself
	return ext
}
