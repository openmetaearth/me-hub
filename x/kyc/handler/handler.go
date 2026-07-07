package handler

import (
	"context"
	"cosmossdk.io/errors"
	"sort"
	"sync"
)

type HandlerFunc func(ctx context.Context, eventType string, data interface{}) error

type SortHandler struct {
	Priority int
	Module   string
	Handler  HandlerFunc
}

// HandlerRegistry holds the registered event handlers and execute them before kyc transaction completion
type HandlerRegistry struct {
	mu       sync.RWMutex
	handlers map[string][]SortHandler
}

// NewEventRegistry creates a new HandlerRegistry
func NewEventRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string][]SortHandler),
	}
}

// RegisterHandler registers an event handler for a specific event type with a priority
func (r *HandlerRegistry) RegisterHandler(eventType string, priority int, module string, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	handlers := append(r.handlers[eventType], SortHandler{Priority: priority, Module: module, Handler: handler})

	// sorting handlers by priority
	sort.SliceStable(handlers, func(i, j int) bool {
		return handlers[i].Priority < handlers[j].Priority
	})

	r.handlers[eventType] = handlers
}

// HandleEvent calls the registered handlers for the given event type in order of their priority
func (r *HandlerRegistry) HandleEvent(ctx context.Context, eventType string, data interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handlers, found := r.handlers[eventType]
	if !found {
		return nil
	}

	for _, handler := range handlers {
		if err := handler.Handler(ctx, eventType, data); err != nil {
			return errors.Wrapf(err,
				"failed to handle event by module %s, priority %d", handler.Module, handler.Priority)
		}
	}

	return nil
}
