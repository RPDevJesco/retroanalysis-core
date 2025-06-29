package mappers

import (
	"cuelang.org/go/cue"
	"strconv"
	"strings"
)

// Parser provides utility functions for CUE parsing
type Parser struct{}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseAddress parses a hex address string to uint32
func (p *Parser) ParseAddress(addressStr string) (uint32, error) {
	if strings.HasPrefix(addressStr, "0x") || strings.HasPrefix(addressStr, "0X") {
		value, err := strconv.ParseUint(addressStr[2:], 16, 32)
		return uint32(value), err
	}

	value, err := strconv.ParseUint(addressStr, 10, 32)
	return uint32(value), err
}

// ParseAnyValue parses a CUE value into appropriate Go type
func (p *Parser) ParseAnyValue(value cue.Value) interface{} {
	if str, err := value.String(); err == nil {
		return str
	} else if num, err := value.Float64(); err == nil {
		return num
	} else if i, err := value.Int64(); err == nil {
		return i
	} else if b, err := value.Bool(); err == nil {
		return b
	} else if list, err := value.List(); err == nil {
		var arr []interface{}
		for list.Next() {
			arr = append(arr, p.ParseAnyValue(list.Value()))
		}
		return arr
	} else if fields, err := value.Fields(); err == nil {
		obj := make(map[string]interface{})
		for fields.Next() {
			obj[fields.Label()] = p.ParseAnyValue(fields.Value())
		}
		return obj
	}
	return nil
}

// Other parsing utilities...
