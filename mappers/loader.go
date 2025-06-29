package mappers

import (
	"fmt"
	"github.com/RPDevJesco/retroanalysis-core/properties"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
)

// Loader handles loading and parsing enhanced CUE mapper files
type Loader struct {
	mappersDir string
	mappers    map[string]*Mapper
}

// NewLoader creates a new enhanced mapper loader
func NewLoader(mappersDir string) *Loader {
	return &Loader{
		mappersDir: mappersDir,
		mappers:    make(map[string]*Mapper),
	}
}

// List returns available mapper names
func (l *Loader) List() []string {
	names := make([]string, 0)

	err := filepath.WalkDir(l.mappersDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".cue") {
			relPath, err := filepath.Rel(l.mappersDir, path)
			if err != nil {
				return nil
			}

			mapperName := strings.TrimSuffix(relPath, ".cue")
			mapperName = strings.ReplaceAll(mapperName, "\\", "/")
			names = append(names, mapperName)
		}

		return nil
	})

	if err != nil {
		return []string{}
	}

	return names
}

// Load loads a mapper by name
func (l *Loader) Load(name string) (*Mapper, error) {
	if mapper, exists := l.mappers[name]; exists {
		return mapper, nil
	}

	filePath := filepath.Join(l.mappersDir, name+".cue")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("mapper file not found: %s", filePath)
	}

	mapper, err := l.loadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	if mapper.Name == "" {
		mapper.Name = name
	}

	l.mappers[name] = mapper
	return mapper, nil
}

// loadFromFile loads a mapper from a CUE file
func (l *Loader) loadFromFile(filePath string) (*Mapper, error) {
	log.Printf("🔍 Loading enhanced mapper from file: %s", filePath)

	ctx := cuecontext.New()

	buildInstances := load.Instances([]string{filePath}, &load.Config{})
	if len(buildInstances) == 0 {
		return nil, fmt.Errorf("no CUE instances found in %s", filePath)
	}

	inst := buildInstances[0]
	if inst.Err != nil {
		log.Printf("❌ CUE load error for %s: %v", filePath, inst.Err)
		return nil, fmt.Errorf("CUE load error: %w", inst.Err)
	}

	log.Printf("✅ CUE instance loaded, building value...")
	value := ctx.BuildInstance(inst)
	if value.Err() != nil {
		log.Printf("❌ CUE build error for %s: %v", filePath, value.Err())
		if err := value.Validate(); err != nil {
			log.Printf("❌ CUE validation errors:")
			for _, e := range errors.Errors(err) {
				log.Printf("   • %s", e)
			}
		}
		return nil, fmt.Errorf("CUE build error: %w", value.Err())
	}

	log.Printf("✅ CUE value built successfully, parsing enhanced mapper...")
	mapper, err := l.parseEnhancedMapper(value)
	if err != nil {
		log.Printf("❌ Enhanced mapper parsing error: %v", err)
		return nil, err
	}

	log.Printf("✅ Enhanced mapper parsed: %d properties, %d groups, %d computed, %d references",
		len(mapper.Properties), len(mapper.Groups), len(mapper.Computed), len(mapper.References))

	return mapper, nil
}

// parseEnhancedMapper parses a CUE value into an enhanced Mapper struct
func (l *Loader) parseEnhancedMapper(value cue.Value) (*Mapper, error) {
	mapper := &Mapper{
		Properties: make(map[string]*properties.Property),
		Groups:     make(map[string]*PropertyGroup),
		Computed:   make(map[string]*properties.ComputedProperty),
		Constants:  make(map[string]interface{}),
		References: make(map[string]*references.Reference),
		CharMaps:   make(map[string]map[uint8]string),
	}

	// Parse each section...
	// This would contain all the parsing logic from your original loader.go
	// but broken down into smaller methods

	return mapper, nil
}
