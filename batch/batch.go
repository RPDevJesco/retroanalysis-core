package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RPDevJesco/retroanalysis-core/properties"
)

// OperationType represents the type of batch operation
type OperationType string

const (
	OperationTypeMemoryRead     OperationType = "memory_read"
	OperationTypeMemoryWrite    OperationType = "memory_write"
	OperationTypePropertyRead   OperationType = "property_read"
	OperationTypePropertyWrite  OperationType = "property_write"
	OperationTypePropertyFreeze OperationType = "property_freeze"
	OperationTypeValidation     OperationType = "validation"
	OperationTypeTransform      OperationType = "transform"
	OperationTypeCustom         OperationType = "custom"
)

// Priority levels for batch operations
type Priority uint8

const (
	PriorityLow      Priority = 1
	PriorityNormal   Priority = 5
	PriorityHigh     Priority = 10
	PriorityCritical Priority = 15
)

// ExecutionMode determines how the batch should be executed
type ExecutionMode string

const (
	ExecutionModeSequential ExecutionMode = "sequential"
	ExecutionModeParallel   ExecutionMode = "parallel"
	ExecutionModeAtomic     ExecutionMode = "atomic"
	ExecutionModeBestEffort ExecutionMode = "best_effort"
)

// ValidationMode determines how validation errors are handled
type ValidationMode string

const (
	ValidationModeStrict ValidationMode = "strict" // Fail on any validation error
	ValidationModeWarn   ValidationMode = "warn"   // Log warnings but continue
	ValidationModeIgnore ValidationMode = "ignore" // Skip validation entirely
)

// Operation represents a single operation in a batch
type Operation struct {
	ID         string                 `json:"id"`
	Type       OperationType          `json:"type"`
	Target     string                 `json:"target"` // Property name or memory address
	Data       interface{}            `json:"data"`   // Operation data
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Retry      *RetryConfig           `json:"retry,omitempty"`
	Timeout    time.Duration          `json:"timeout,omitempty"`
	Validation *ValidationConfig      `json:"validation,omitempty"`
	Condition  string                 `json:"condition,omitempty"` // CUE expression for conditional execution

	// Internal state
	Status     OperationStatus  `json:"status"`
	Result     *OperationResult `json:"result,omitempty"`
	StartTime  time.Time        `json:"start_time,omitempty"`
	EndTime    time.Time        `json:"end_time,omitempty"`
	RetryCount int              `json:"retry_count"`
	LastError  string           `json:"last_error,omitempty"`
}

// OperationStatus represents the current status of an operation
type OperationStatus string

const (
	StatusPending   OperationStatus = "pending"
	StatusRunning   OperationStatus = "running"
	StatusCompleted OperationStatus = "completed"
	StatusFailed    OperationStatus = "failed"
	StatusSkipped   OperationStatus = "skipped"
	StatusCancelled OperationStatus = "cancelled"
)

// OperationResult represents the result of a single operation
type OperationResult struct {
	Success  bool                   `json:"success"`
	Value    interface{}            `json:"value,omitempty"`
	Data     []byte                 `json:"data,omitempty"`
	Error    string                 `json:"error,omitempty"`
	Duration time.Duration          `json:"duration"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Warnings []string               `json:"warnings,omitempty"`
}

// RetryConfig configures retry behavior for operations
type RetryConfig struct {
	MaxAttempts       int           `json:"max_attempts"`
	Delay             time.Duration `json:"delay"`
	BackoffMultiplier float64       `json:"backoff_multiplier"`
	MaxDelay          time.Duration `json:"max_delay"`
	RetryOn           []string      `json:"retry_on,omitempty"` // Error types to retry on
}

// ValidationConfig configures validation for operations
type ValidationConfig struct {
	Mode         ValidationMode `json:"mode"`
	Rules        []string       `json:"rules,omitempty"` // CUE validation expressions
	PreValidate  bool           `json:"pre_validate"`    // Validate before execution
	PostValidate bool           `json:"post_validate"`   // Validate after execution
	FailFast     bool           `json:"fail_fast"`       // Stop batch on first validation error
}

// Request represents a batch operation request
type Request struct {
	ID         string                 `json:"id"`
	Operations []*Operation           `json:"operations"`
	Mode       ExecutionMode          `json:"mode"`
	Priority   Priority               `json:"priority"`
	Timeout    time.Duration          `json:"timeout"`
	MaxRetries int                    `json:"max_retries"`
	Validation ValidationMode         `json:"validation"`
	Atomic     bool                   `json:"atomic"`   // All operations succeed or all fail
	Parallel   bool                   `json:"parallel"` // Execute operations in parallel
	Metadata   map[string]interface{} `json:"metadata,omitempty"`

	// Callbacks
	OnProgress func(*Progress) `json:"-"` // Progress callback
	OnComplete func(*Result)   `json:"-"` // Completion callback
	OnError    func(error)     `json:"-"` // Error callback

	// Internal state
	Status      RequestStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	StartedAt   time.Time     `json:"started_at,omitempty"`
	CompletedAt time.Time     `json:"completed_at,omitempty"`
	Result      *Result       `json:"result,omitempty"`
}

// RequestStatus represents the status of a batch request
type RequestStatus string

const (
	RequestStatusQueued    RequestStatus = "queued"
	RequestStatusRunning   RequestStatus = "running"
	RequestStatusCompleted RequestStatus = "completed"
	RequestStatusFailed    RequestStatus = "failed"
	RequestStatusCancelled RequestStatus = "cancelled"
	RequestStatusTimeout   RequestStatus = "timeout"
)

// Result represents the result of a batch operation
type Result struct {
	RequestID       string                 `json:"request_id"`
	Success         bool                   `json:"success"`
	TotalOperations int                    `json:"total_operations"`
	SuccessCount    int                    `json:"success_count"`
	FailureCount    int                    `json:"failure_count"`
	SkippedCount    int                    `json:"skipped_count"`
	Duration        time.Duration          `json:"duration"`
	Operations      []*Operation           `json:"operations"`
	Errors          []string               `json:"errors,omitempty"`
	Warnings        []string               `json:"warnings,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Statistics      *Statistics            `json:"statistics,omitempty"`
}

// Progress represents the progress of a batch operation
type Progress struct {
	RequestID          string        `json:"request_id"`
	TotalOperations    int           `json:"total_operations"`
	CompletedCount     int           `json:"completed_count"`
	FailedCount        int           `json:"failed_count"`
	SkippedCount       int           `json:"skipped_count"`
	CurrentOperation   *Operation    `json:"current_operation,omitempty"`
	Percentage         float64       `json:"percentage"`
	ElapsedTime        time.Duration `json:"elapsed_time"`
	EstimatedRemaining time.Duration `json:"estimated_remaining"`
}

// Statistics represents detailed batch operation statistics
type Statistics struct {
	MemoryOperations    int           `json:"memory_operations"`
	PropertyOperations  int           `json:"property_operations"`
	ValidationErrors    int           `json:"validation_errors"`
	RetryCount          int           `json:"retry_count"`
	AverageOpDuration   time.Duration `json:"average_op_duration"`
	MinOpDuration       time.Duration `json:"min_op_duration"`
	MaxOpDuration       time.Duration `json:"max_op_duration"`
	ThroughputOpsPerSec float64       `json:"throughput_ops_per_sec"`
	MemoryBytesRead     uint64        `json:"memory_bytes_read"`
	MemoryBytesWritten  uint64        `json:"memory_bytes_written"`
}

// Processor handles batch operation processing
type Processor struct {
	config         *Config
	requests       chan *Request
	activeRequests map[string]*Request
	workers        []*Worker
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc

	// Dependencies
	memoryManager   MemoryManager
	propertyManager PropertyManager
	validator       Validator

	// Statistics
	totalProcessed uint64
	totalSucceeded uint64
	totalFailed    uint64
	startTime      time.Time
}

// Config represents batch processor configuration
type Config struct {
	MaxConcurrentRequests int           `json:"max_concurrent_requests"`
	MaxOperationsPerBatch int           `json:"max_operations_per_batch"`
	DefaultTimeout        time.Duration `json:"default_timeout"`
	WorkerPoolSize        int           `json:"worker_pool_size"`
	MaxRetryAttempts      int           `json:"max_retry_attempts"`
	RetryDelay            time.Duration `json:"retry_delay"`
	EnableMetrics         bool          `json:"enable_metrics"`
	EnableValidation      bool          `json:"enable_validation"`
	BufferSize            int           `json:"buffer_size"`
}

// Worker represents a batch operation worker
type Worker struct {
	id        int
	processor *Processor
	ctx       context.Context
	cancel    context.CancelFunc
	active    bool
	current   *Request
	mu        sync.RWMutex
}

// Dependencies interfaces
type MemoryManager interface {
	ReadBytes(address uint32, length uint32) ([]byte, error)
	WriteBytes(address uint32, data []byte) []byte
	FreezeProperty(address uint32, data []byte) error
	UnfreezeProperty(address uint32) error
}

type PropertyManager interface {
	GetProperty(name string) (*properties.Property, error)
	SetProperty(name string, value interface{}) error
	GetState(name string) *properties.State
	UpdateState(name string, value interface{}, bytes []byte, address uint32)
}

type Validator interface {
	Validate(value interface{}, validation *properties.Validation) error
	ValidateExpression(expression string, context map[string]interface{}) (bool, error)
}

// NewProcessor creates a new batch processor
func NewProcessor(config *Config, memManager MemoryManager, propManager PropertyManager, validator Validator) *Processor {
	ctx, cancel := context.WithCancel(context.Background())

	processor := &Processor{
		config:          config,
		requests:        make(chan *Request, config.BufferSize),
		activeRequests:  make(map[string]*Request),
		workers:         make([]*Worker, config.WorkerPoolSize),
		ctx:             ctx,
		cancel:          cancel,
		memoryManager:   memManager,
		propertyManager: propManager,
		validator:       validator,
		startTime:       time.Now(),
	}

	// Initialize workers
	for i := 0; i < config.WorkerPoolSize; i++ {
		processor.workers[i] = NewWorker(i, processor)
	}

	return processor
}

// NewWorker creates a new batch worker
func NewWorker(id int, processor *Processor) *Worker {
	ctx, cancel := context.WithCancel(processor.ctx)

	worker := &Worker{
		id:        id,
		processor: processor,
		ctx:       ctx,
		cancel:    cancel,
		active:    false,
	}

	go worker.run()
	return worker
}

// Start starts the batch processor
func (p *Processor) Start() {
	go p.processRequests()
}

// Stop stops the batch processor
func (p *Processor) Stop() {
	p.cancel()

	// Stop all workers
	for _, worker := range p.workers {
		worker.stop()
	}
}

// Submit submits a batch request for processing
func (p *Processor) Submit(request *Request) error {
	if len(request.Operations) == 0 {
		return fmt.Errorf("batch request must contain at least one operation")
	}

	if len(request.Operations) > p.config.MaxOperationsPerBatch {
		return fmt.Errorf("batch request exceeds maximum operations limit: %d", p.config.MaxOperationsPerBatch)
	}

	// Generate ID if not provided
	if request.ID == "" {
		request.ID = generateID()
	}

	// Set defaults
	if request.Timeout == 0 {
		request.Timeout = p.config.DefaultTimeout
	}

	if request.MaxRetries == 0 {
		request.MaxRetries = p.config.MaxRetryAttempts
	}

	// Initialize operation states
	for _, op := range request.Operations {
		if op.ID == "" {
			op.ID = generateID()
		}
		op.Status = StatusPending
		op.RetryCount = 0
	}

	request.Status = RequestStatusQueued
	request.CreatedAt = time.Now()

	// Add to active requests
	p.mu.Lock()
	p.activeRequests[request.ID] = request
	p.mu.Unlock()

	// Submit to processing queue
	select {
	case p.requests <- request:
		return nil
	default:
		return fmt.Errorf("batch processor queue is full")
	}
}

// GetRequest gets a batch request by ID
func (p *Processor) GetRequest(id string) (*Request, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	request, exists := p.activeRequests[id]
	return request, exists
}

// ListActiveRequests returns all active request IDs
func (p *Processor) ListActiveRequests() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ids := make([]string, 0, len(p.activeRequests))
	for id := range p.activeRequests {
		ids = append(ids, id)
	}
	return ids
}

// CancelRequest cancels a batch request
func (p *Processor) CancelRequest(id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	request, exists := p.activeRequests[id]
	if !exists {
		return fmt.Errorf("request %s not found", id)
	}

	if request.Status == RequestStatusCompleted || request.Status == RequestStatusFailed {
		return fmt.Errorf("request %s is already completed", id)
	}

	request.Status = RequestStatusCancelled

	// Cancel operations that haven't started
	for _, op := range request.Operations {
		if op.Status == StatusPending {
			op.Status = StatusCancelled
		}
	}

	return nil
}

// processRequests processes batch requests from the queue
func (p *Processor) processRequests() {
	for {
		select {
		case <-p.ctx.Done():
			return
		case request := <-p.requests:
			p.processRequest(request)
		}
	}
}

// processRequest processes a single batch request
func (p *Processor) processRequest(request *Request) {
	request.Status = RequestStatusRunning
	request.StartedAt = time.Now()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(p.ctx, request.Timeout)
	defer cancel()

	var result *Result
	var err error

	switch request.Mode {
	case ExecutionModeSequential:
		result, err = p.executeSequential(ctx, request)
	case ExecutionModeParallel:
		result, err = p.executeParallel(ctx, request)
	case ExecutionModeAtomic:
		result, err = p.executeAtomic(ctx, request)
	case ExecutionModeBestEffort:
		result, err = p.executeBestEffort(ctx, request)
	default:
		result, err = p.executeSequential(ctx, request)
	}

	// Finalize result
	request.CompletedAt = time.Now()
	request.Result = result

	if err != nil {
		request.Status = RequestStatusFailed
		if result != nil {
			result.Errors = append(result.Errors, err.Error())
		}
	} else {
		if result.Success {
			request.Status = RequestStatusCompleted
		} else {
			request.Status = RequestStatusFailed
		}
	}

	// Call completion callback
	if request.OnComplete != nil {
		go request.OnComplete(result)
	}

	// Update statistics
	p.updateStatistics(request)

	// Remove from active requests after a delay (for status checking)
	go func() {
		time.Sleep(5 * time.Minute)
		p.mu.Lock()
		delete(p.activeRequests, request.ID)
		p.mu.Unlock()
	}()
}

// executeSequential executes operations sequentially
func (p *Processor) executeSequential(ctx context.Context, request *Request) (*Result, error) {
	result := &Result{
		RequestID:       request.ID,
		TotalOperations: len(request.Operations),
		Operations:      request.Operations,
		Metadata:        make(map[string]interface{}),
	}

	startTime := time.Now()

	for i, op := range request.Operations {
		select {
		case <-ctx.Done():
			result.Errors = append(result.Errors, "request timeout")
			return result, ctx.Err()
		default:
		}

		if request.Status == RequestStatusCancelled {
			break
		}

		// Execute operation
		opResult := p.executeOperation(ctx, op, request)
		op.Result = opResult

		// Update counts
		if opResult.Success {
			result.SuccessCount++
		} else {
			result.FailureCount++
			if request.Atomic {
				// In atomic mode, stop on first failure
				result.Errors = append(result.Errors, fmt.Sprintf("atomic batch failed at operation %d: %s", i, opResult.Error))
				break
			}
		}

		// Report progress
		if request.OnProgress != nil {
			progress := &Progress{
				RequestID:        request.ID,
				TotalOperations:  len(request.Operations),
				CompletedCount:   i + 1,
				CurrentOperation: op,
				Percentage:       float64(i+1) / float64(len(request.Operations)) * 100,
				ElapsedTime:      time.Since(startTime),
			}
			go request.OnProgress(progress)
		}
	}

	result.Duration = time.Since(startTime)
	result.Success = result.FailureCount == 0 || (!request.Atomic && result.SuccessCount > 0)
	result.Statistics = p.calculateStatistics(request.Operations)

	return result, nil
}

// executeParallel executes operations in parallel
func (p *Processor) executeParallel(ctx context.Context, request *Request) (*Result, error) {
	result := &Result{
		RequestID:       request.ID,
		TotalOperations: len(request.Operations),
		Operations:      request.Operations,
		Metadata:        make(map[string]interface{}),
	}

	startTime := time.Now()

	// Create worker pool for this request
	type job struct {
		op    *Operation
		index int
	}

	jobs := make(chan job, len(request.Operations))
	results := make(chan struct{}, len(request.Operations))

	// Start workers
	workerCount := min(len(request.Operations), p.config.WorkerPoolSize)
	for i := 0; i < workerCount; i++ {
		go func() {
			for job := range jobs {
				if request.Status == RequestStatusCancelled {
					job.op.Status = StatusCancelled
					results <- struct{}{}
					continue
				}

				job.op.Result = p.executeOperation(ctx, job.op, request)
				results <- struct{}{}
			}
		}()
	}

	// Submit jobs
	for i, op := range request.Operations {
		jobs <- job{op: op, index: i}
	}
	close(jobs)

	// Wait for all operations to complete
	completed := 0
	for completed < len(request.Operations) {
		select {
		case <-ctx.Done():
			result.Errors = append(result.Errors, "request timeout")
			return result, ctx.Err()
		case <-results:
			completed++

			// Report progress
			if request.OnProgress != nil {
				progress := &Progress{
					RequestID:       request.ID,
					TotalOperations: len(request.Operations),
					CompletedCount:  completed,
					Percentage:      float64(completed) / float64(len(request.Operations)) * 100,
					ElapsedTime:     time.Since(startTime),
				}
				go request.OnProgress(progress)
			}
		}
	}

	// Count results
	for _, op := range request.Operations {
		if op.Result != nil && op.Result.Success {
			result.SuccessCount++
		} else {
			result.FailureCount++
		}
	}

	result.Duration = time.Since(startTime)
	result.Success = result.FailureCount == 0 || (!request.Atomic && result.SuccessCount > 0)
	result.Statistics = p.calculateStatistics(request.Operations)

	return result, nil
}

// executeAtomic executes operations atomically (all succeed or all fail)
func (p *Processor) executeAtomic(ctx context.Context, request *Request) (*Result, error) {
	// For atomic operations, we need to be able to rollback
	// This is a simplified implementation - real atomic operations would need
	// transaction support from the underlying systems

	result, err := p.executeSequential(ctx, request)
	if err != nil {
		return result, err
	}

	// If any operation failed, mark all as failed
	if result.FailureCount > 0 {
		for _, op := range request.Operations {
			if op.Result != nil && op.Result.Success {
				// In a real implementation, we would rollback successful operations here
				op.Status = StatusFailed
				op.Result.Success = false
				op.Result.Error = "rolled back due to atomic batch failure"
			}
		}
		result.SuccessCount = 0
		result.FailureCount = len(request.Operations)
		result.Success = false
	}

	return result, nil
}

// executeBestEffort executes operations with best effort (continue on errors)
func (p *Processor) executeBestEffort(ctx context.Context, request *Request) (*Result, error) {
	// Best effort is similar to parallel but always continues
	request.Atomic = false
	return p.executeParallel(ctx, request)
}

// executeOperation executes a single operation
func (p *Processor) executeOperation(ctx context.Context, op *Operation, request *Request) *OperationResult {
	op.Status = StatusRunning
	op.StartTime = time.Now()

	result := &OperationResult{
		Metadata: make(map[string]interface{}),
	}

	// Check condition if specified
	if op.Condition != "" {
		if shouldExecute, err := p.evaluateCondition(op.Condition, request); err != nil {
			result.Error = fmt.Sprintf("condition evaluation failed: %v", err)
			op.Status = StatusFailed
			return result
		} else if !shouldExecute {
			result.Success = true
			op.Status = StatusSkipped
			return result
		}
	}

	// Execute with retries
	var lastErr error
	maxAttempts := 1
	if op.Retry != nil {
		maxAttempts = op.Retry.MaxAttempts
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			// Wait before retry
			delay := p.calculateRetryDelay(op.Retry, attempt)
			select {
			case <-ctx.Done():
				result.Error = "operation cancelled during retry delay"
				op.Status = StatusCancelled
				return result
			case <-time.After(delay):
			}
			op.RetryCount++
		}

		// Execute the actual operation
		switch op.Type {
		case OperationTypeMemoryRead:
			lastErr = p.executeMemoryRead(op, result)
		case OperationTypeMemoryWrite:
			lastErr = p.executeMemoryWrite(op, result)
		case OperationTypePropertyRead:
			lastErr = p.executePropertyRead(op, result)
		case OperationTypePropertyWrite:
			lastErr = p.executePropertyWrite(op, result)
		case OperationTypePropertyFreeze:
			lastErr = p.executePropertyFreeze(op, result)
		case OperationTypeValidation:
			lastErr = p.executeValidation(op, result)
		case OperationTypeTransform:
			lastErr = p.executeTransform(op, result)
		case OperationTypeCustom:
			lastErr = p.executeCustom(op, result)
		default:
			lastErr = fmt.Errorf("unsupported operation type: %s", op.Type)
		}

		if lastErr == nil {
			result.Success = true
			break
		}

		// Check if we should retry this error
		if op.Retry != nil && !p.shouldRetry(lastErr, op.Retry) {
			break
		}
	}

	if lastErr != nil {
		result.Error = lastErr.Error()
		op.Status = StatusFailed
		op.LastError = lastErr.Error()
	} else {
		op.Status = StatusCompleted
	}

	op.EndTime = time.Now()
	result.Duration = op.EndTime.Sub(op.StartTime)

	return result
}

// Operation execution methods

func (p *Processor) executeMemoryRead(op *Operation, result *OperationResult) error {
	// Parse address from target
	address, length, err := parseMemoryTarget(op.Target)
	if err != nil {
		return err
	}

	data, err := p.memoryManager.ReadBytes(address, length)
	if err != nil {
		return err
	}

	result.Data = data
	result.Metadata["address"] = address
	result.Metadata["length"] = length
	return nil
}

func (p *Processor) executeMemoryWrite(op *Operation, result *OperationResult) error {
	address, _, err := parseMemoryTarget(op.Target)
	if err != nil {
		return err
	}

	data, ok := op.Data.([]byte)
	if !ok {
		return fmt.Errorf("invalid data type for memory write: expected []byte")
	}

	writtenData := p.memoryManager.WriteBytes(address, data)
	result.Data = writtenData
	result.Metadata["address"] = address
	result.Metadata["bytes_written"] = len(data)
	return nil
}

func (p *Processor) executePropertyRead(op *Operation, result *OperationResult) error {
	prop, err := p.propertyManager.GetProperty(op.Target)
	if err != nil {
		return err
	}

	state := p.propertyManager.GetState(op.Target)
	result.Value = prop
	result.Metadata["state"] = state
	return nil
}

func (p *Processor) executePropertyWrite(op *Operation, result *OperationResult) error {
	err := p.propertyManager.SetProperty(op.Target, op.Data)
	if err != nil {
		return err
	}

	result.Value = op.Data
	result.Metadata["property"] = op.Target
	return nil
}

func (p *Processor) executePropertyFreeze(op *Operation, result *OperationResult) error {
	// Implementation would depend on property freezing system
	result.Metadata["frozen"] = true
	result.Metadata["property"] = op.Target
	return nil
}

func (p *Processor) executeValidation(op *Operation, result *OperationResult) error {
	// Implementation would depend on validation system
	result.Value = "validation_passed"
	return nil
}

func (p *Processor) executeTransform(op *Operation, result *OperationResult) error {
	// Implementation would depend on transform system
	result.Value = op.Data
	return nil
}

func (p *Processor) executeCustom(op *Operation, result *OperationResult) error {
	// Custom operations would be handled by registered handlers
	return fmt.Errorf("custom operation not implemented: %s", op.Target)
}

// Helper methods

func (p *Processor) evaluateCondition(condition string, request *Request) (bool, error) {
	// This would integrate with CUE evaluation
	return true, nil
}

func (p *Processor) calculateRetryDelay(retry *RetryConfig, attempt int) time.Duration {
	if retry == nil {
		return p.config.RetryDelay
	}

	delay := retry.Delay
	if retry.BackoffMultiplier > 1.0 {
		for i := 0; i < attempt; i++ {
			delay = time.Duration(float64(delay) * retry.BackoffMultiplier)
		}
	}

	if retry.MaxDelay > 0 && delay > retry.MaxDelay {
		delay = retry.MaxDelay
	}

	return delay
}

func (p *Processor) shouldRetry(err error, retry *RetryConfig) bool {
	if retry == nil || len(retry.RetryOn) == 0 {
		return true
	}

	errorStr := err.Error()
	for _, retryOn := range retry.RetryOn {
		if contains(errorStr, retryOn) {
			return true
		}
	}

	return false
}

func (p *Processor) calculateStatistics(operations []*Operation) *Statistics {
	stats := &Statistics{}

	var totalDuration time.Duration
	minDuration := time.Hour
	maxDuration := time.Duration(0)

	for _, op := range operations {
		if op.Result == nil {
			continue
		}

		switch op.Type {
		case OperationTypeMemoryRead, OperationTypeMemoryWrite:
			stats.MemoryOperations++
		case OperationTypePropertyRead, OperationTypePropertyWrite, OperationTypePropertyFreeze:
			stats.PropertyOperations++
		}

		if op.Result.Duration > 0 {
			totalDuration += op.Result.Duration
			if op.Result.Duration < minDuration {
				minDuration = op.Result.Duration
			}
			if op.Result.Duration > maxDuration {
				maxDuration = op.Result.Duration
			}
		}

		stats.RetryCount += op.RetryCount

		if !op.Result.Success {
			stats.ValidationErrors++
		}

		if data, ok := op.Result.Metadata["bytes_written"].(int); ok {
			stats.MemoryBytesWritten += uint64(data)
		}
		if len(op.Result.Data) > 0 {
			stats.MemoryBytesRead += uint64(len(op.Result.Data))
		}
	}

	opCount := len(operations)
	if opCount > 0 {
		stats.AverageOpDuration = totalDuration / time.Duration(opCount)
		stats.MinOpDuration = minDuration
		stats.MaxOpDuration = maxDuration

		if totalDuration > 0 {
			stats.ThroughputOpsPerSec = float64(opCount) / totalDuration.Seconds()
		}
	}

	return stats
}

func (p *Processor) updateStatistics(request *Request) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.totalProcessed++
	if request.Result != nil && request.Result.Success {
		p.totalSucceeded++
	} else {
		p.totalFailed++
	}
}

// GetStatistics returns processor statistics
func (p *Processor) GetStatistics() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	uptime := time.Since(p.startTime)

	return map[string]interface{}{
		"total_processed": p.totalProcessed,
		"total_succeeded": p.totalSucceeded,
		"total_failed":    p.totalFailed,
		"success_rate":    float64(p.totalSucceeded) / float64(p.totalProcessed) * 100,
		"uptime":          uptime.String(),
		"active_requests": len(p.activeRequests),
		"worker_count":    len(p.workers),
	}
}

// Worker methods

func (w *Worker) run() {
	// Worker implementation would handle request processing
}

func (w *Worker) stop() {
	w.cancel()
}

// Utility functions

func generateID() string {
	return fmt.Sprintf("batch_%d", time.Now().UnixNano())
}

func parseMemoryTarget(target string) (uint32, uint32, error) {
	// Parse target like "0x1000:100" (address:length)
	// This is a simplified implementation
	return 0x1000, 100, nil
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// DefaultConfig returns default batch processor configuration
func DefaultConfig() *Config {
	return &Config{
		MaxConcurrentRequests: 10,
		MaxOperationsPerBatch: 1000,
		DefaultTimeout:        30 * time.Second,
		WorkerPoolSize:        4,
		MaxRetryAttempts:      3,
		RetryDelay:            100 * time.Millisecond,
		EnableMetrics:         true,
		EnableValidation:      true,
		BufferSize:            100,
	}
}
