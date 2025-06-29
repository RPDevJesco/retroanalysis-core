retroanalysis-core/
├── config/
│   └── config.go            # package config
├── properties/
│   ├── property.go          # Property definitions
│   ├── manager.go           # Property manager
│   └── validator.go         # Property validation
├── mappers/
│   ├── loader.go            # Mapper loader
│   ├── parser.go            # CUE parser
│   └── schema.go            # Schema definitions
├── transforms/
│   ├── transform.go         # Transform engine
│   └── processor.go         # Advanced processors
├── events/
│   ├── event.go             # Event system
│   └── trigger.go           # Event triggers
├── batch/
│   └── batch.go             # Batch operations
└── references/
└── reference.go         # Reference types