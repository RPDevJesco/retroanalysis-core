package events

import (
	"fmt"
	"sync"
	"time"
)

// Event represents an enhanced event definition
type Event struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	Trigger      string                 `json:"trigger"`      // CUE expression for when to trigger
	Action       string                 `json:"action"`       // CUE expression for what to do
	Dependencies []string               `json:"dependencies"` // Properties this event depends on
	Enabled      bool                   `json:"enabled"`
	Priority     uint                   `json:"priority,omitempty"`
	Cooldown     *time.Duration         `json:"cooldown,omitempty"`     // Minimum time between triggers
	MaxTriggers  *uint                  `json:"max_triggers,omitempty"` // Maximum number of times to trigger
	Metadata     map[string]interface{} `json:"metadata,omitempty"`

	// Internal state
	TriggerCount    uint      `json:"trigger_count"`
	LastTriggered   time.Time `json:"last_triggered,omitempty"`
	LastEvaluated   time.Time `json:"last_evaluated,omitempty"`
	CurrentlyActive bool      `json:"currently_active"`
}

// EventType represents different types of events
type EventType string

const (
	EventTypePropertyChanged EventType = "property_changed"
	EventTypeValueThreshold  EventType = "value_threshold"
	EventTypeTimer           EventType = "timer"
	EventTypeSequence        EventType = "sequence"
	EventTypeCustom          EventType = "custom"
	EventTypeSystem          EventType = "system"
)

// Config represents the events configuration for a mapper
type Config struct {
	OnLoad            string                  `json:"on_load,omitempty"`             // CUE expression run when mapper loads
	OnUnload          string                  `json:"on_unload,omitempty"`           // CUE expression run when mapper unloads
	OnPropertyChanged string                  `json:"on_property_changed,omitempty"` // CUE expression run when any property changes
	Custom            map[string]*CustomEvent `json:"custom,omitempty"`              // Custom event definitions
	Global            *GlobalEventConfig      `json:"global,omitempty"`              // Global event settings
}

// CustomEvent represents a custom event definition from mapper
type CustomEvent struct {
	Trigger      string   `json:"trigger"`
	Action       string   `json:"action"`
	Dependencies []string `json:"dependencies"`
	Description  string   `json:"description,omitempty"`
	Enabled      bool     `json:"enabled"`
	Priority     uint     `json:"priority,omitempty"`
}

// GlobalEventConfig represents global event system settings
type GlobalEventConfig struct {
	MaxEventHistory    int           `json:"max_event_history"`
	EventBatchSize     int           `json:"event_batch_size"`
	ProcessingTimeout  time.Duration `json:"processing_timeout"`
	EnableCustomEvents bool          `json:"enable_custom_events"`
	LogEventTriggers   bool          `json:"log_event_triggers"`
	DefaultCooldown    time.Duration `json:"default_cooldown,omitempty"`
}

// HistoryEntry represents an event in the history
type HistoryEntry struct {
	EventName string                 `json:"event_name"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Triggered bool                   `json:"triggered"`
	Source    string                 `json:"source"`   // "automatic", "manual", "api"
	Duration  time.Duration          `json:"duration"` // How long the event took to process
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"` // Property values at trigger time
}

// Result represents the result of event execution
type Result struct {
	Success  bool                   `json:"success"`
	Duration time.Duration          `json:"duration"`
	Output   interface{}            `json:"output,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Condition represents an event trigger condition
type Condition struct {
	Expression   string   `json:"expression"`   // CUE expression to evaluate
	Dependencies []string `json:"dependencies"` // Properties this condition depends on
	Description  string   `json:"description,omitempty"`
}

// Action represents an event action
type Action struct {
	Expression  string                 `json:"expression"` // CUE expression to execute
	Description string                 `json:"description,omitempty"`
	Async       bool                   `json:"async"` // Whether to execute asynchronously
	Timeout     time.Duration          `json:"timeout,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Manager handles event management and execution
type Manager struct {
	events    map[string]*Event
	history   []HistoryEntry
	active    []string
	config    *Config
	mu        sync.RWMutex
	historyMu sync.RWMutex

	// Event processing
	enabled         bool
	maxHistory      int
	processingQueue chan *ProcessRequest
	stopChan        chan bool
	running         bool
}

// ProcessRequest represents a request to process an event
type ProcessRequest struct {
	EventName string
	Force     bool
	Context   map[string]interface{}
	Source    string
	Response  chan *Result
}

// NewManager creates a new event manager
func NewManager(config *Config) *Manager {
	maxHistory := 1000
	if config != nil && config.Global != nil && config.Global.MaxEventHistory > 0 {
		maxHistory = config.Global.MaxEventHistory
	}

	return &Manager{
		events:          make(map[string]*Event),
		history:         make([]HistoryEntry, 0),
		active:          make([]string, 0),
		config:          config,
		enabled:         true,
		maxHistory:      maxHistory,
		processingQueue: make(chan *ProcessRequest, 100),
		stopChan:        make(chan bool),
	}
}

// Start starts the event manager
func (m *Manager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return
	}

	m.running = true
	go m.processEvents()
}

// Stop stops the event manager
func (m *Manager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false
	m.stopChan <- true
}

// AddEvent adds an event to the manager
func (m *Manager) AddEvent(event *Event) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events[event.Name] = event
}

// RemoveEvent removes an event from the manager
func (m *Manager) RemoveEvent(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.events, name)
}

// GetEvent gets an event by name
func (m *Manager) GetEvent(name string) (*Event, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	event, exists := m.events[name]
	return event, exists
}

// ListEvents returns all event names
func (m *Manager) ListEvents() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.events))
	for name := range m.events {
		names = append(names, name)
	}
	return names
}

// GetActiveEvents returns currently active event names
func (m *Manager) GetActiveEvents() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]string, len(m.active))
	copy(result, m.active)
	return result
}

// GetRecentlyTriggeredEvents returns events triggered in the last specified duration
func (m *Manager) GetRecentlyTriggeredEvents(since time.Duration) []string {
	m.historyMu.RLock()
	defer m.historyMu.RUnlock()

	cutoff := time.Now().Add(-since)
	result := make([]string, 0)
	seen := make(map[string]bool)

	for _, entry := range m.history {
		if entry.Timestamp.After(cutoff) && entry.Triggered && !seen[entry.EventName] {
			result = append(result, entry.EventName)
			seen[entry.EventName] = true
		}
	}

	return result
}

// TriggerEvent manually triggers an event
func (m *Manager) TriggerEvent(name string, force bool, context map[string]interface{}) error {
	if !m.enabled {
		return fmt.Errorf("event manager is disabled")
	}

	request := &ProcessRequest{
		EventName: name,
		Force:     force,
		Context:   context,
		Source:    "manual",
		Response:  make(chan *Result, 1),
	}

	select {
	case m.processingQueue <- request:
		result := <-request.Response
		if !result.Success {
			return fmt.Errorf("event execution failed: %s", result.Error)
		}
		return nil
	default:
		return fmt.Errorf("event processing queue is full")
	}
}

// EvaluateEvent evaluates an event's trigger condition
func (m *Manager) EvaluateEvent(name string, context map[string]interface{}) (bool, error) {
	m.mu.RLock()
	event, exists := m.events[name]
	m.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("event %s not found", name)
	}

	if !event.Enabled {
		return false, nil
	}

	// Check cooldown
	if event.Cooldown != nil && time.Since(event.LastTriggered) < *event.Cooldown {
		return false, nil
	}

	// Check max triggers
	if event.MaxTriggers != nil && event.TriggerCount >= *event.MaxTriggers {
		return false, nil
	}

	// Evaluate trigger condition (this would integrate with CUE evaluation)
	// For now, return true as placeholder
	return true, nil
}

// ExecuteEvent executes an event's action
func (m *Manager) ExecuteEvent(name string, context map[string]interface{}) *Result {
	start := time.Now()

	m.mu.RLock()
	event, exists := m.events[name]
	m.mu.RUnlock()

	if !exists {
		return &Result{
			Success: false,
			Error:   fmt.Sprintf("event %s not found", name),
		}
	}

	// Execute the event action (this would integrate with CUE evaluation)
	// For now, just simulate execution
	result := &Result{
		Success:  true,
		Duration: time.Since(start),
		Output:   fmt.Sprintf("Executed event %s", name),
	}

	// Update event state
	m.mu.Lock()
	event.TriggerCount++
	event.LastTriggered = time.Now()
	event.CurrentlyActive = true
	m.mu.Unlock()

	// Add to history
	m.addToHistory(HistoryEntry{
		EventName: name,
		Type:      EventTypeCustom,
		Timestamp: time.Now(),
		Triggered: true,
		Source:    "execution",
		Duration:  result.Duration,
		Success:   result.Success,
		Error:     result.Error,
		Context:   context,
	})

	// Update active events
	m.updateActiveEvents(name, true)

	return result
}

// GetHistory returns event history
func (m *Manager) GetHistory(limit int) []HistoryEntry {
	m.historyMu.RLock()
	defer m.historyMu.RUnlock()

	if limit <= 0 || limit > len(m.history) {
		limit = len(m.history)
	}

	// Return most recent entries
	start := len(m.history) - limit
	result := make([]HistoryEntry, limit)
	copy(result, m.history[start:])
	return result
}

// ClearHistory clears event history
func (m *Manager) ClearHistory() {
	m.historyMu.Lock()
	defer m.historyMu.Unlock()
	m.history = make([]HistoryEntry, 0)
}

// SetEnabled enables or disables the event system
func (m *Manager) SetEnabled(enabled bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.enabled = enabled
}

// IsEnabled returns whether the event system is enabled
func (m *Manager) IsEnabled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.enabled
}

// processEvents processes events from the queue
func (m *Manager) processEvents() {
	for {
		select {
		case <-m.stopChan:
			return
		case request := <-m.processingQueue:
			result := m.ExecuteEvent(request.EventName, request.Context)
			request.Response <- result
		}
	}
}

// addToHistory adds an entry to the event history
func (m *Manager) addToHistory(entry HistoryEntry) {
	m.historyMu.Lock()
	defer m.historyMu.Unlock()

	m.history = append(m.history, entry)

	// Maintain history size limit
	if len(m.history) > m.maxHistory {
		m.history = m.history[1:]
	}
}

// updateActiveEvents updates the active events list
func (m *Manager) updateActiveEvents(eventName string, active bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if active {
		// Add to active list if not already present
		for _, name := range m.active {
			if name == eventName {
				return
			}
		}
		m.active = append(m.active, eventName)
	} else {
		// Remove from active list
		for i, name := range m.active {
			if name == eventName {
				m.active = append(m.active[:i], m.active[i+1:]...)
				break
			}
		}
	}
}

// LoadFromConfig loads events from configuration
func (m *Manager) LoadFromConfig(config *Config) error {
	if config == nil {
		return nil
	}

	m.config = config

	// Load custom events
	if config.Custom != nil {
		for name, customEvent := range config.Custom {
			event := &Event{
				Name:         name,
				Description:  customEvent.Description,
				Trigger:      customEvent.Trigger,
				Action:       customEvent.Action,
				Dependencies: customEvent.Dependencies,
				Enabled:      customEvent.Enabled,
				Priority:     customEvent.Priority,
				Metadata:     make(map[string]interface{}),
			}

			m.AddEvent(event)
		}
	}

	return nil
}

// GetConfig returns the current configuration
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.config
}

// GetStats returns event system statistics
func (m *Manager) GetStats() map[string]interface{} {
	m.mu.RLock()
	totalEvents := len(m.events)
	activeCount := len(m.active)
	m.mu.RUnlock()

	m.historyMu.RLock()
	historyCount := len(m.history)
	m.historyMu.RUnlock()

	return map[string]interface{}{
		"total_events":    totalEvents,
		"active_events":   activeCount,
		"history_entries": historyCount,
		"enabled":         m.enabled,
		"running":         m.running,
	}
}
