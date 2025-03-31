package formats

import (
	"github.com/epilande/codegrab/internal/generator"
)

// formatRegistry maps format names to their constructors
var formatRegistry = map[string]func() generator.Format{
	"markdown": func() generator.Format { return &MarkdownFormat{} },
	"text":     func() generator.Format { return &TxtFormat{} },
	"xml":      func() generator.Format { return &XMLFormat{} },
}

// GetFormat returns a format by name, or the default if not found
func GetFormat(name string) generator.Format {
	if constructor, exists := formatRegistry[name]; exists {
		return constructor()
	}
	// Default to markdown if format not found
	return &MarkdownFormat{}
}

// GetFormatNames returns a list of available format names
func GetFormatNames() []string {
	names := make([]string, 0, len(formatRegistry))
	for name := range formatRegistry {
		names = append(names, name)
	}
	return names
}
