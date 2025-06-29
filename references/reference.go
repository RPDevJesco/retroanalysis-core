package references

import (
	"fmt"
)

// Reference represents a reference type (like enums, items, moves, etc.)
type Reference struct {
	Name        string                 `json:"name"`
	Type        ReferenceType          `json:"type"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Values      map[string]*Value      `json:"values"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ReferenceType represents the type of reference
type ReferenceType string

const (
	ReferenceTypeEnum         ReferenceType = "enum"
	ReferenceTypePokemon      ReferenceType = "pokemon"
	ReferenceTypePokemonTypes ReferenceType = "pokemon_types"
	ReferenceTypeMoves        ReferenceType = "moves"
	ReferenceTypeItems        ReferenceType = "items"
	ReferenceTypeAbilities    ReferenceType = "abilities"
	ReferenceTypeLocations    ReferenceType = "locations"
	ReferenceTypeTrainers     ReferenceType = "trainers"
	ReferenceTypeCustom       ReferenceType = "custom"
)

// Value represents a single value in a reference type
type Value struct {
	ID          uint32                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Color       string                 `json:"color,omitempty"`
	Icon        string                 `json:"icon,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Deprecated  bool                   `json:"deprecated,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PokemonSpecies represents Pokemon species reference data
type PokemonSpecies struct {
	*Value
	Type1      string    `json:"type1"`
	Type2      string    `json:"type2,omitempty"`
	BaseStats  BaseStats `json:"base_stats"`
	Height     float64   `json:"height,omitempty"` // in meters
	Weight     float64   `json:"weight,omitempty"` // in kg
	Generation uint      `json:"generation,omitempty"`
}

// BaseStats represents Pokemon base statistics
type BaseStats struct {
	HP        uint `json:"hp"`
	Attack    uint `json:"attack"`
	Defense   uint `json:"defense"`
	Speed     uint `json:"speed"`
	Special   uint `json:"special,omitempty"`    // Gen 1 special stat
	SpAttack  uint `json:"sp_attack,omitempty"`  // Gen 2+ special attack
	SpDefense uint `json:"sp_defense,omitempty"` // Gen 2+ special defense
}

// PokemonType represents Pokemon type reference data
type PokemonType struct {
	*Value
	Effectiveness map[string]float64 `json:"effectiveness,omitempty"` // type effectiveness chart
	Physical      bool               `json:"physical,omitempty"`      // true if physical type (Gen 1-3)
}

// Move represents Pokemon move reference data
type Move struct {
	*Value
	Type         string `json:"type"`
	Power        *uint  `json:"power,omitempty"`    // nil for status moves
	Accuracy     *uint  `json:"accuracy,omitempty"` // nil for moves that never miss
	PP           uint   `json:"pp"`
	Effect       string `json:"effect,omitempty"`
	EffectChance *uint  `json:"effect_chance,omitempty"` // percentage
	Priority     int    `json:"priority,omitempty"`
	Generation   uint   `json:"generation,omitempty"`
	Physical     *bool  `json:"physical,omitempty"` // nil for status, true/false for physical/special
}

// Item represents game item reference data
type Item struct {
	*Value
	Price        *uint  `json:"price,omitempty"` // nil for key items
	Effect       string `json:"effect,omitempty"`
	Consumable   bool   `json:"consumable,omitempty"`
	BattleUse    bool   `json:"battle_use,omitempty"`
	OverworldUse bool   `json:"overworld_use,omitempty"`
	Generation   uint   `json:"generation,omitempty"`
}

// ItemCategory represents item categories
type ItemCategory string

const (
	ItemCategoryMedicine  ItemCategory = "medicine"
	ItemCategoryPokeball  ItemCategory = "pokeball"
	ItemCategoryTM        ItemCategory = "tm"
	ItemCategoryBerry     ItemCategory = "berry"
	ItemCategoryKey       ItemCategory = "key"
	ItemCategoryBattle    ItemCategory = "battle"
	ItemCategoryEvolution ItemCategory = "evolution"
	ItemCategoryMisc      ItemCategory = "misc"
)

// Ability represents Pokemon ability reference data
type Ability struct {
	*Value
	Effect     string `json:"effect,omitempty"`
	Generation uint   `json:"generation,omitempty"`
	Hidden     bool   `json:"hidden,omitempty"` // true if this is a hidden ability
}

// Location represents game location reference data
type Location struct {
	*Value
	Region     string   `json:"region,omitempty"`
	Type       string   `json:"type,omitempty"`    // "route", "city", "cave", etc.
	Pokemon    []uint32 `json:"pokemon,omitempty"` // Pokemon IDs found here
	Items      []uint32 `json:"items,omitempty"`   // Item IDs found here
	Generation uint     `json:"generation,omitempty"`
}

// Trainer represents trainer reference data
type Trainer struct {
	*Value
	Class      string   `json:"class,omitempty"`    // "Gym Leader", "Elite Four", etc.
	Pokemon    []uint32 `json:"pokemon,omitempty"`  // Pokemon IDs in their party
	Location   string   `json:"location,omitempty"` // Where they're found
	Rematch    bool     `json:"rematch,omitempty"`  // Can be rematched
	Generation uint     `json:"generation,omitempty"`
}

// Manager handles reference type management
type Manager struct {
	references map[string]*Reference
}

// NewManager creates a new reference manager
func NewManager() *Manager {
	return &Manager{
		references: make(map[string]*Reference),
	}
}

// AddReference adds a reference type
func (m *Manager) AddReference(ref *Reference) {
	m.references[ref.Name] = ref
}

// GetReference gets a reference by name
func (m *Manager) GetReference(name string) (*Reference, bool) {
	ref, exists := m.references[name]
	return ref, exists
}

// GetValue gets a specific value from a reference
func (m *Manager) GetValue(referenceName string, valueID uint32) (*Value, error) {
	ref, exists := m.GetReference(referenceName)
	if !exists {
		return nil, fmt.Errorf("reference %s not found", referenceName)
	}

	for _, value := range ref.Values {
		if value.ID == valueID {
			return value, nil
		}
	}

	return nil, fmt.Errorf("value %d not found in reference %s", valueID, referenceName)
}

// GetValueByName gets a specific value by name from a reference
func (m *Manager) GetValueByName(referenceName string, valueName string) (*Value, error) {
	ref, exists := m.GetReference(referenceName)
	if !exists {
		return nil, fmt.Errorf("reference %s not found", referenceName)
	}

	if value, exists := ref.Values[valueName]; exists {
		return value, nil
	}

	return nil, fmt.Errorf("value %s not found in reference %s", valueName, referenceName)
}

// ListReferences returns all reference names
func (m *Manager) ListReferences() []string {
	names := make([]string, 0, len(m.references))
	for name := range m.references {
		names = append(names, name)
	}
	return names
}

// GetReferencesByType returns all references of a specific type
func (m *Manager) GetReferencesByType(refType ReferenceType) []*Reference {
	refs := make([]*Reference, 0)
	for _, ref := range m.references {
		if ref.Type == refType {
			refs = append(refs, ref)
		}
	}
	return refs
}

// ValidateValue validates that a value exists in a reference
func (m *Manager) ValidateValue(referenceName string, valueID uint32) bool {
	_, err := m.GetValue(referenceName, valueID)
	return err == nil
}

// GetAllReferences returns all references
func (m *Manager) GetAllReferences() map[string]*Reference {
	// Return a copy to prevent modification
	result := make(map[string]*Reference)
	for name, ref := range m.references {
		result[name] = ref
	}
	return result
}

// Factory functions for creating specialized reference types

// NewPokemonSpeciesReference creates a Pokemon species reference
func NewPokemonSpeciesReference() *Reference {
	return &Reference{
		Name:        "pokemon_species",
		Type:        ReferenceTypePokemon,
		Description: "Pokemon species data",
		Category:    "pokemon",
		Values:      make(map[string]*Value),
		Metadata:    make(map[string]interface{}),
	}
}

// NewPokemonTypesReference creates a Pokemon types reference
func NewPokemonTypesReference() *Reference {
	return &Reference{
		Name:        "pokemon_types",
		Type:        ReferenceTypePokemonTypes,
		Description: "Pokemon type data with effectiveness",
		Category:    "pokemon",
		Values:      make(map[string]*Value),
		Metadata:    make(map[string]interface{}),
	}
}

// NewMovesReference creates a moves reference
func NewMovesReference() *Reference {
	return &Reference{
		Name:        "moves",
		Type:        ReferenceTypeMoves,
		Description: "Pokemon move data",
		Category:    "pokemon",
		Values:      make(map[string]*Value),
		Metadata:    make(map[string]interface{}),
	}
}

// NewItemsReference creates an items reference
func NewItemsReference() *Reference {
	return &Reference{
		Name:        "items",
		Type:        ReferenceTypeItems,
		Description: "Game item data",
		Category:    "items",
		Values:      make(map[string]*Value),
		Metadata:    make(map[string]interface{}),
	}
}

// Helper functions for working with specialized types

// AsPokemonSpecies converts a Value to PokemonSpecies if possible
func (v *Value) AsPokemonSpecies() (*PokemonSpecies, bool) {
	if species, ok := v.Metadata["species"].(*PokemonSpecies); ok {
		return species, true
	}
	return nil, false
}

// AsPokemonType converts a Value to PokemonType if possible
func (v *Value) AsPokemonType() (*PokemonType, bool) {
	if pokemonType, ok := v.Metadata["type"].(*PokemonType); ok {
		return pokemonType, true
	}
	return nil, false
}

// AsMove converts a Value to Move if possible
func (v *Value) AsMove() (*Move, bool) {
	if move, ok := v.Metadata["move"].(*Move); ok {
		return move, true
	}
	return nil, false
}

// AsItem converts a Value to Item if possible
func (v *Value) AsItem() (*Item, bool) {
	if item, ok := v.Metadata["item"].(*Item); ok {
		return item, true
	}
	return nil, false
}
