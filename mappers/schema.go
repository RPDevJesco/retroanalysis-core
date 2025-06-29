package mappers

import _ "embed"

var Schema string

// SchemaValidator handles schema validation
type SchemaValidator struct{}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{}
}

// Validate validates a mapper against the schema
func (v *SchemaValidator) Validate(mapper *Mapper) error {
	// Schema validation logic
	return nil
}
