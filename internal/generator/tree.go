package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

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

	var normalizedOrder []string
	seen := make(map[string]bool)
	origForNormalized := make(map[string]string)

	// DEBUG: Add debug output for Windows troubleshooting
	if len(g.SelectedFiles) > 0 {
		fmt.Fprintf(os.Stderr, "DEBUG buildTree: Processing %d selected files, RootPath: %q\n", len(g.SelectedFiles), g.RootPath)
	}

       for origPath, selected := range g.SelectedFiles {
	       if !selected || origPath == "" {
		       continue
	       }

	       // DEBUG: Show what we're processing
	       fmt.Fprintf(os.Stderr, "DEBUG buildTree: Processing file %q\n", origPath)

	       // Always normalize path to use forward slashes for tree structure
	       normPath := strings.ReplaceAll(origPath, "\\", "/")
	       normPath = filepath.ToSlash(filepath.Clean(normPath))

	       // DEBUG: Show normalization result
	       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Normalized to %q\n", normPath)

	       // For filesystem access, convert normalized path to OS-specific separators
	       // This handles Windows-style backslashes correctly on both Windows and Linux
	       osSpecificPath := filepath.FromSlash(normPath)
	       origFull := filepath.Join(g.RootPath, osSpecificPath)

	       // DEBUG: Show filesystem path
	       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   OS-specific: %q, Full: %q\n", osSpecificPath, origFull)

	       info, err := os.Stat(origFull)
	       if err != nil {
		       // DEBUG: Show stat error
		       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Stat error: %v\n", err)
		       // If normalized path fails, try the original path as fallback
		       origFull = filepath.Join(g.RootPath, origPath)
		       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Trying fallback: %q\n", origFull)
		       info, err = os.Stat(origFull)
		       if err != nil {
		       	fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Fallback failed: %v\n", err)
		       	continue
		       }
		       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Fallback SUCCESS\n")
	       } else {
		       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Stat SUCCESS\n")
	       }
	       if info.IsDir() {
		       fmt.Fprintf(os.Stderr, "DEBUG buildTree:   SKIPPING: is directory\n")
		       continue // we only add files; directories inferred from files
	       }
	       fileCache := cache.GetGlobalFileCache()
	       if ok, err := fileCache.GetTextFileStatus(origFull, utils.IsTextFile); err != nil || !ok {
		       continue
	       }

	       normalized := normPath
		if filepath.IsAbs(normalized) {
			if rel, err := filepath.Rel(g.RootPath, normalized); err == nil {
				normalized = filepath.ToSlash(filepath.Clean(rel))
			}
		}
		if strings.HasPrefix(normalized, "../") || normalized == ".." {
			continue
		}
		if seen[normalized] {
			continue
		}
		seen[normalized] = true
		origForNormalized[normalized] = origPath // preserve original (with backslashes) for Node.Path
		normalizedOrder = append(normalizedOrder, normalized)

		// DEBUG: Show successful processing
		fmt.Fprintf(os.Stderr, "DEBUG buildTree:   SUCCESS: Added %q (normalized: %q)\n", origPath, normalized)
	}

	sort.Strings(normalizedOrder)
	
	// DEBUG: Show final processing list
	fmt.Fprintf(os.Stderr, "DEBUG buildTree: Final processing list (%d items):\n", len(normalizedOrder))
	for _, norm := range normalizedOrder {
		fmt.Fprintf(os.Stderr, "DEBUG buildTree:   Will process: %q\n", norm)
	}
	
	dirSet := make(map[string]bool)
	for _, normalized := range normalizedOrder {
		parts := strings.Split(normalized, "/")
		current := root
		accum := ""
		for i, part := range parts {
			if part == "" {
				continue
			}
			if accum == "" {
				accum = part
			} else {
				accum = accum + "/" + part
			}
			isLast := i == len(parts)-1
			found := false
			for _, child := range current.Children {
				if child.Name == part {
					current = child
					found = true
					break
				}
			}
			if !found {
				isDir := !isLast
				storedPath := accum
				if isLast { // file node: use original path form
					storedPath = origForNormalized[normalized]
				}
				newNode := &Node{Name: part, IsDir: isDir, Children: []*Node{}, Path: storedPath}
				current.Children = append(current.Children, newNode)
				current = newNode
				if !isDir {
					// Use normalized path for filesystem operations, with fallback
					osSpecificPath := filepath.FromSlash(origForNormalized[normalized])
					absolutePath := filepath.Join(g.RootPath, osSpecificPath)
					if err := cache.GetGlobalFileCache().CacheMetadataOnly(absolutePath); err != nil {
						// Try fallback with original path
						absolutePath = filepath.Join(g.RootPath, origForNormalized[normalized])
						if err := cache.GetGlobalFileCache().CacheMetadataOnly(absolutePath); err == nil {
							newNode.Language = determineLanguage(part)
						}
					} else {
						newNode.Language = determineLanguage(part)
					}
				}
			}
			if !isLast {
				dirSet[accum] = true
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

// fileWorkItem represents a file to be processed
type fileWorkItem struct {
	node *Node
	path string
}

// ConcurrentFileCollector handles concurrent file content reading
type ConcurrentFileCollector struct {
	maxWorkers int
	rootPath   string
}

// NewConcurrentFileCollector creates a new concurrent file collector
func NewConcurrentFileCollector(rootPath string) *ConcurrentFileCollector {
	return &ConcurrentFileCollector{maxWorkers: runtime.NumCPU(), rootPath: rootPath}
}

// collectFiles intelligently chooses between concurrent and sequential based on file count
func collectFiles(node *Node, files *[]FileData, rootPath string, secretScanner *secrets.Scanner) {
	fileCount := countFiles(node)
	const concurrentThreshold = 50
	if fileCount >= concurrentThreshold {
		collector := NewConcurrentFileCollector(rootPath)
		result, err := collector.CollectFilesConcurrent(node, secretScanner)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: concurrent file collection failed, falling back to sequential: %v\n", err)
			collectFilesSequential(node, files, rootPath, secretScanner)
			return
		}
		*files = result
	} else {
		collectFilesSequential(node, files, rootPath, secretScanner)
	}
}

func countFiles(node *Node) int {
	if !node.IsDir {
		return 1
	}
	count := 0
	for _, c := range node.Children {
		count += countFiles(c)
	}
	return count
}

// CollectFilesConcurrent performs concurrent file content reading
func (c *ConcurrentFileCollector) CollectFilesConcurrent(node *Node, secretScanner *secrets.Scanner) ([]FileData, error) {
	var workItems []fileWorkItem
	c.collectWorkItems(node, &workItems)
	if len(workItems) == 0 {
		return []FileData{}, nil
	}
	workQueue := make(chan fileWorkItem, len(workItems))
	resultQueue := make(chan FileData, len(workItems))
	errorChan := make(chan error, c.maxWorkers)
	for _, item := range workItems {
		workQueue <- item
	}
	close(workQueue)
	var wg sync.WaitGroup
	var firstError error
	for i := 0; i < c.maxWorkers; i++ {
		wg.Add(1)
		go c.fileWorker(workQueue, resultQueue, errorChan, &wg)
	}
	errorDone := make(chan struct{})
	go func() {
		defer close(errorDone)
		for err := range errorChan {
			if firstError == nil {
				firstError = err
			}
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}()
	go func() { wg.Wait(); close(resultQueue); close(errorChan) }()
	var out []FileData
	for fd := range resultQueue {
		out = append(out, fd)
	}
	<-errorDone
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out, firstError
}

func (c *ConcurrentFileCollector) collectWorkItems(node *Node, workItems *[]fileWorkItem) {
	if !node.IsDir {
		*workItems = append(*workItems, fileWorkItem{node: node, path: filepath.Join(c.rootPath, node.Path)})
	}
	for _, child := range node.Children {
		c.collectWorkItems(child, workItems)
	}
}

func (c *ConcurrentFileCollector) fileWorker(workQueue <-chan fileWorkItem, resultQueue chan<- FileData, errorChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	fileCache := cache.GetGlobalFileCache()
	for item := range workQueue {
		content, err := fileCache.GetLazy(item.path)
		if err != nil {
			select {
			case errorChan <- fmt.Errorf("failed to read file %s: %w", item.node.Path, err):
			default:
			}
			continue
		}
		item.node.Content = content
		select {
		case resultQueue <- FileData{Path: item.node.Path, Content: content, Language: item.node.Language, Findings: nil}:
		default:
		}
	}
}

func collectFilesSequential(node *Node, files *[]FileData, rootPath string, secretScanner *secrets.Scanner) {
	if !node.IsDir {
		fileCache := cache.GetGlobalFileCache()
		absolute := filepath.Join(rootPath, node.Path)
		if content, err := fileCache.GetLazy(absolute); err == nil {
			node.Content = content
			*files = append(*files, FileData{Path: node.Path, Content: content, Language: node.Language, Findings: nil})
		} else {
			fmt.Fprintf(os.Stderr, "Warning: failed to read file %s: %v\n", node.Path, err)
		}
	}
	for _, child := range node.Children {
		collectFilesSequential(child, files, rootPath, secretScanner)
	}
}

func determineLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
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
	mapping := map[string]string{"html": "html", "htm": "html", "xhtml": "html", "css": "css", "scss": "scss", "sass": "sass", "less": "less", "js": "javascript", "jsx": "jsx", "ts": "typescript", "tsx": "tsx", "vue": "vue", "svelte": "svelte", "go": "go", "py": "python", "rb": "ruby", "php": "php", "java": "java", "c": "c", "h": "c", "cpp": "cpp", "cc": "cpp", "cxx": "cpp", "hpp": "cpp", "hxx": "cpp", "cs": "csharp", "fs": "fsharp", "fsx": "fsharp", "rs": "rust", "swift": "swift", "kt": "kotlin", "kts": "kotlin", "scala": "scala", "clj": "clojure", "cljs": "clojure", "cljc": "clojure", "edn": "clojure", "hs": "haskell", "lhs": "haskell", "elm": "elm", "ex": "elixir", "exs": "elixir", "erl": "erlang", "hrl": "erlang", "dart": "dart", "pl": "perl", "pm": "perl", "r": "r", "lua": "lua", "groovy": "groovy", "tcl": "tcl", "m": "objectivec", "mm": "objectivec", "d": "d", "jl": "julia", "cr": "crystal", "nim": "nim", "zig": "zig", "v": "v", "sh": "bash", "bash": "bash", "zsh": "bash", "fish": "fish", "ps1": "powershell", "psm1": "powershell", "bat": "batch", "cmd": "batch", "awk": "awk", "json": "json", "yaml": "yaml", "yml": "yaml", "toml": "toml", "xml": "xml", "csv": "csv", "tsv": "tsv", "ini": "ini", "conf": "conf", "cfg": "conf", "plist": "xml", "md": "markdown", "markdown": "markdown", "rst": "restructuredtext", "tex": "latex", "latex": "latex", "txt": "text", "adoc": "asciidoc", "asciidoc": "asciidoc", "org": "org", "sql": "sql", "graphql": "graphql", "gql": "graphql", "proto": "protobuf", "tf": "terraform", "tfvars": "terraform", "hcl": "hcl", "dockerfile": "dockerfile", "lock": "text", "gradle": "gradle", "properties": "properties", "diff": "diff", "patch": "diff", "svg": "svg", "log": "log"}
	if lang, ok := mapping[ext]; ok {
		return lang
	}
	return ext
}
