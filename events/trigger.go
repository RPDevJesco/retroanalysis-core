package events

import (
	"fmt"
	"sync"
	"time"
)

// TriggerType represents different types of triggers
type TriggerType string

const (
	TriggerTypeImmediate TriggerType = "immediate"
	TriggerTypeCondition TriggerType = "condition"
	TriggerTypeThreshold TriggerType = "threshold"
	TriggerTypeChange    TriggerType = "change"
	TriggerTypeSequence  TriggerType = "sequence"
	TriggerTypeTimer     TriggerType = "timer"
)

// Trigger represents an event trigger with evaluation logic
type Trigger struct {
	Type         TriggerType            `json:"type"`
	Expression   string                 `json:"expression"`   // CUE expression
	Dependencies []string               `json:"dependencies"` // Properties to watch
	Threshold    *ThresholdConfig       `json:"threshold,omitempty"`
	Timer        *TimerConfig           `json:"timer,omitempty"`
	Sequence     *SequenceConfig        `json:"sequence,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ThresholdConfig represents threshold-based trigger configuration
type ThresholdConfig struct {
	Property   string      `json:"property"`
	Operator   string      `json:"operator"` // "gt", "lt", "eq", "ne", "gte", "lte"
	Value      interface{} `json:"value"`
	Hysteresis *float64    `json:"hysteresis,omitempty"` // Prevents oscillation
}

// TimerConfig represents timer-based trigger configuration
type TimerConfig struct {
	Interval time.Duration `json:"interval"`
	Repeat   bool          `json:"repeat"`
	MaxCount *uint         `json:"max_count,omitempty"`
}

// SequenceConfig represents sequence-based trigger configuration
type SequenceConfig struct {
	Steps   []SequenceStep `json:"steps"`
	Timeout time.Duration  `json:"timeout,omitempty"`
	Reset   bool           `json:"reset"` // Reset sequence on timeout
}

// SequenceStep represents a step in a sequence trigger
type SequenceStep struct {
	Condition string        `json:"condition"`
	Timeout   time.Duration `json:"timeout,omitempty"`
	Optional  bool          `json:"optional"`
}

// TriggerEngine handles trigger evaluation and management
type TriggerEngine struct {
	triggers        map[string]*Trigger
	sequences       map[string]*SequenceState
	timers          map[string]*TimerState
	lastValues      map[string]interface{}
	propertyWatches map[string][]string // property -> list of trigger names
	mu              sync.RWMutex

	// Evaluation context
	evaluator TriggerEvaluator
}

// SequenceState tracks the state of a sequence trigger
type SequenceState struct {
	CurrentStep int
	StartTime   time.Time
	StepTimes   []time.Time
	Completed   bool
	Failed      bool
}

// TimerState tracks the state of a timer trigger
type TimerState struct {
	StartTime    time.Time
	LastTrigger  time.Time
	TriggerCount uint
	Active       bool
}

// TriggerEvaluator interface for evaluating trigger expressions
type TriggerEvaluator interface {
	Evaluate(expression string, context map[string]interface{}) (bool, error)
}

// NewTriggerEngine creates a new trigger engine
func NewTriggerEngine(evaluator TriggerEvaluator) *TriggerEngine {
	return &TriggerEngine{
		triggers:        make(map[string]*Trigger),
		sequences:       make(map[string]*SequenceState),
		timers:          make(map[string]*TimerState),
		lastValues:      make(map[string]interface{}),
		propertyWatches: make(map[string][]string),
		evaluator:       evaluator,
	}
}

// AddTrigger adds a trigger to the engine
func (te *TriggerEngine) AddTrigger(name string, trigger *Trigger) {
	te.mu.Lock()
	defer te.mu.Unlock()

	te.triggers[name] = trigger

	// Set up property watches
	for _, property := range trigger.Dependencies {
		if te.propertyWatches[property] == nil {
			te.propertyWatches[property] = make([]string, 0)
		}
		te.propertyWatches[property] = append(te.propertyWatches[property], name)
	}

	// Initialize timer if needed
	if trigger.Timer != nil {
		te.timers[name] = &TimerState{
			StartTime: time.Now(),
			Active:    true,
		}
	}

	// Initialize sequence if needed
	if trigger.Sequence != nil {
		te.sequences[name] = &SequenceState{
			CurrentStep: 0,
			StartTime:   time.Now(),
			StepTimes:   make([]time.Time, len(trigger.Sequence.Steps)),
		}
	}
}

// RemoveTrigger removes a trigger from the engine
func (te *TriggerEngine) RemoveTrigger(name string) {
	te.mu.Lock()
	defer te.mu.Unlock()

	trigger, exists := te.triggers[name]
	if !exists {
		return
	}

	// Remove property watches
	for _, property := range trigger.Dependencies {
		watches := te.propertyWatches[property]
		for i, watchName := range watches {
			if watchName == name {
				te.propertyWatches[property] = append(watches[:i], watches[i+1:]...)
				break
			}
		}
	}

	delete(te.triggers, name)
	delete(te.sequences, name)
	delete(te.timers, name)
}

// EvaluateTrigger evaluates a specific trigger
func (te *TriggerEngine) EvaluateTrigger(name string, context map[string]interface{}) (bool, error) {
	te.mu.RLock()
	trigger, exists := te.triggers[name]
	te.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("trigger %s not found", name)
	}

	switch trigger.Type {
	case TriggerTypeImmediate:
		return true, nil
	case TriggerTypeCondition:
		return te.evaluateConditionTrigger(trigger, context)
	case TriggerTypeThreshold:
		return te.evaluateThresholdTrigger(name, trigger, context)
	case TriggerTypeChange:
		return te.evaluateChangeTrigger(trigger, context)
	case TriggerTypeSequence:
		return te.evaluateSequenceTrigger(name, trigger, context)
	case TriggerTypeTimer:
		return te.evaluateTimerTrigger(name, trigger)
	default:
		return false, fmt.Errorf("unsupported trigger type: %s", trigger.Type)
	}
}

// EvaluatePropertyChange evaluates triggers when a property changes
func (te *TriggerEngine) EvaluatePropertyChange(property string, oldValue, newValue interface{}) []string {
	te.mu.RLock()
	triggerNames := te.propertyWatches[property]
	te.mu.RUnlock()

	if len(triggerNames) == 0 {
		return nil
	}

	// Update last value
	te.mu.Lock()
	te.lastValues[property] = newValue
	te.mu.Unlock()

	triggeredEvents := make([]string, 0)
	context := map[string]interface{}{
		property:           newValue,
		property + "_old":  oldValue,
		"changed_property": property,
	}

	for _, triggerName := range triggerNames {
		if triggered, err := te.EvaluateTrigger(triggerName, context); err == nil && triggered {
			triggeredEvents = append(triggeredEvents, triggerName)
		}
	}

	return triggeredEvents
}

// evaluateConditionTrigger evaluates a condition-based trigger
func (te *TriggerEngine) evaluateConditionTrigger(trigger *Trigger, context map[string]interface{}) (bool, error) {
	if te.evaluator == nil {
		return false, fmt.Errorf("no trigger evaluator available")
	}

	return te.evaluator.Evaluate(trigger.Expression, context)
}

// evaluateThresholdTrigger evaluates a threshold-based trigger
func (te *TriggerEngine) evaluateThresholdTrigger(name string, trigger *Trigger, context map[string]interface{}) (bool, error) {
	if trigger.Threshold == nil {
		return false, fmt.Errorf("threshold trigger missing threshold config")
	}

	propertyValue, exists := context[trigger.Threshold.Property]
	if !exists {
		return false, nil
	}

	// Convert values to float64 for comparison
	currentValue, ok := te.toFloat64(propertyValue)
	if !ok {
		return false, fmt.Errorf("cannot convert property value to numeric")
	}

	thresholdValue, ok := te.toFloat64(trigger.Threshold.Value)
	if !ok {
		return false, fmt.Errorf("cannot convert threshold value to numeric")
	}

	// Apply hysteresis if configured
	if trigger.Threshold.Hysteresis != nil {
		// TODO: Implement hysteresis logic
	}

	// Evaluate threshold condition
	switch trigger.Threshold.Operator {
	case "gt":
		return currentValue > thresholdValue, nil
	case "lt":
		return currentValue < thresholdValue, nil
	case "eq":
		return currentValue == thresholdValue, nil
	case "ne":
		return currentValue != thresholdValue, nil
	case "gte":
		return currentValue >= thresholdValue, nil
	case "lte":
		return currentValue <= thresholdValue, nil
	default:
		return false, fmt.Errorf("unsupported threshold operator: %s", trigger.Threshold.Operator)
	}
}

// evaluateChangeTrigger evaluates a change-based trigger
func (te *TriggerEngine) evaluateChangeTrigger(trigger *Trigger, context map[string]interface{}) (bool, error) {
	// For change triggers, we evaluate if any dependency actually changed
	for _, property := range trigger.Dependencies {
		if oldValue, exists := context[property+"_old"]; exists {
			currentValue := context[property]
			if !te.valuesEqual(oldValue, currentValue) {
				return true, nil
			}
		}
	}
	return false, nil
}

// evaluateSequenceTrigger evaluates a sequence-based trigger
func (te *TriggerEngine) evaluateSequenceTrigger(name string, trigger *Trigger, context map[string]interface{}) (bool, error) {
	te.mu.Lock()
	defer te.mu.Unlock()

	state, exists := te.sequences[name]
	if !exists {
		return false, fmt.Errorf("sequence state not found for trigger %s", name)
	}

	if state.Completed || state.Failed {
		return state.Completed, nil
	}

	// Check timeout
	if trigger.Sequence.Timeout > 0 && time.Since(state.StartTime) > trigger.Sequence.Timeout {
		state.Failed = true
		if trigger.Sequence.Reset {
			state.CurrentStep = 0
			state.StartTime = time.Now()
			state.Failed = false
		}
		return false, nil
	}

	// Evaluate current step
	if state.CurrentStep < len(trigger.Sequence.Steps) {
		step := trigger.Sequence.Steps[state.CurrentStep]

		if triggered, err := te.evaluator.Evaluate(step.Condition, context); err == nil && triggered {
			state.StepTimes[state.CurrentStep] = time.Now()
			state.CurrentStep++

			// Check if sequence is complete
			if state.CurrentStep >= len(trigger.Sequence.Steps) {
				state.Completed = true
				return true, nil
			}
		}
	}

	return false, nil
}

// evaluateTimerTrigger evaluates a timer-based trigger
func (te *TriggerEngine) evaluateTimerTrigger(name string, trigger *Trigger) (bool, error) {
	te.mu.Lock()
	defer te.mu.Unlock()

	state, exists := te.timers[name]
	if !exists {
		return false, fmt.Errorf("timer state not found for trigger %s", name)
	}

	if !state.Active {
		return false, nil
	}

	now := time.Now()

	// Check if it's time to trigger
	if state.LastTrigger.IsZero() {
		// First trigger
		if now.Sub(state.StartTime) >= trigger.Timer.Interval {
			state.LastTrigger = now
			state.TriggerCount++

			// Check if we should stop repeating
			if !trigger.Timer.Repeat ||
				(trigger.Timer.MaxCount != nil && state.TriggerCount >= *trigger.Timer.MaxCount) {
				state.Active = false
			}

			return true, nil
		}
	} else {
		// Subsequent triggers
		if now.Sub(state.LastTrigger) >= trigger.Timer.Interval {
			state.LastTrigger = now
			state.TriggerCount++

			// Check if we should stop repeating
			if !trigger.Timer.Repeat ||
				(trigger.Timer.MaxCount != nil && state.TriggerCount >= *trigger.Timer.MaxCount) {
				state.Active = false
			}

			return true, nil
		}
	}

	return false, nil
}

// Helper methods

// toFloat64 converts various numeric types to float64
func (te *TriggerEngine) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// valuesEqual compares two values for equality
func (te *TriggerEngine) valuesEqual(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// GetTriggerStats returns statistics about triggers
func (te *TriggerEngine) GetTriggerStats() map[string]interface{} {
	te.mu.RLock()
	defer te.mu.RUnlock()

	activeTimers := 0
	activeSequences := 0

	for _, state := range te.timers {
		if state.Active {
			activeTimers++
		}
	}

	for _, state := range te.sequences {
		if !state.Completed && !state.Failed {
			activeSequences++
		}
	}

	return map[string]interface{}{
		"total_triggers":     len(te.triggers),
		"active_timers":      activeTimers,
		"active_sequences":   activeSequences,
		"watched_properties": len(te.propertyWatches),
	}
}
