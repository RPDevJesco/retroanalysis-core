package properties

import (
	"sync"
	"time"
)

// State tracks enhanced state and history of a property
type State struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Bytes       []byte      `json:"bytes"`
	Address     uint32      `json:"address"`
	Type        string      `json:"type"`
	Frozen      bool        `json:"frozen"`
	LastChanged time.Time   `json:"last_changed"`
	LastRead    time.Time   `json:"last_read"`
	LastWrite   time.Time   `json:"last_write"`

	// Enhanced tracking
	Performance *PerformanceMetrics `json:"performance"`
	Events      []Event             `json:"events"`

	// Value history
	ValueHistory   []ValueHistoryEntry `json:"value_history"`
	MaxHistorySize uint                `json:"max_history_size"`

	// Dependencies and relationships
	Dependencies []string `json:"dependencies"`
	Dependents   []string `json:"dependents"`

	// Caching
	CachedValue interface{} `json:"cached_value,omitempty"`
	CacheExpiry *time.Time  `json:"cache_expiry,omitempty"`
	CacheValid  bool        `json:"cache_valid"`
}

// PerformanceMetrics represents detailed performance tracking
type PerformanceMetrics struct {
	ReadCount    uint64        `json:"read_count"`
	WriteCount   uint64        `json:"write_count"`
	ErrorCount   uint64        `json:"error_count"`
	AvgReadTime  time.Duration `json:"avg_read_time"`
	AvgWriteTime time.Duration `json:"avg_write_time"`
	FirstAccess  time.Time     `json:"first_access"`
	LastAccess   time.Time     `json:"last_access"`
}

// Event represents an event that occurred for a property
type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Value     interface{}            `json:"value,omitempty"`
	OldValue  interface{}            `json:"old_value,omitempty"`
	Duration  time.Duration          `json:"duration,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Source    string                 `json:"source,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

type ValueHistoryEntry struct {
	Value     interface{}   `json:"value"`
	Timestamp time.Time     `json:"timestamp"`
	Source    string        `json:"source"`
	Duration  time.Duration `json:"duration,omitempty"`
}

// Manager handles property state management
type Manager struct {
	states map[string]*State
	mu     sync.RWMutex
}

// NewManager creates a new property manager
func NewManager() *Manager {
	return &Manager{
		states: make(map[string]*State),
	}
}

// GetState returns property state
func (m *Manager) GetState(name string) *State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.states[name]
}

// UpdateState updates property state
func (m *Manager) UpdateState(name string, value interface{}, bytes []byte, address uint32) {
	m.mu.Lock()
	defer m.mu.Unlock()

	state := m.states[name]
	if state == nil {
		state = &State{
			Name:           name,
			Address:        address,
			Performance:    &PerformanceMetrics{FirstAccess: time.Now()},
			Events:         make([]Event, 0),
			ValueHistory:   make([]ValueHistoryEntry, 0),
			MaxHistorySize: 100,
		}
		m.states[name] = state
	}

	// Update state...
	state.Value = value
	state.Bytes = bytes
	state.LastRead = time.Now()
	state.Performance.ReadCount++
}
