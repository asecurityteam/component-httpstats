package httpstats

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/asecurityteam/settings/v2"
	"github.com/rs/xstats"
	"github.com/stretchr/testify/assert"
)

const (
	testdependency = "testdependency"
)

func TestMetricsConfigName(t *testing.T) {
	Config := MetricsConfig{}
	assert.Equal(t, "metrics", Config.Name())
}

func TestMetricsComponentSettings(t *testing.T) {
	metricsComponent := &MetricsComponent{}
	assert.IsType(t, &MetricsConfig{}, metricsComponent.Settings())
}

func TestMetricsComponentNew(t *testing.T) {
	component := NewComponent()
	assert.IsType(t, &MetricsComponent{}, component)
}

func TestNew(t *testing.T) {
	src := settings.NewMapSource(map[string]interface{}{
		"metrics": map[string]interface{}{
			"backend": testdependency,
		},
	})
	wrapper, err := New(context.Background(), src)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	assert.Equal(
		t, fmt.Sprintf("client_dependency:%s", testdependency),
		reflect.Indirect(reflect.ValueOf(tr)).FieldByName("tags").Index(0).String(),
	)
}

func TestMetricsComponentNewNoDependencyError(t *testing.T) {
	cmp := NewComponent()
	conf := cmp.Settings()
	wrapper, err := cmp.New(context.Background(), conf)

	assert.IsType(t, err, ErrNoBackend)
	assert.Nil(t, wrapper)
}

func TestMetricsComponentNewStaticPath(t *testing.T) {
	cmp := NewComponent()
	conf := cmp.Settings()
	conf.Backend = testdependency
	conf.Path = "/foo"
	wrapper, err := cmp.New(context.Background(), conf)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	tags := reflect.Indirect(reflect.ValueOf(tr)).FieldByName("tags")
	for i := 0; i < tags.Len(); i++ {
		assert.Contains(t, []string{fmt.Sprintf("client_dependency:%s", testdependency), "client_path:/foo"}, tags.Index(i).String())
	}
}

func TestMetricsComponentNonStaticPath(t *testing.T) {
	cmp := NewComponent()
	conf := cmp.Settings()
	conf.Backend = testdependency
	wrapper, err := cmp.New(context.Background(), conf)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)

	request := http.Request{URL: &url.URL{Path: `/some/random/path`}}

	xStater := XStater{}
	xStaterContext := xstats.NewContext(context.Background(), &xStater)

	request = *request.WithContext(xStaterContext)

	_, _ = tr.RoundTrip(&request)

	expectedTags := []string{
		"client_path:/some/random/path",
		fmt.Sprintf("client_dependency:%s", testdependency),
		"method:",
		"status_code:502",
		"status:error",
	}

	assert.Equal(t, len(expectedTags), len(xStater.TimingTags))

	for i := 0; i < len(expectedTags); i++ {
		assert.Contains(t, expectedTags, xStater.TimingTags[i])
	}
}

func TestMetricsComponentNewStaticPathOmitPath(t *testing.T) {
	cmp := NewComponent()
	conf := cmp.Settings()
	conf.Backend = testdependency
	conf.Path = "/foo"
	conf.OmitPathTag = true
	wrapper, err := cmp.New(context.Background(), conf)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	tags := reflect.Indirect(reflect.ValueOf(tr)).FieldByName("tags")
	for i := 0; i < tags.Len(); i++ {
		assert.Contains(t, []string{fmt.Sprintf("client_dependency:%s", testdependency)}, tags.Index(i).String())
	}
}

func TestMetricsComponentNewNoPathGiven(t *testing.T) {
	cmp := NewComponent()
	conf := cmp.Settings()
	conf.Backend = testdependency
	wrapper, err := cmp.New(context.Background(), conf)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	assert.Equal(t, 1, reflect.Indirect(reflect.ValueOf(tr)).FieldByName("requestTaggers").Len())
}

// XStater a rudimentary recorder of emitted stats things.  Add recordings for the stat name and values if you need them
type XStater struct {
	AddTagsTags   []string
	CountTags     []string
	GaugeTags     []string
	HistogramTags []string
	TimingTags    []string
}

func (x *XStater) AddTags(tags ...string)                                  { x.AddTagsTags = tags }
func (x *XStater) GetTags() []string                                       { return x.AddTagsTags }
func (x *XStater) Gauge(stat string, value float64, tags ...string)        { x.GaugeTags = tags }
func (x *XStater) Count(stat string, count float64, tags ...string)        { x.CountTags = tags }
func (x *XStater) Histogram(stat string, value float64, tags ...string)    { x.HistogramTags = tags }
func (x *XStater) Timing(stat string, value time.Duration, tags ...string) { x.TimingTags = tags }
