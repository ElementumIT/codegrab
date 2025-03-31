package generator

// FileData holds file content for the generated sections
type FileData struct {
	Path     string
	Content  string
	Language string
}

// TemplateData is injected into the templates
type TemplateData struct {
	Structure string
	Files     []FileData
}

// Format defines the interface for different output formats
type Format interface {
	// Render converts the template data into the specific format
	Render(data TemplateData) (string, int, error)
	// Extension returns the file extension for this format
	Extension() string
	// Name returns the name of the format
	Name() string
}
