package simple

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

type FileParser func([]byte, any) error

var FileRegistry = map[string]FileParser{
	"json": json.Unmarshal,
	"yaml": yaml.Unmarshal,
	"toml": toml.Unmarshal,
	"xml":  xml.Unmarshal,
}

func LoadFile[S any](filePath string, limit int, validation bool) ([]S, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	ext := strings.TrimPrefix(strings.ToLower(path.Ext(filePath)), ".")
	parser, ok := FileRegistry[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	var items []S
	if err := parser(data, &items); err != nil {
		var single S
		if err2 := parser(data, &single); err2 != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", ext, err)
		}
		items = append(items, single)
	}

	if validation {
		validate := validator.New()
		for i, elem := range items {
			if err := validate.Struct(elem); err != nil {
				return nil, fmt.Errorf("validation failed for element %d: %w", i, err)
			}
		}
	}

	if limit > 0 && limit < len(items) {
		items = items[:limit]
	}

	return items, nil
}
