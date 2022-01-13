package pipekit

import (
	meta "github.com/pipekit/pipekit-go/v0/meta/v1alpha1"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

// Pipe holds Pipekit specific metadata and a workflow definition to
// be run.
type Pipe struct {
	Pipekit meta.PipekitMeta
	Argo    *v1alpha1.Workflow
}

// GetUserId is a convenience method to get the user id
func (p *Pipe) GetUserId() string {
	return p.Pipekit.UserId
}

// GetPipeId is a convenience method to get the pipe id
func (p *Pipe) GetPipeId() string {
	return p.Pipekit.PipeId
}

// GetRunId is a convenience mehtod to get the run id
func (p *Pipe) GetRunId() string {
	return p.Pipekit.RunId
}
