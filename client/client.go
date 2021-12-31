package client

import (
	"github.com/pipekit/pipekit-go/v0"
	"github.com/pipekit/pipekit-go/v0/pipes"
	workflows "github.com/pipekit/pipekit-go/v0/pipes"
)

// API holds a composite of all pipekit apis
type API struct {
	Pipes pipes.Client
}

// Config holds api configuration
type Config struct{}

// New creates a new API instance
func New(config *Config) API {
	backend := pipekit.NewBackend(&pipekit.BackendConfig{})
	pipesClient := pipes.New(&workflows.Config{Backend: backend})

	return API{
		Pipes: pipesClient,
	}
}
