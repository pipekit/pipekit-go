package pipes

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pipekit/pipekit-go/v0"
	meta "github.com/pipekit/pipekit-go/v0/meta/v1alpha1"
)

// Client holds methods related to workflows
type Client interface {
	Create(ctx context.Context, pipe *pipekit.Pipe, opts meta.CreateOptions) (*pipekit.Pipe, error)
	Stop(ctx context.Context, userId, pipeId, runId string, opts meta.DeleteOptions) error
}

type client struct {
	backend pipekit.Backend
}

// Config holds config for the workflows client
type Config struct {
	Backend pipekit.Backend
}

// New instantiates a workflows client
func New(config *Config) Client {
	return &client{
		backend: config.Backend,
	}
}

// Create runs a workflow either through the Pipekit endpoint or directly
// on the cluster using the Argo server.
func (c *client) Create(ctx context.Context, pipe *pipekit.Pipe, opts meta.CreateOptions) (*pipekit.Pipe, error) {
	var err error

	if opts.IsInCluster {
		err = c.createOnCluster(ctx, pipe)
	} else {
		err = c.createThroughPipekit(ctx, pipe)
	}

	return pipe, err
}

func (c *client) createThroughPipekit(ctx context.Context, pipe *pipekit.Pipe) error {
	path := "events-handler/v1/users/%s/pipes/%s/runs"

	path = pipekit.FormatURLPath(path, pipe.GetUserId(), pipe.GetPipeId())

	return c.backend.Call(ctx, http.MethodPost, path, nil, pipe)
}

func (c *client) createOnCluster(ctx context.Context, pipe *pipekit.Pipe) error {
	path := "plumbing/v1/users/%s/pipes/%s/runs"

	path = pipekit.FormatURLPath(path, pipe.GetUserId(), pipe.GetPipeId())

	return c.backend.Call(ctx, http.MethodPost, path, nil, pipe)
}

// Stop stops a running workflow either through the Pipekit endpoint or directly
// on the cluster using the Argo server.
func (c *client) Stop(ctx context.Context, userId, pipeId, runId string, opts meta.DeleteOptions) error {
	if opts.IsInCluster {
		return c.stopOnCluster(ctx, userId, pipeId, runId, opts.ShouldKill)
	}

	return c.stopThroughPipekit(ctx, userId, pipeId, runId, opts.ShouldKill)
}

func (c *client) stopThroughPipekit(ctx context.Context, userId, pipeId, runId string, shouldKill bool) error {
	path := "events-handler/v1/users/%s/pipes/%s/runs/%s"

	path = pipekit.FormatURLPath(path, userId, pipeId, runId)

	params := formatSigTermOrKillParams(shouldKill)

	return c.backend.Call(ctx, http.MethodDelete, path, params, nil)
}

func (c *client) stopOnCluster(ctx context.Context, userId, pipeId, runId string, shouldKill bool) error {
	path := "plumbing/v1/users/%s/pipes/%s/runs/%s"

	path = pipekit.FormatURLPath(path, userId, pipeId, runId)

	params := formatSigTermOrKillParams(shouldKill)

	return c.backend.Call(ctx, http.MethodDelete, path, params, nil)
}

func formatSigTermOrKillParams(shouldKill bool) pipekit.ParamsContainer {
	params := pipekit.Params{
		Values: url.Values{},
	}
	params.Values.Add("should-kill", strconv.FormatBool(shouldKill))
	return &params
}
