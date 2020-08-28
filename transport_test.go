package httpstats

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/asecurityteam/settings"
	"github.com/stretchr/testify/assert"
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
			"backend": "testdependency",
		},
	})
	wrapper, err := New(context.Background(), src)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	assert.Equal(
		t, "client_dependency:testdependency",
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
	conf.Backend = "testdependency"
	conf.Path = "/foo"
	wrapper, err := cmp.New(context.Background(), conf)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	tags := reflect.Indirect(reflect.ValueOf(tr)).FieldByName("tags")
	for i := 0; i < tags.Len(); i++ {
		assert.Contains(t, []string{"client_dependency:testdependency", "client_path:/foo"}, tags.Index(i).String())
	}
}

func TestMetricsComponentNewNoPathGiven(t *testing.T) {
	cmp := NewComponent()
	conf := cmp.Settings()
	conf.Backend = "testdependency"
	wrapper, err := cmp.New(context.Background(), conf)
	assert.Nil(t, err)
	assert.IsType(t, wrapper, func(next http.RoundTripper) http.RoundTripper { return nil }, wrapper)

	tr := wrapper(http.DefaultTransport)
	assert.Equal(t, 1, reflect.Indirect(reflect.ValueOf(tr)).FieldByName("requestTaggers").Len())
}
