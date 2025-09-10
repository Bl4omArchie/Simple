package simple

import (
	"os"
	"fmt"
	"strings"
	"encoding/xml"
	"encoding/json"

	"gopkg.in/yaml.v3"
	"github.com/pelletier/go-toml"
	"github.com/go-playground/validator/v10"
)

// This feature can open, deserialize its content and validate tags.
// Supported file format : json, yaml, toml and xml
// Usage : LoadFileSingle[Struct type]("json", "filename", true or false for tags)

type FileFormatFactory func(data []byte, v any) error

var FileFormatRegistry = map[string]FileFormatFactory {
	"json": json.Unmarshal,
	"yaml": yaml.Unmarshal,
	"toml": toml.Unmarshal,
	"xml": xml.Unmarshal,
}


// The file content must represent a single object in the specified format
func LoadFileSingle[S any](fileFormat string, file string, validation bool) (S, error) {
	var item S

	data, err := os.ReadFile(file)
	if err != nil {
		return item, fmt.Errorf("failed to read file: %w", err)
	}

	factory, ok := FileFormatRegistry[strings.ToLower(fileFormat)]
	if !ok {
		return item, fmt.Errorf("unsupported file format: %s", fileFormat)
	}

	if err := factory(data, &item); err != nil {
		return item, fmt.Errorf("failed to parse %s: %w", fileFormat, err)
	}

	if validation {
		validate := validator.New()
		if err := validate.Struct(item); err != nil {
			return item, fmt.Errorf("validation failed : %w", err)
		}
	}

	return item, nil
}

// The file content must represent an array/list of objects in the specified format
func LoadFileMultiple[S any](fileFormat string, file string, validation bool) ([]S, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var items []S
	factory, ok := FileFormatRegistry[strings.ToLower(fileFormat)]
	if !ok {
		return nil, fmt.Errorf("unsupported file format: %s", fileFormat)
	}

	if err := factory(data, &items); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", fileFormat, err)
	}

	if validation {
		validate := validator.New()
		for i, src := range items {
			if err := validate.Struct(src); err != nil {
				return nil, fmt.Errorf("validation failed for item %d: %w", i, err)
			}
		}
	}

	return items, nil
}
