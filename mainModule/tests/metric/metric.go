package metric

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{1, 5, 10, 30, 60, 300},
		},
		[]string{"code", "method", "path"},
	)

	metricRegisterOnce sync.Once
)

type router struct {
	httprouter.Router
}

type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(BeforeFunc)
}

// BeforeFunc is a function that is called before the ResponseWriter has been written to.
type BeforeFunc func(ResponseWriter)

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	newRw := responseWriter{rw, 0, 0, nil}
	return &newRw
}

func (rw *responseWriter) Before(before BeforeFunc) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []BeforeFunc
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func HandlerWrapper(path string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := NewResponseWriter(w)
		now := time.Now()
		next.ServeHTTP(d, r)
		requestCounter.With(
			prometheus.Labels{
				"code":   fmt.Sprintf("%d", d.Status()),
				"method": r.Method,
				"path":   path,
			},
		).Inc()

		requestDuration.With(
			prometheus.Labels{
				"code":   fmt.Sprintf("%d", d.Status()),
				"method": r.Method,
				"path":   path,
			},
		).Observe(time.Since(now).Seconds())
	}
}

func (r *router) Handler(method string, path string, handler http.Handler) {
	r.Router.Handler(
		method,
		path,
		promhttp.InstrumentHandlerInFlight(
			inFlightGauge,
			HandlerWrapper(
				path,
				handler,
			),
		),
	)
}

func NewRouterWithMetrics() *router {

	metricRegisterOnce.Do(func() {
		prometheus.MustRegister(inFlightGauge, requestCounter, requestDuration)
	})

	return &router{
		*httprouter.New(),
	}
}
