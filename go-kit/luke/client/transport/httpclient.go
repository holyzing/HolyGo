package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/holyzing/HolyGo/go-kit/luke/client/api"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

// NewHTTPClient returns an AddService backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, options ...httptransport.RequestFunc) (Set, error) {
	var clientOptions = []httptransport.ClientOption{}

	for _, option := range options {
		clientOptions = append(clientOptions, httptransport.ClientBefore(option))
	}
	u, err := url.Parse(instance)
	if err != nil {
		return Set{}, err
	}

	// Each individual endpoint is an http/transport.Client (which implements
	// endpoint.Endpoint) that gets wrapped with various middlewares. If you
	// made your own client library, you'd do this work there, so your server
	// could rely on a consistent set of client behavior.
	var createJobEndpoint endpoint.Endpoint
	{
		createJobEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/v1/jobs"),
			encodeHTTPCreateJobRequest,
			decodeHTTPCreateJobResponse,
			clientOptions...,
		).Endpoint()
	}

	var getJobEndpoint endpoint.Endpoint
	{
		getJobEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/v1/jobs"),
			encodeHTTPGetJobRequest,
			decodeHTTPGetJobResponse,
			clientOptions...,
		).Endpoint()
	}

	// Returning the endpoint.Set as a service.Service relies on the
	// endpoint.Set implementing the Service methods. That's just a simple bit
	// of glue code.
	return Set{
		CreateJobEndpoint: createJobEndpoint,
		GetJobEndpoint:    getJobEndpoint,
	}, nil
}

// encodeHTTPCreateJobRequest is a transport/http.EncodeRequestFunc
// that encodes a generic request into the various portions of
// the http request (path, query, and body).
func encodeHTTPCreateJobRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*api.CreateJobInput)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	// path := strings.Join([]string{
	// 	"",
	// 	"luke",
	// 	"v1",
	// 	"jobs",
	// }, "/")
	// u, err := url.Parse(path)
	// if err != nil {
	// 	return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	// }
	// r.URL.RawPath = u.RawPath
	// r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	values.Set("user", req.User)
	values.Set("organization", req.Organization)

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(req); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", req)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeHTTPGetJobRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*api.GetJobInput)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters

	// Set the query parameters
	values := r.URL.Query()
	values.Set("user", req.User)
	values.Set("organization", req.Organization)
	values.Set("id", strconv.FormatInt(req.ID, 10))
	values.Set("handle", req.Handle)
	values.Set("fields", strings.Join(req.Fields, ","))

	r.URL.RawQuery = values.Encode()

	return nil
}

// decodeHTTPCreateJobResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded LukeResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func decodeHTTPCreateJobResponse(_ context.Context, r *http.Response) (interface{}, error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if len(buf) == 0 {
		return nil, errors.New("response http body empty")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp api.CreateJobOutput
	if err = json.Unmarshal(buf, &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// decodeHTTPGetJobResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded LukeResponse response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func decodeHTTPGetJobResponse(_ context.Context, r *http.Response) (interface{}, error) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if len(buf) == 0 {
		return nil, errors.New("response http body empty")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp api.GetJobOutput
	if err = json.Unmarshal(buf, &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// ------------------------------------------------

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path += path
	return &next
}

func errorDecoder(buf []byte) error {
	var w errorWrapper
	if err := json.Unmarshal(buf, &w); err != nil {
		const size = 8196
		if len(buf) > size {
			buf = buf[:size]
		}
		return fmt.Errorf("response body '%s': cannot parse non-json request body", buf)
	}

	return errors.New(w.String())
}

type ErrMsg struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

type errorWrapper struct {
	Error ErrMsg `json:"Error"`
}

func (e errorWrapper) String() string {
	return e.Error.Code + ":" + e.Error.Message
}
