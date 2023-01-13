package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/holyzing/HolyGo/go-kit/luke/api/v1"
	"github.com/holyzing/HolyGo/go-kit/luke/endpoint"
	"github.com/holyzing/HolyGo/go-kit/luke/metric"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gogo/protobuf/jsonpb"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const contentType = "application/json; charset=utf-8"

// --------------------------
var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.NewLukeServiceClient
	_ = io.Copy
	// _ = errors.Wrap
)

type errorWrapper struct {
	Error string `json:"error"`
}

var pbMarshaler jsonpb.Marshaler

func init() {
	pbMarshaler = jsonpb.Marshaler{
		//EmitDefaults: true,
		// OrigName: true,
		// EnumsAsInts: true,
	}
}

// --------------------------

type httpError struct {
	error
	statusCode int
	headers    map[string][]string
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	println("---- transport_http errorEncoder ----")
	body, _ := jsoniter.Marshal(errorWrapper{Error: err.Error()})
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

func MakeHTTPHandler(endpoints endpoint.LukeEndPoints) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext, remoteAddrToContext),
		// opentracing.HTTPToContext(trace.OtTracer, "MakeHTTPHandler", log.KLog)

		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	}

	r := metric.NewRouterWithMetrics()
	r.Handler("GET", "/v1/job", httptransport.NewServer(
		endpoints.JobReadEndPoint,
		DecodeGetJobRequest,
		EncodeHTTPResponse,
		serverOptions...,
	))

	// r.Handler("GET", "/v1/jobs", httptransport.NewServer(
	// 	endpoints.JobReadEndPoint,
	// 	DecodeGetJobsRequest,
	// 	EncodeHTTPResponse,
	// 	serverOptions...,
	// ))

	r.Handler("POST", "/v1/job", httptransport.NewServer(
		endpoints.JobWriteEndPoint,
		DecodeCreateJobRequest,
		EncodeHTTPResponse,
		serverOptions...,
	))

	return r

}

func DecodeGetJobRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	println("---- transport_http DecodeHTTPGenericGetJobRequest ----")
	req := new(pb.GetJobRequest)

	pathParams := httprouter.ParamsFromContext(ctx)
	_ = pathParams
	queryParams := r.URL.Query()

	ids := queryParams.Get("id")
	if ids != "" {
		id, err := strconv.ParseInt(ids, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting ids: %v", ids))
		}
		req.Id = id
	}
	req.Handle = queryParams.Get("handle")
	field := queryParams.Get("fields")
	if field != "" {
		req.Fields = strings.Split(field, ",")
	}
	GenericCombined := queryParams.Get("combined")
	if strings.EqualFold("true", GenericCombined) {
		req.Combined = true
	}

	return &pb.LukeRequest{
		Method:       "GetJob",
		User:         queryParams.Get("user"),
		Organization: queryParams.Get("organization"),
		TenantName:   queryParams.Get("tenant_name"),
		Body: &pb.LukeRequest_GetRequest{
			GetRequest: req,
		},
	}, nil
}

func DecodeCreateJobRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	println("---- transport_http DecodeHTTPGenericCreateJobRequest ----")
	var req = new(pb.CreateJobRequest)
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		if err = jsonpb.Unmarshal(bytes.NewBuffer(buf), req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{fmt.Errorf("request body '%s': %s", buf, err.Error()),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := httprouter.ParamsFromContext(ctx)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &pb.LukeRequest{
		Method:       "CreateJob",
		User:         queryParams.Get("user"),
		Organization: queryParams.Get("organization"),
		TenantName:   queryParams.Get("tenant_name"),
		Body: &pb.LukeRequest_CreateRequest{
			CreateRequest: req,
		},
	}, err
}

func EncodeHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	println("---- transport_http EncodeHTTPJobResponse ----")
	resp := response.(*pb.LukeResponse)
	w.WriteHeader(int(resp.Retcode))
	// 将 response 序列化为 JSON 写入到 http response
	return pbMarshaler.Marshal(w, resp)
}

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	println("---- transport_http headersToContext ----")
	for k := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	// Tune specific change.
	// also add the request url
	// collisions
	ctx = context.WithValue(ctx, "request-url", r.URL.Path)
	ctx = context.WithValue(ctx, "transport", "HTTPJSON")

	return ctx
}

func remoteAddrToContext(ctx context.Context, r *http.Request) context.Context {
	println("---- transport_http remoteAddrToContext ----")
	return context.WithValue(ctx, "http-remote-addr", r.RemoteAddr)
}
