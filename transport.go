package httpstats

import (
	"context"
	"errors"
	"net/http"

	statstransport "github.com/asecurityteam/httpstats/v2"
	"github.com/asecurityteam/settings/v2"
)

// MetricsConfig contains settings for request metrics emissions.
type MetricsConfig struct {
	Timing            string `description:"Name of overall timing metric."`
	DNS               string `description:"Name of DNS timing metric."`
	TCP               string `description:"Name of TCP timing metric."`
	ConnectionIdle    string `description:"Name of idle timing metric."`
	TLS               string `description:"Name of TLS timing metric."`
	WroteHeaders      string `description:"Name of time to write headers metric."`
	FirstResponseByte string `description:"Name of time to first resposne byte metrics."`
	BytesReceived     string `description:"Name of bytes received metric."`
	BytesSent         string `description:"Name of bytes sent metric."`
	BytesTotal        string `description:"Name of bytes sent and received metric."`
	PutIdle           string `description:"Name of idle connection return count metric."`
	BackendTag        string `description:"Name of the tag containing the backend reference."`
	PathTag           string `description:"Name of the tag containing the path reference."`
	Backend           string `description:"Static value for the backend tag metric."`
	Path              string `description:"Static value for the path tag metric. If not specified, will be generated for each request."`
	OmitPathTag       bool   `description:"Boolean to omit path tag metric.  Omitting the path tag may be desirable to avoid cardinality explosions for paths that vary greatly."`
}

// Name of the config root.
func (*MetricsConfig) Name() string {
	return "metrics"
}

// MetricsComponent implements the settings.Component interface.
type MetricsComponent struct{}

// NewComponent populates default values.
func NewComponent() *MetricsComponent {
	return &MetricsComponent{}
}

// Settings generates a config populated with defaults.
func (*MetricsComponent) Settings() *MetricsConfig {
	return &MetricsConfig{
		Timing:            "http.client.timing",
		DNS:               "http.client.dns.timing",
		TCP:               "http.client.tcp.timing",
		ConnectionIdle:    "http.client.connection_idle.timing",
		TLS:               "http.client.tls.timing",
		WroteHeaders:      "http.client.wrote_headers.timing",
		FirstResponseByte: "http.client.first_response_byte.timing",
		BytesReceived:     "http.client.bytes_received",
		BytesSent:         "http.client.bytes_sent",
		BytesTotal:        "http.client.bytes_total",
		PutIdle:           "http.client.put_idle",
		BackendTag:        "client_dependency",
		PathTag:           "client_path",
		Backend:           "",
		Path:              "",
		OmitPathTag:       false,
	}
}

// ErrNoBackend is returned when attempting to initialize a MetricsComponent
// without configuring a value for the Backend tag.
var ErrNoBackend = errors.New("no backend configured")

// New generates the HTTP metrics transport decorator.
//
// If `Dependency` is zero-valued, this constructor will return an error.
//
// If `Path` is zero-valued, this constructor adds an
// `httpstats.TransportOptionRequestTag` that dynamically tags the request
// path on each request made with the decorated `http.Transport`.
func (c *MetricsComponent) New(_ context.Context, conf *MetricsConfig) (func(http.RoundTripper) http.RoundTripper, error) { // nolint
	if conf.Backend == "" {
		return nil, ErrNoBackend
	}

	options := []statstransport.TransportOption{
		statstransport.TransportOptionBytesInName(conf.BytesReceived),
		statstransport.TransportOptionBytesOutName(conf.BytesSent),
		statstransport.TransportOptionBytesTotalName(conf.BytesTotal),
		statstransport.TransportOptionConnectionIdleName(conf.ConnectionIdle),
		statstransport.TransportOptionDNSName(conf.DNS),
		statstransport.TransportOptionFirstResponseByteName(conf.FirstResponseByte),
		statstransport.TransportOptionGotConnectionName(conf.TCP),
		statstransport.TransportOptionPutIdleName(conf.PutIdle),
		statstransport.TransportOptionRequestTimeName(conf.Timing),
		statstransport.TransportOptionTLSName(conf.TLS),
		statstransport.TransportOptionWroteHeadersName(conf.WroteHeaders),
		statstransport.TransportOptionTag(conf.BackendTag, conf.Backend),
	}

	if !conf.OmitPathTag {
		if conf.Path == "" {
			options = append(options, statstransport.TransportOptionRequestTag(func(r *http.Request) (string, string) {
				return conf.PathTag, r.URL.EscapedPath()
			}))
		} else {
			options = append(options, statstransport.TransportOptionTag(conf.PathTag, conf.Path))
		}
	}

	return statstransport.NewTransport(options...), nil
}

// New is the top-level entrypoint for creating an `http.Transport` decorator
// that emits HTTP metrics on every `RoundTrip`.
//
// Useful when configuring a `MetricsComponent` outside of the hierarchy
// of a surrounding application.
func New(ctx context.Context, source settings.Source) (func(http.RoundTripper) http.RoundTripper, error) {
	var dst func(http.RoundTripper) http.RoundTripper
	err := settings.NewComponent(ctx, source, NewComponent(), &dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
