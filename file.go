package simple

import (
	"io"
	"os"
	"fmt"
	"path"
	"context"
	"strings"
	"archive/zip"
	"encoding/xml"
	"path/filepath"
	"encoding/json"

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

// Deserialize data from the given data type (json, yaml, toml or xml)
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

// Adapted from Gosamples website
func Unzip(ctx context.Context, source, destination string) error {
    reader, err := zip.OpenReader(source)
    if err != nil {
        return err
    }
    defer reader.Close()

    destination, err = filepath.Abs(destination)
    if err != nil {
        return err
    }

    for _, f := range reader.File {
        select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				err := unzipFile(f, destination)
        		if err != nil {
            		return err
        		}
        }
    }
    return nil
}

func unzipFile(f *zip.File, destination string) error {
    // Check if file paths are not vulnerable to Zip Slip
    filePath := filepath.Join(destination, f.Name)
    if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
        return fmt.Errorf("invalid file path: %s", filePath)
    }

    if f.FileInfo().IsDir() {
        if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
            return err
        }
        return nil
    }

    if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
        return err
    }

    destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
    if err != nil {
        return err
    }
    defer destinationFile.Close()

    zippedFile, err := f.Open()
    if err != nil {
        return err
    }
    defer zippedFile.Close()

    if _, err := io.Copy(destinationFile, zippedFile); err != nil {
        return err
    }
    return nil
}
