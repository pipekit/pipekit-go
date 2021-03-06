package pipekit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	defaultBaseURI        string = "https://pipekit.io/api"
	defaultClusterBaseURI string = "http://localhost:8080/api"
)

// Backend exposes for methods for calling pipekit apis
type Backend interface {
	Call(ctx context.Context, method, path string, params ParamsContainer, obj interface{}) error
}

type backend struct {
	baseURI        string
	clusterBaseURI string
	httpClient     *http.Client
}

// BackendConfig helps configure a new pipekit backend
type BackendConfig struct {
	BaseURI        *string
	ClusterBaseURI *string
	HTTPClient     *http.Client
}

// NewBackend instantiates a new backend
func NewBackend(config *BackendConfig) Backend {
	if config == nil {
		config = newDefaultConfig()
	}

	if config.BaseURI == nil {
		config.BaseURI = &defaultBaseURI
	}

	if config.ClusterBaseURI == nil {
		config.ClusterBaseURI = &defaultClusterBaseURI
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{}
	}

	return &backend{
		baseURI:        *config.BaseURI,
		clusterBaseURI: *config.ClusterBaseURI,
		httpClient:     config.HTTPClient,
	}
}

func newDefaultConfig() *BackendConfig {
	httpClient := http.Client{}

	return &BackendConfig{
		BaseURI:        &defaultBaseURI,
		ClusterBaseURI: &defaultClusterBaseURI,
		HTTPClient:     &httpClient,
	}
}

// Call invokes the backend for a pipekit api
func (b *backend) Call(ctx context.Context, method, path string, params ParamsContainer, body interface{}) error {
	formattedPath, err := b.formatURL(path, params)
	if err != nil {
		return err
	}

	req, err := formatJSONRequest(ctx, method, formattedPath, body)
	if err != nil {
		return err
	}

	_, err = executeRequest(b.httpClient, req, body)
	return err
}

// formatURL does two things:
// 1. prepends the base URI to the path - any calls to the pipekit daemon
// are routed to the cluster base uri, and
// 2. appends the query params
func (b *backend) formatURL(urlPath string, params ParamsContainer) (string, error) {
	uri := b.baseURI
	if strings.Contains(urlPath, "plumbing") {
		uri = b.clusterBaseURI
	}

	u, err := url.Parse(uri)
	if err != nil {
		return "", nil
	}

	u.Path = path.Join(u.Path, urlPath)
	u.RawQuery = params.GetParams().Values.Encode()

	return u.String(), nil
}

// executeRequest executes the formatted request
func executeRequest(client *http.Client, request *http.Request, obj interface{}) (*int, error) {
	var err error

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
	}()

	if obj != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, obj)
		if err != nil {
			return nil, err
		}
	}

	err = makeErrorIf400sOr500sStatus(resp.StatusCode)
	if err != nil {
		return nil, err
	}

	return &resp.StatusCode, nil
}

// makeErrorIf400sOr500sStatus generates an error if the http status code reads
// an error
func makeErrorIf400sOr500sStatus(statusCode int) error {
	if statusCode >= http.StatusBadRequest {
		return fmt.Errorf("HTTP status: %d", statusCode)
	}

	return nil
}

// formatJSONRequest formats a JSON request
func formatJSONRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Close = true

	formatJSONHeader(ctx, req)

	return req, nil
}

func formatJSONHeader(ctx context.Context, req *http.Request) {
	req.Header.Set("Content-Type", "application/json")

	// TODO: Get proper context/token passing
	// tokenString, ok := GetAuthTokenFromContext(ctx)
	// if !ok {
	// 	return
	// }

	// bearerToken := guaranteeBearerPrefix(tokenString)

	bearerToken := "placeholder"
	req.Header.Set("Authorization", bearerToken)
}

// FormatURLPath is a convenience method for ensuring that paths are properly
// formatted
func FormatURLPath(format string, params ...string) string {
	// Convert parameters to interface{} and URL-escape them
	untypedParams := make([]interface{}, len(params))
	for i, param := range params {
		untypedParams[i] = interface{}(url.QueryEscape(param))
	}

	return fmt.Sprintf(format, untypedParams...)
}
