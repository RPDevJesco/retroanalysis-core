package mappers

import (
	"github.com/RPDevJesco/retroanalysis-core/events"
	"github.com/RPDevJesco/retroanalysis-core/properties"
	"github.com/RPDevJesco/retroanalysis-core/references"
	"github.com/RPDevJesco/retroanalysis-drivers/memory"
)

// Mapper represents an enhanced complete mapper with properties
type Mapper struct {
	Name        string
	Game        string
	Version     string
	MinVersion  string // Minimum RetroAnalysis version
	Author      string
	Description string
	Website     string
	License     string
	Metadata    *Metadata
	Platform    Platform
	Properties  map[string]*properties.Property
	Groups      map[string]*PropertyGroup
	Computed    map[string]*properties.ComputedProperty
	Constants   map[string]interface{}           // Global constants
	Preprocess  []string                         // CUE expressions run before property evaluation
	Postprocess []string                         // CUE expressions run after property evaluation
	References  map[string]*references.Reference // Reference types
	CharMaps    map[string]map[uint8]string      // Character maps
	Events      *events.Config                   // Events configuration
	Validation  *GlobalValidation                // Global validation
	Debug       *DebugConfig                     // Debug configuration
}

// Metadata represents mapper metadata
type Metadata struct {
	Created  string   `json:"created"`
	Modified string   `json:"modified"`
	Tags     []string `json:"tags"`
	Category string   `json:"category"`
	Language string   `json:"language"`
	Region   string   `json:"region"`
	Revision string   `json:"revision,omitempty"`
}

// Platform represents enhanced platform configuration
type Platform struct {
	Name          string
	Endian        string
	MemoryBlocks  []memory.Block
	Constants     map[string]interface{} // Platform-specific constants
	BaseAddresses map[string]string      // Named base addresses
	Description   string                 `json:"description,omitempty"`
	Version       string                 `json:"version,omitempty"`
	Manufacturer  string                 `json:"manufacturer,omitempty"`
	ReleaseYear   *uint                  `json:"release_year,omitempty"`
	Capabilities  *PlatformCapabilities  `json:"capabilities,omitempty"`
	Performance   *PlatformPerformance   `json:"performance,omitempty"`
}

// PropertyGroup represents enhanced property groups for UI organization
type PropertyGroup struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description,omitempty"`
	Icon        string                    `json:"icon,omitempty"`
	Properties  []string                  `json:"properties"`
	Collapsed   bool                      `json:"collapsed"`
	Color       string                    `json:"color,omitempty"`
	DisplayMode string                    `json:"display_mode,omitempty"` // "table", "cards", "tree", "custom"
	SortBy      string                    `json:"sort_by,omitempty"`
	FilterBy    string                    `json:"filter_by,omitempty"`
	Conditional *ConditionalDisplay       `json:"conditional_display,omitempty"`
	Subgroups   map[string]*PropertyGroup `json:"subgroups,omitempty"`
	Renderer    string                    `json:"custom_renderer,omitempty"`
	Priority    *uint                     `json:"priority,omitempty"`
}

// ConditionalDisplay represents conditional group display
type ConditionalDisplay struct {
	Expression   string   `json:"expression"`
	Dependencies []string `json:"dependencies"`
}

// Supporting platform structs
type PlatformCapabilities struct {
	MaxMemorySize    *uint32 `json:"max_memory_size,omitempty"`
	AddressBusWidth  *uint   `json:"address_bus_width,omitempty"`
	DataBusWidth     *uint   `json:"data_bus_width,omitempty"`
	HasMemoryMapping *bool   `json:"has_memory_mapping,omitempty"`
	SupportsBanking  *bool   `json:"supports_banking,omitempty"`
}

type PlatformPerformance struct {
	ReadLatency  *uint `json:"read_latency,omitempty"`  // milliseconds
	WriteLatency *uint `json:"write_latency,omitempty"` // milliseconds
	BatchSize    *uint `json:"batch_size,omitempty"`    // optimal batch size
}

// GlobalValidation represents global validation rules
type GlobalValidation struct {
	MemoryLayout    *MemoryLayoutValidation `json:"memory_layout,omitempty"`
	CrossValidation []CrossValidationRule   `json:"cross_validation,omitempty"`
	Performance     *PerformanceValidation  `json:"performance,omitempty"`
}

type MemoryLayoutValidation struct {
	CheckOverlaps  *bool `json:"check_overlaps,omitempty"`
	CheckBounds    *bool `json:"check_bounds,omitempty"`
	CheckAlignment *bool `json:"check_alignment,omitempty"`
}

type CrossValidationRule struct {
	Name         string   `json:"name"`
	Expression   string   `json:"expression"`
	Dependencies []string `json:"dependencies"`
	Message      string   `json:"message,omitempty"`
}

type PerformanceValidation struct {
	MaxProperties      *uint `json:"max_properties,omitempty"`
	MaxComputedDepth   *uint `json:"max_computed_depth,omitempty"`
	WarnSlowProperties *bool `json:"warn_slow_properties,omitempty"`
}

type DebugConfig struct {
	Enabled             *bool    `json:"enabled,omitempty"`
	LogLevel            string   `json:"log_level,omitempty"`
	LogProperties       []string `json:"log_properties,omitempty"`
	BenchmarkProperties []string `json:"benchmark_properties,omitempty"`
	HotReload           *bool    `json:"hot_reload,omitempty"`
	TypeChecking        *bool    `json:"type_checking,omitempty"`
	MemoryDumps         *bool    `json:"memory_dumps,omitempty"`
}
