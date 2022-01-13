package v1alpha1

// CreateOptions are options that may be provided when creating an API object
type CreateOptions struct {
	IsInCluster bool
}

// DeleteOptions are options that may be provided when deleting an API object
type DeleteOptions struct {
	IsInCluster bool
	ShouldKill  bool
}

// PipekitMeta holds metadata specific to interacting with the Pipekit API
type PipekitMeta struct {
	PipeName           string
	UserId             string
	PipeId             string
	RunId              string
	Cluster            string
	SecretsEnvironment string
	Namespace          string
	Tags               []string
}
