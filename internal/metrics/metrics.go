package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Gatherer interface {
	Inc(string, string)
}
type HttpRequestsCounter struct {
	counter *prometheus.CounterVec
}

func NewHttpRequestsCounter() *HttpRequestsCounter {
	httpReqs := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_http_requests_total",
			Help: "How many http request processed, partioned by resource and method",
		},
		[]string{"resource", "method"},
	)
	//prometheus.MustRegister(httpReqs)
	return &HttpRequestsCounter{
		counter: httpReqs,
	}
}

func (hrc *HttpRequestsCounter) Inc(resource, method string) {
	hrc.counter.WithLabelValues(resource, method).Inc()
}
