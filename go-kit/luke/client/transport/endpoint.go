package transport

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreateJobEndpoint endpoint.Endpoint
	GetJobEndpoint    endpoint.Endpoint
}

// LabeledMiddleware will get passed the endpoint name when passed to
// WrapAllLabeledExcept, this can be used to write a generic metrics
// middleware which can send the endpoint name to the metrics collector.
type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

func (s *Set) WrapAllLabeledExcept(middleware LabeledMiddleware, excluded ...string) {
	included := map[string]struct{}{
		"CreateJob": {},
		"GetJob":    struct{}{},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "CreateJob" {
			s.CreateJobEndpoint = middleware("CreateJob", s.CreateJobEndpoint)
		}

		if inc == "GetJob" {
			s.GetJobEndpoint = middleware("GetJob", s.GetJobEndpoint)
		}
	}
}
