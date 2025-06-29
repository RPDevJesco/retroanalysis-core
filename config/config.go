package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the enhanced core configuration
type Config struct {
	// Core driver configuration
	RetroArch   RetroArchConfig   `mapstructure:"retroarch"`
	BizHawk     BizHawkConfig     `mapstructure:"bizhawk"`
	Paths       PathsConfig       `mapstructure:"paths"`
	Performance PerformanceConfig `mapstructure:"performance"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	Features    FeaturesConfig    `mapstructure:"features"`

	// Enhanced core configuration sections
	PropertyMonitoring PropertyMonitoringConfig `mapstructure:"property_monitoring"`
	BatchOperations    BatchOperationsConfig    `mapstructure:"batch_operations"`
	Validation         ValidationConfig         `mapstructure:"validation"`
	Events             EventsConfig             `mapstructure:"events"`
	Memory             AdvancedMemoryConfig     `mapstructure:"memory"`
	Transforms         TransformsConfig         `mapstructure:"transforms"`
	References         ReferencesConfig         `mapstructure:"references"`
	Mappers            MappersConfig            `mapstructure:"mappers"`
}

// RetroArch driver configuration
type RetroArchConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	MaxRetries     int           `mapstructure:"max_retries"`
	RetryDelay     time.Duration `mapstructure:"retry_delay"`
	BufferSize     int           `mapstructure:"buffer_size"`
	ChunkSize      uint32        `mapstructure:"chunk_size"`
}

// BizHawk driver configuration
type BizHawkConfig struct {
	MemoryMapName string        `mapstructure:"memory_map_name"`
	DataMapName   string        `mapstructure:"data_map_name"`
	Timeout       time.Duration `mapstructure:"timeout"`
}

// Core paths configuration
type PathsConfig struct {
	MappersDir string `mapstructure:"mappers_dir"`
	DataDir    string `mapstructure:"data_dir"`
	LogDir     string `mapstructure:"log_dir"`
	CacheDir   string `mapstructure:"cache_dir"`
	TempDir    string `mapstructure:"temp_dir"`
}

// Performance configuration
type PerformanceConfig struct {
	UpdateInterval   time.Duration `mapstructure:"update_interval"`
	MemoryBufferSize int           `mapstructure:"memory_buffer_size"`
	GCTargetPercent  int           `mapstructure:"gc_target_percent"`
	MaxMemoryUsageMB int           `mapstructure:"max_memory_usage_mb"`
	WorkerPoolSize   int           `mapstructure:"worker_pool_size"`
	BatchSize        int           `mapstructure:"batch_size"`
}

// Logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

// Enhanced features configuration
type FeaturesConfig struct {
	Metrics           bool `mapstructure:"metrics"`
	Profiling         bool `mapstructure:"profiling"`
	AutoMapperReload  bool `mapstructure:"auto_mapper_reload"`
	CacheProperties   bool `mapstructure:"cache_properties"`
	BackgroundSave    bool `mapstructure:"background_save"`
	MemoryCompression bool `mapstructure:"memory_compression"`

	// Enhanced features
	AdvancedPropertyTypes bool `mapstructure:"advanced_property_types"`
	ComputedProperties    bool `mapstructure:"computed_properties"`
	PropertyFreezing      bool `mapstructure:"property_freezing"`
	RealTimeValidation    bool `mapstructure:"real_time_validation"`
	TransformPipeline     bool `mapstructure:"transform_pipeline"`
	EventSystem           bool `mapstructure:"event_system"`
	ReferenceTypes        bool `mapstructure:"reference_types"`
}

// Enhanced property monitoring configuration
type PropertyMonitoringConfig struct {
	UpdateInterval         time.Duration `mapstructure:"update_interval"`
	EnableStatistics       bool          `mapstructure:"enable_statistics"`
	HistorySize            int           `mapstructure:"history_size"`
	ChangeThreshold        float64       `mapstructure:"change_threshold"`
	EnableChangeDetection  bool          `mapstructure:"enable_change_detection"`
	BatchChangeEvents      bool          `mapstructure:"batch_change_events"`
	MaxEventsPerBatch      int           `mapstructure:"max_events_per_batch"`
	EnablePatternDetection bool          `mapstructure:"enable_pattern_detection"`
	StatisticsInterval     time.Duration `mapstructure:"statistics_interval"`
}

// Batch operations configuration
type BatchOperationsConfig struct {
	MaxBatchSize      int           `mapstructure:"max_batch_size"`
	Timeout           time.Duration `mapstructure:"timeout"`
	EnableAtomic      bool          `mapstructure:"enable_atomic"`
	ParallelExecution bool          `mapstructure:"parallel_execution"`
	ValidationMode    string        `mapstructure:"validation_mode"` // "strict", "warn", "ignore"
	RetryAttempts     int           `mapstructure:"retry_attempts"`
	RetryDelay        time.Duration `mapstructure:"retry_delay"`
}

// Enhanced validation configuration
type ValidationConfig struct {
	EnableStrict            bool          `mapstructure:"enable_strict"`
	LogValidation           bool          `mapstructure:"log_validation"`
	FailOnError             bool          `mapstructure:"fail_on_error"`
	CacheValidationResults  bool          `mapstructure:"cache_validation_results"`
	ValidationTimeout       time.Duration `mapstructure:"validation_timeout"`
	CustomValidatorsEnabled bool          `mapstructure:"custom_validators_enabled"`
	CrossValidationEnabled  bool          `mapstructure:"cross_validation_enabled"`
	AsyncValidation         bool          `mapstructure:"async_validation"`
}

// Events system configuration
type EventsConfig struct {
	Enabled             bool          `mapstructure:"enabled"`
	MaxEventHistory     int           `mapstructure:"max_event_history"`
	EventBatchSize      int           `mapstructure:"event_batch_size"`
	ProcessingTimeout   time.Duration `mapstructure:"processing_timeout"`
	EnableCustomEvents  bool          `mapstructure:"enable_custom_events"`
	LogEventTriggers    bool          `mapstructure:"log_event_triggers"`
	DefaultCooldown     time.Duration `mapstructure:"default_cooldown"`
	MaxConcurrentEvents int           `mapstructure:"max_concurrent_events"`
	EnableSequences     bool          `mapstructure:"enable_sequences"`
	EnableTimers        bool          `mapstructure:"enable_timers"`
}

// Advanced memory configuration
type AdvancedMemoryConfig struct {
	EnableCompression    bool          `mapstructure:"enable_compression"`
	CompressionAlgorithm string        `mapstructure:"compression_algorithm"` // gzip, lz4, snappy
	CacheSizeMB          int           `mapstructure:"cache_size_mb"`
	EnableMemoryMapping  bool          `mapstructure:"enable_memory_mapping"`
	PrefetchEnabled      bool          `mapstructure:"prefetch_enabled"`
	MemoryAlignment      int           `mapstructure:"memory_alignment"`
	GCInterval           time.Duration `mapstructure:"gc_interval"`
	EnableSnapshots      bool          `mapstructure:"enable_snapshots"`
	SnapshotInterval     time.Duration `mapstructure:"snapshot_interval"`
}

// Transforms configuration
type TransformsConfig struct {
	EnableCaching    bool          `mapstructure:"enable_caching"`
	CacheSize        int           `mapstructure:"cache_size"`
	CacheTTL         time.Duration `mapstructure:"cache_ttl"`
	EnablePipeline   bool          `mapstructure:"enable_pipeline"`
	MaxPipelineDepth int           `mapstructure:"max_pipeline_depth"`
	EnableAsync      bool          `mapstructure:"enable_async"`
	AsyncTimeout     time.Duration `mapstructure:"async_timeout"`
	EnableValidation bool          `mapstructure:"enable_validation"`
	LogTransforms    bool          `mapstructure:"log_transforms"`
}

// References configuration
type ReferencesConfig struct {
	EnableCaching      bool          `mapstructure:"enable_caching"`
	CacheSize          int           `mapstructure:"cache_size"`
	CacheTTL           time.Duration `mapstructure:"cache_ttl"`
	EnableValidation   bool          `mapstructure:"enable_validation"`
	AutoLoad           bool          `mapstructure:"auto_load"`
	WatchForChanges    bool          `mapstructure:"watch_for_changes"`
	DefaultLanguage    string        `mapstructure:"default_language"`
	EnableLocalization bool          `mapstructure:"enable_localization"`
}

// Mappers configuration
type MappersConfig struct {
	EnableCaching      bool          `mapstructure:"enable_caching"`
	CacheSize          int           `mapstructure:"cache_size"`
	CacheTTL           time.Duration `mapstructure:"cache_ttl"`
	EnableValidation   bool          `mapstructure:"enable_validation"`
	StrictValidation   bool          `mapstructure:"strict_validation"`
	EnableHotReload    bool          `mapstructure:"enable_hot_reload"`
	WatchForChanges    bool          `mapstructure:"watch_for_changes"`
	PreloadMappers     []string      `mapstructure:"preload_mappers"`
	EnableOptimization bool          `mapstructure:"enable_optimization"`
}

// DefaultConfig returns a configuration with enhanced sensible defaults for the core library
func DefaultConfig() *Config {
	return &Config{
		RetroArch: RetroArchConfig{
			Host:           "127.0.0.1",
			Port:           55355,
			RequestTimeout: 64 * time.Millisecond,
			MaxRetries:     3,
			RetryDelay:     100 * time.Millisecond,
			BufferSize:     1024 * 1024, // 1MB
			ChunkSize:      2048,
		},
		BizHawk: BizHawkConfig{
			MemoryMapName: "RETROANALYSIS_BIZHAWK.bin",
			DataMapName:   "RETROANALYSIS_BIZHAWK_DATA.bin",
			Timeout:       1 * time.Second,
		},
		Paths: PathsConfig{
			MappersDir: "./mappers",
			DataDir:    "./data",
			LogDir:     "./logs",
			CacheDir:   "./cache",
			TempDir:    "./temp",
		},
		Performance: PerformanceConfig{
			UpdateInterval:   5 * time.Millisecond,
			MemoryBufferSize: 1024 * 1024, // 1MB
			GCTargetPercent:  100,
			MaxMemoryUsageMB: 512,
			WorkerPoolSize:   4,
			BatchSize:        100,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			Output:     "stdout",
			MaxSize:    100, // MB
			MaxAge:     7,   // days
			MaxBackups: 3,
			Compress:   true,
		},
		Features: FeaturesConfig{
			Metrics:               true,
			Profiling:             false,
			AutoMapperReload:      true,
			CacheProperties:       true,
			BackgroundSave:        false,
			MemoryCompression:     false,
			AdvancedPropertyTypes: true,
			ComputedProperties:    true,
			PropertyFreezing:      true,
			RealTimeValidation:    true,
			TransformPipeline:     true,
			EventSystem:           true,
			ReferenceTypes:        true,
		},
		PropertyMonitoring: PropertyMonitoringConfig{
			UpdateInterval:         16 * time.Millisecond, // 60 FPS
			EnableStatistics:       true,
			HistorySize:            100000,
			ChangeThreshold:        0.01,
			EnableChangeDetection:  true,
			BatchChangeEvents:      true,
			MaxEventsPerBatch:      50,
			EnablePatternDetection: true,
			StatisticsInterval:     1 * time.Second,
		},
		BatchOperations: BatchOperationsConfig{
			MaxBatchSize:      50,
			Timeout:           5 * time.Second,
			EnableAtomic:      true,
			ParallelExecution: false,
			ValidationMode:    "strict",
			RetryAttempts:     3,
			RetryDelay:        100 * time.Millisecond,
		},
		Validation: ValidationConfig{
			EnableStrict:            true,
			LogValidation:           true,
			FailOnError:             false,
			CacheValidationResults:  true,
			ValidationTimeout:       1 * time.Second,
			CustomValidatorsEnabled: true,
			CrossValidationEnabled:  true,
			AsyncValidation:         false,
		},
		Events: EventsConfig{
			Enabled:             true,
			MaxEventHistory:     1000,
			EventBatchSize:      10,
			ProcessingTimeout:   5 * time.Second,
			EnableCustomEvents:  true,
			LogEventTriggers:    true,
			DefaultCooldown:     100 * time.Millisecond,
			MaxConcurrentEvents: 20,
			EnableSequences:     true,
			EnableTimers:        true,
		},
		Memory: AdvancedMemoryConfig{
			EnableCompression:    false,
			CompressionAlgorithm: "gzip",
			CacheSizeMB:          64,
			EnableMemoryMapping:  false,
			PrefetchEnabled:      true,
			MemoryAlignment:      4,
			GCInterval:           30 * time.Second,
			EnableSnapshots:      false,
			SnapshotInterval:     5 * time.Minute,
		},
		Transforms: TransformsConfig{
			EnableCaching:    true,
			CacheSize:        1000,
			CacheTTL:         5 * time.Minute,
			EnablePipeline:   true,
			MaxPipelineDepth: 10,
			EnableAsync:      false,
			AsyncTimeout:     1 * time.Second,
			EnableValidation: true,
			LogTransforms:    false,
		},
		References: ReferencesConfig{
			EnableCaching:      true,
			CacheSize:          1000,
			CacheTTL:           10 * time.Minute,
			EnableValidation:   true,
			AutoLoad:           true,
			WatchForChanges:    true,
			DefaultLanguage:    "en",
			EnableLocalization: false,
		},
		Mappers: MappersConfig{
			EnableCaching:      true,
			CacheSize:          10,
			CacheTTL:           30 * time.Minute,
			EnableValidation:   true,
			StrictValidation:   true,
			EnableHotReload:    true,
			WatchForChanges:    true,
			PreloadMappers:     []string{},
			EnableOptimization: true,
		},
	}
}

// LoadConfig loads enhanced configuration from files and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	v := viper.New()

	// Set default values
	setDefaults(v, config)

	// Configure viper
	v.SetConfigName("retroanalysis-core")
	v.SetConfigType("yaml")

	// Add config paths
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.retroanalysis")
	v.AddConfigPath("/etc/retroanalysis")

	// Environment variables
	v.SetEnvPrefix("RETROANALYSIS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is not an error, we'll use defaults
	}

	// Unmarshal to struct
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate and normalize configuration
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	if err := NormalizeConfig(config); err != nil {
		return nil, fmt.Errorf("error normalizing config: %w", err)
	}

	return config, nil
}

// setDefaults sets all default values in viper
func setDefaults(v *viper.Viper, config *Config) {
	// RetroArch defaults
	v.SetDefault("retroarch.host", config.RetroArch.Host)
	v.SetDefault("retroarch.port", config.RetroArch.Port)
	v.SetDefault("retroarch.request_timeout", config.RetroArch.RequestTimeout)
	v.SetDefault("retroarch.max_retries", config.RetroArch.MaxRetries)
	v.SetDefault("retroarch.retry_delay", config.RetroArch.RetryDelay)
	v.SetDefault("retroarch.buffer_size", config.RetroArch.BufferSize)
	v.SetDefault("retroarch.chunk_size", config.RetroArch.ChunkSize)

	// BizHawk defaults
	v.SetDefault("bizhawk.memory_map_name", config.BizHawk.MemoryMapName)
	v.SetDefault("bizhawk.data_map_name", config.BizHawk.DataMapName)
	v.SetDefault("bizhawk.timeout", config.BizHawk.Timeout)

	// Paths defaults
	v.SetDefault("paths.mappers_dir", config.Paths.MappersDir)
	v.SetDefault("paths.data_dir", config.Paths.DataDir)
	v.SetDefault("paths.log_dir", config.Paths.LogDir)
	v.SetDefault("paths.cache_dir", config.Paths.CacheDir)
	v.SetDefault("paths.temp_dir", config.Paths.TempDir)

	// Performance defaults
	v.SetDefault("performance.update_interval", config.Performance.UpdateInterval)
	v.SetDefault("performance.memory_buffer_size", config.Performance.MemoryBufferSize)
	v.SetDefault("performance.gc_target_percent", config.Performance.GCTargetPercent)
	v.SetDefault("performance.max_memory_usage_mb", config.Performance.MaxMemoryUsageMB)
	v.SetDefault("performance.worker_pool_size", config.Performance.WorkerPoolSize)
	v.SetDefault("performance.batch_size", config.Performance.BatchSize)

	// Logging defaults
	v.SetDefault("logging.level", config.Logging.Level)
	v.SetDefault("logging.format", config.Logging.Format)
	v.SetDefault("logging.output", config.Logging.Output)
	v.SetDefault("logging.max_size", config.Logging.MaxSize)
	v.SetDefault("logging.max_age", config.Logging.MaxAge)
	v.SetDefault("logging.max_backups", config.Logging.MaxBackups)
	v.SetDefault("logging.compress", config.Logging.Compress)

	// Features defaults
	v.SetDefault("features.metrics", config.Features.Metrics)
	v.SetDefault("features.profiling", config.Features.Profiling)
	v.SetDefault("features.auto_mapper_reload", config.Features.AutoMapperReload)
	v.SetDefault("features.cache_properties", config.Features.CacheProperties)
	v.SetDefault("features.background_save", config.Features.BackgroundSave)
	v.SetDefault("features.memory_compression", config.Features.MemoryCompression)
	v.SetDefault("features.advanced_property_types", config.Features.AdvancedPropertyTypes)
	v.SetDefault("features.computed_properties", config.Features.ComputedProperties)
	v.SetDefault("features.property_freezing", config.Features.PropertyFreezing)
	v.SetDefault("features.real_time_validation", config.Features.RealTimeValidation)
	v.SetDefault("features.transform_pipeline", config.Features.TransformPipeline)
	v.SetDefault("features.event_system", config.Features.EventSystem)
	v.SetDefault("features.reference_types", config.Features.ReferenceTypes)

	// Property monitoring defaults
	v.SetDefault("property_monitoring.update_interval", config.PropertyMonitoring.UpdateInterval)
	v.SetDefault("property_monitoring.enable_statistics", config.PropertyMonitoring.EnableStatistics)
	v.SetDefault("property_monitoring.history_size", config.PropertyMonitoring.HistorySize)
	v.SetDefault("property_monitoring.change_threshold", config.PropertyMonitoring.ChangeThreshold)
	v.SetDefault("property_monitoring.enable_change_detection", config.PropertyMonitoring.EnableChangeDetection)
	v.SetDefault("property_monitoring.batch_change_events", config.PropertyMonitoring.BatchChangeEvents)
	v.SetDefault("property_monitoring.max_events_per_batch", config.PropertyMonitoring.MaxEventsPerBatch)
	v.SetDefault("property_monitoring.enable_pattern_detection", config.PropertyMonitoring.EnablePatternDetection)
	v.SetDefault("property_monitoring.statistics_interval", config.PropertyMonitoring.StatisticsInterval)

	// Batch operations defaults
	v.SetDefault("batch_operations.max_batch_size", config.BatchOperations.MaxBatchSize)
	v.SetDefault("batch_operations.timeout", config.BatchOperations.Timeout)
	v.SetDefault("batch_operations.enable_atomic", config.BatchOperations.EnableAtomic)
	v.SetDefault("batch_operations.parallel_execution", config.BatchOperations.ParallelExecution)
	v.SetDefault("batch_operations.validation_mode", config.BatchOperations.ValidationMode)
	v.SetDefault("batch_operations.retry_attempts", config.BatchOperations.RetryAttempts)
	v.SetDefault("batch_operations.retry_delay", config.BatchOperations.RetryDelay)

	// Validation defaults
	v.SetDefault("validation.enable_strict", config.Validation.EnableStrict)
	v.SetDefault("validation.log_validation", config.Validation.LogValidation)
	v.SetDefault("validation.fail_on_error", config.Validation.FailOnError)
	v.SetDefault("validation.cache_validation_results", config.Validation.CacheValidationResults)
	v.SetDefault("validation.validation_timeout", config.Validation.ValidationTimeout)
	v.SetDefault("validation.custom_validators_enabled", config.Validation.CustomValidatorsEnabled)
	v.SetDefault("validation.cross_validation_enabled", config.Validation.CrossValidationEnabled)
	v.SetDefault("validation.async_validation", config.Validation.AsyncValidation)

	// Events defaults
	v.SetDefault("events.enabled", config.Events.Enabled)
	v.SetDefault("events.max_event_history", config.Events.MaxEventHistory)
	v.SetDefault("events.event_batch_size", config.Events.EventBatchSize)
	v.SetDefault("events.processing_timeout", config.Events.ProcessingTimeout)
	v.SetDefault("events.enable_custom_events", config.Events.EnableCustomEvents)
	v.SetDefault("events.log_event_triggers", config.Events.LogEventTriggers)
	v.SetDefault("events.default_cooldown", config.Events.DefaultCooldown)
	v.SetDefault("events.max_concurrent_events", config.Events.MaxConcurrentEvents)
	v.SetDefault("events.enable_sequences", config.Events.EnableSequences)
	v.SetDefault("events.enable_timers", config.Events.EnableTimers)

	// Memory defaults
	v.SetDefault("memory.enable_compression", config.Memory.EnableCompression)
	v.SetDefault("memory.compression_algorithm", config.Memory.CompressionAlgorithm)
	v.SetDefault("memory.cache_size_mb", config.Memory.CacheSizeMB)
	v.SetDefault("memory.enable_memory_mapping", config.Memory.EnableMemoryMapping)
	v.SetDefault("memory.prefetch_enabled", config.Memory.PrefetchEnabled)
	v.SetDefault("memory.memory_alignment", config.Memory.MemoryAlignment)
	v.SetDefault("memory.gc_interval", config.Memory.GCInterval)
	v.SetDefault("memory.enable_snapshots", config.Memory.EnableSnapshots)
	v.SetDefault("memory.snapshot_interval", config.Memory.SnapshotInterval)

	// Transforms defaults
	v.SetDefault("transforms.enable_caching", config.Transforms.EnableCaching)
	v.SetDefault("transforms.cache_size", config.Transforms.CacheSize)
	v.SetDefault("transforms.cache_ttl", config.Transforms.CacheTTL)
	v.SetDefault("transforms.enable_pipeline", config.Transforms.EnablePipeline)
	v.SetDefault("transforms.max_pipeline_depth", config.Transforms.MaxPipelineDepth)
	v.SetDefault("transforms.enable_async", config.Transforms.EnableAsync)
	v.SetDefault("transforms.async_timeout", config.Transforms.AsyncTimeout)
	v.SetDefault("transforms.enable_validation", config.Transforms.EnableValidation)
	v.SetDefault("transforms.log_transforms", config.Transforms.LogTransforms)

	// References defaults
	v.SetDefault("references.enable_caching", config.References.EnableCaching)
	v.SetDefault("references.cache_size", config.References.CacheSize)
	v.SetDefault("references.cache_ttl", config.References.CacheTTL)
	v.SetDefault("references.enable_validation", config.References.EnableValidation)
	v.SetDefault("references.auto_load", config.References.AutoLoad)
	v.SetDefault("references.watch_for_changes", config.References.WatchForChanges)
	v.SetDefault("references.default_language", config.References.DefaultLanguage)
	v.SetDefault("references.enable_localization", config.References.EnableLocalization)

	// Mappers defaults
	v.SetDefault("mappers.enable_caching", config.Mappers.EnableCaching)
	v.SetDefault("mappers.cache_size", config.Mappers.CacheSize)
	v.SetDefault("mappers.cache_ttl", config.Mappers.CacheTTL)
	v.SetDefault("mappers.enable_validation", config.Mappers.EnableValidation)
	v.SetDefault("mappers.strict_validation", config.Mappers.StrictValidation)
	v.SetDefault("mappers.enable_hot_reload", config.Mappers.EnableHotReload)
	v.SetDefault("mappers.watch_for_changes", config.Mappers.WatchForChanges)
	v.SetDefault("mappers.preload_mappers", config.Mappers.PreloadMappers)
	v.SetDefault("mappers.enable_optimization", config.Mappers.EnableOptimization)
}

// ValidateConfig performs comprehensive validation of the configuration
func ValidateConfig(config *Config) error {
	// Core validation
	if err := validateCoreConfig(config); err != nil {
		return fmt.Errorf("core configuration validation failed: %w", err)
	}

	// Enhanced feature validation
	if err := validateEnhancedConfig(config); err != nil {
		return fmt.Errorf("enhanced configuration validation failed: %w", err)
	}

	// Cross-validation between sections
	if err := validateCrossReferences(config); err != nil {
		return fmt.Errorf("cross-reference validation failed: %w", err)
	}

	return nil
}

func validateCoreConfig(config *Config) error {
	// Validate port ranges
	if config.RetroArch.Port < 1 || config.RetroArch.Port > 65535 {
		return fmt.Errorf("invalid RetroArch port: %d", config.RetroArch.Port)
	}

	// Validate timeouts
	if config.Performance.UpdateInterval < time.Millisecond {
		return fmt.Errorf("update interval too small: %v", config.Performance.UpdateInterval)
	}

	if config.RetroArch.RequestTimeout < time.Millisecond {
		return fmt.Errorf("RetroArch request timeout too small: %v", config.RetroArch.RequestTimeout)
	}

	return nil
}

func validateEnhancedConfig(config *Config) error {
	// Validate enhanced features are compatible
	if config.Features.PropertyFreezing && !config.Features.AdvancedPropertyTypes {
		return fmt.Errorf("property freezing requires advanced property types to be enabled")
	}

	if config.Events.Enabled && config.Events.ProcessingTimeout < time.Millisecond {
		return fmt.Errorf("event processing timeout must be at least 1ms when events are enabled")
	}

	if config.PropertyMonitoring.BatchChangeEvents && config.PropertyMonitoring.MaxEventsPerBatch < 1 {
		return fmt.Errorf("max events per batch must be at least 1 when batch change events are enabled")
	}

	// Validate validation mode
	validModes := []string{"strict", "warn", "ignore"}
	if !contains(validModes, config.BatchOperations.ValidationMode) {
		return fmt.Errorf("invalid batch operations validation mode: %s", config.BatchOperations.ValidationMode)
	}

	// Validate compression algorithm
	if config.Memory.EnableCompression {
		validAlgorithms := []string{"gzip", "lz4", "snappy"}
		if !contains(validAlgorithms, config.Memory.CompressionAlgorithm) {
			return fmt.Errorf("invalid memory compression algorithm: %s", config.Memory.CompressionAlgorithm)
		}
	}

	// Validate logging level
	validLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLevels, config.Logging.Level) {
		return fmt.Errorf("invalid logging level: %s", config.Logging.Level)
	}

	return nil
}

func validateCrossReferences(config *Config) error {
	// Validate that related settings are compatible
	if config.PropertyMonitoring.UpdateInterval > config.Performance.UpdateInterval {
		return fmt.Errorf("property monitoring update interval cannot be longer than performance update interval")
	}

	return nil
}

// NormalizeConfig validates and normalizes enhanced configuration values
func NormalizeConfig(config *Config) error {
	// Convert relative paths to absolute paths
	var err error

	if config.Paths.MappersDir, err = filepath.Abs(config.Paths.MappersDir); err != nil {
		return fmt.Errorf("invalid mappers directory: %w", err)
	}

	if config.Paths.DataDir, err = filepath.Abs(config.Paths.DataDir); err != nil {
		return fmt.Errorf("invalid data directory: %w", err)
	}

	if config.Paths.LogDir, err = filepath.Abs(config.Paths.LogDir); err != nil {
		return fmt.Errorf("invalid log directory: %w", err)
	}

	if config.Paths.CacheDir, err = filepath.Abs(config.Paths.CacheDir); err != nil {
		return fmt.Errorf("invalid cache directory: %w", err)
	}

	if config.Paths.TempDir, err = filepath.Abs(config.Paths.TempDir); err != nil {
		return fmt.Errorf("invalid temp directory: %w", err)
	}

	// Create directories if they don't exist
	dirs := []string{
		config.Paths.MappersDir,
		config.Paths.DataDir,
		config.Paths.LogDir,
		config.Paths.CacheDir,
		config.Paths.TempDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetConfigSummary returns a summary of the current configuration for logging
func GetConfigSummary(config *Config) map[string]interface{} {
	return map[string]interface{}{
		"retroarch_host":      fmt.Sprintf("%s:%d", config.RetroArch.Host, config.RetroArch.Port),
		"update_interval":     config.Performance.UpdateInterval.String(),
		"property_monitoring": config.PropertyMonitoring.UpdateInterval.String(),
		"features_enabled":    countEnabledFeatures(config),
		"enhanced_features":   config.Features.AdvancedPropertyTypes,
		"event_system":        config.Events.Enabled,
		"validation_enabled":  config.Validation.EnableStrict,
		"memory_cache_mb":     config.Memory.CacheSizeMB,
		"batch_max_size":      config.BatchOperations.MaxBatchSize,
	}
}

func countEnabledFeatures(config *Config) int {
	count := 0
	if config.Features.Metrics {
		count++
	}
	if config.Features.Profiling {
		count++
	}
	if config.Features.AutoMapperReload {
		count++
	}
	if config.Features.CacheProperties {
		count++
	}
	if config.Features.BackgroundSave {
		count++
	}
	if config.Features.MemoryCompression {
		count++
	}
	if config.Features.AdvancedPropertyTypes {
		count++
	}
	if config.Features.ComputedProperties {
		count++
	}
	if config.Features.PropertyFreezing {
		count++
	}
	if config.Features.RealTimeValidation {
		count++
	}
	if config.Features.TransformPipeline {
		count++
	}
	if config.Features.EventSystem {
		count++
	}
	if config.Features.ReferenceTypes {
		count++
	}
	return count
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Environment helpers

// GetEnvOrDefault gets an environment variable or returns a default value
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsProduction returns true if running in production environment
func IsProduction() bool {
	env := strings.ToLower(GetEnvOrDefault("RETROANALYSIS_ENV", "development"))
	return env == "production" || env == "prod"
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	env := strings.ToLower(GetEnvOrDefault("RETROANALYSIS_ENV", "development"))
	return env == "development" || env == "dev"
}
