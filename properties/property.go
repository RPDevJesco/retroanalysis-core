package properties

// PropertyType represents the enhanced property types
type PropertyType string

const (
	// Basic types
	PropertyTypeUint8    PropertyType = "uint8"
	PropertyTypeUint16   PropertyType = "uint16"
	PropertyTypeUint32   PropertyType = "uint32"
	PropertyTypeInt8     PropertyType = "int8"
	PropertyTypeInt16    PropertyType = "int16"
	PropertyTypeInt32    PropertyType = "int32"
	PropertyTypeString   PropertyType = "string"
	PropertyTypeBool     PropertyType = "bool"
	PropertyTypeBitfield PropertyType = "bitfield"
	PropertyTypeBCD      PropertyType = "bcd"

	// Advanced types
	PropertyTypeBit        PropertyType = "bit"
	PropertyTypeNibble     PropertyType = "nibble"
	PropertyTypeFloat32    PropertyType = "float32"
	PropertyTypeFloat64    PropertyType = "float64"
	PropertyTypePointer    PropertyType = "pointer"
	PropertyTypeArray      PropertyType = "array"
	PropertyTypeStruct     PropertyType = "struct"
	PropertyTypeEnum       PropertyType = "enum"
	PropertyTypeFlags      PropertyType = "flags"
	PropertyTypeTime       PropertyType = "time"
	PropertyTypeVersion    PropertyType = "version"
	PropertyTypeChecksum   PropertyType = "checksum"
	PropertyTypeCoordinate PropertyType = "coordinate"
	PropertyTypeColor      PropertyType = "color"
	PropertyTypePercentage PropertyType = "percentage"
)

// Property represents an enhanced property definition
type Property struct {
	Name        string
	Type        PropertyType
	Address     uint32
	Length      uint32
	Position    *uint32 // for bit/nibble properties
	Size        *uint32 // element size for arrays/structs
	Endian      string
	Description string
	ReadOnly    bool
	CharMap     map[uint8]string

	// Enhanced features
	UIHints     *UIHints
	Advanced    *AdvancedConfig
	Performance *PerformanceHints
	Debug       *DebugConfig

	// Freezing support
	Freezable     bool
	DefaultFrozen bool
	Frozen        bool
	FrozenData    []byte

	// Computed properties
	DependsOn []string
	Computed  *ComputedProperty

	// Custom expressions
	ReadExpression  string
	WriteExpression string

	// Transform reference (handled by transforms package)
	TransformConfig interface{} // Will be populated by transforms package
}

// UIHints represents enhanced UI presentation hints
type UIHints struct {
	DisplayFormat string `json:"display_format,omitempty"`
	Unit          string `json:"unit,omitempty"`
	Precision     *uint  `json:"precision,omitempty"`
	ShowInList    *bool  `json:"show_in_list,omitempty"`
	Category      string `json:"category,omitempty"`
	Priority      *uint  `json:"priority,omitempty"`
	Tooltip       string `json:"tooltip,omitempty"`
	Color         string `json:"color,omitempty"`
	Icon          string `json:"icon,omitempty"`
	Badge         string `json:"badge,omitempty"`
	Editable      *bool  `json:"editable,omitempty"`
	Copyable      *bool  `json:"copyable,omitempty"`
	Watchable     *bool  `json:"watchable,omitempty"`
	Chartable     *bool  `json:"chartable,omitempty"`
	ChartType     string `json:"chart_type,omitempty"`
	ChartColor    string `json:"chart_color,omitempty"`
}

// ComputedProperty represents a property derived from other properties
type ComputedProperty struct {
	Expression        string       `json:"expression"`
	Dependencies      []string     `json:"dependencies"`
	Type              PropertyType `json:"type,omitempty"`
	Cached            *bool        `json:"cached,omitempty"`
	CacheInvalidation []string     `json:"cache_invalidation,omitempty"`
}

// AdvancedConfig represents type-specific advanced configuration
type AdvancedConfig struct {
	// Pointer types
	TargetType      *PropertyType `json:"target_type,omitempty"`
	MaxDereferences *uint         `json:"max_dereferences,omitempty"`
	NullValue       *uint32       `json:"null_value,omitempty"`

	// Array types
	ElementType    *PropertyType `json:"element_type,omitempty"`
	ElementSize    *uint         `json:"element_size,omitempty"`
	DynamicLength  *bool         `json:"dynamic_length,omitempty"`
	LengthProperty string        `json:"length_property,omitempty"`
	MaxElements    *uint         `json:"max_elements,omitempty"`
	IndexOffset    *uint         `json:"index_offset,omitempty"`
	Stride         *uint         `json:"stride,omitempty"`

	// Struct types
	Fields  map[string]*StructField `json:"fields,omitempty"`
	Extends string                  `json:"extends,omitempty"`

	// Enum types
	EnumValues         map[string]*EnumValue `json:"enum_values,omitempty"`
	AllowUnknownValues *bool                 `json:"allow_unknown_values,omitempty"`
	DefaultValue       *uint32               `json:"default_value,omitempty"`

	// Flags/bitfield types
	FlagDefinitions map[string]*FlagDefinition `json:"flag_definitions,omitempty"`

	// Time types
	TimeFormat string   `json:"time_format,omitempty"`
	FrameRate  *float64 `json:"frame_rate,omitempty"`
	Epoch      string   `json:"epoch,omitempty"`

	// Coordinate types
	CoordinateSystem string `json:"coordinate_system,omitempty"`
	Dimensions       *uint  `json:"dimensions,omitempty"`
	Units            string `json:"units,omitempty"`

	// Color types
	ColorFormat  string `json:"color_format,omitempty"`
	AlphaChannel *bool  `json:"alpha_channel,omitempty"`
	PaletteRef   string `json:"palette_ref,omitempty"`

	// Percentage types
	MaxValue  *float64 `json:"max_value,omitempty"`
	Precision *uint    `json:"precision,omitempty"`

	// Version types
	VersionFormat string `json:"version_format,omitempty"`

	// Checksum types
	ChecksumAlgorithm string         `json:"checksum_algorithm,omitempty"`
	ChecksumRange     *ChecksumRange `json:"checksum_range,omitempty"`
}

// Supporting structs for AdvancedConfig...
type StructField struct {
	Type            PropertyType      `json:"type"`
	Offset          uint              `json:"offset"`
	Size            *uint             `json:"size,omitempty"`
	TransformConfig interface{}       `json:"transform,omitempty"` // Populated by transforms package
	Validation      *Validation       `json:"validation,omitempty"`
	Description     string            `json:"description,omitempty"`
	Computed        *ComputedProperty `json:"computed,omitempty"`
}

type EnumValue struct {
	Value       uint32 `json:"value"`
	Description string `json:"description"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Deprecated  *bool  `json:"deprecated,omitempty"`
}

type FlagDefinition struct {
	Bit               uint     `json:"bit"`
	Description       string   `json:"description"`
	InvertLogic       *bool    `json:"invert_logic,omitempty"`
	Group             string   `json:"group,omitempty"`
	MutuallyExclusive []string `json:"mutually_exclusive,omitempty"`
}

type ChecksumRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type PerformanceHints struct {
	Cacheable     *bool  `json:"cacheable,omitempty"`
	CacheTimeout  *uint  `json:"cache_timeout,omitempty"`
	ReadFrequency string `json:"read_frequency,omitempty"`
	Critical      *bool  `json:"critical,omitempty"`
	Batchable     *bool  `json:"batchable,omitempty"`
	BatchGroup    string `json:"batch_group,omitempty"`
}

type DebugConfig struct {
	LogReads        *bool  `json:"log_reads,omitempty"`
	LogWrites       *bool  `json:"log_writes,omitempty"`
	BreakOnRead     *bool  `json:"break_on_read,omitempty"`
	BreakOnWrite    *bool  `json:"break_on_write,omitempty"`
	WatchExpression string `json:"watch_expression,omitempty"`
}
