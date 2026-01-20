package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	openapicommon "k8s.io/kube-openapi/pkg/common"
)

// testPortCounter is used to generate unique ports for each test to avoid conflicts
var testPortCounter = 9000

// newTestConfig creates a minimal valid Config for testing
func newTestConfig() *Config {
	testScheme := runtime.NewScheme()
	codecFactory := serializer.NewCodecFactory(testScheme)

	// Use unique ports for each test to avoid "address already in use" errors
	testPortCounter++
	httpPort := testPortCounter
	httpsPort := testPortCounter + 1
	testPortCounter++ // increment again to keep ports separate

	return &Config{
		Name:            "test-server",
		Version:         "v1",
		Scheme:          testScheme,
		CodecFactory:    &codecFactory,
		HTTPListenPort:  httpPort,
		HTTPSListenPort: httpsPort,
		OpenAPIConfig: func(ref openapicommon.ReferenceCallback) map[string]openapicommon.OpenAPIDefinition {
			return map[string]openapicommon.OpenAPIDefinition{}
		},
	}
}

func TestRateLimitingConfigCustomValues(t *testing.T) {
	config := newTestConfig()
	config.MaxRequestsInFlight = 600
	config.MaxMutatingRequestsInFlight = 300

	server, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, server)

	// Verify custom rate limiting values are applied
	assert.Equal(t, 600, server.Config.Config.MaxRequestsInFlight)
	assert.Equal(t, 300, server.Config.Config.MaxMutatingRequestsInFlight)
}

func TestRateLimitingConfigDefaultValues(t *testing.T) {
	config := newTestConfig()
	config.MaxRequestsInFlight = 0
	config.MaxMutatingRequestsInFlight = 0

	server, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, server)

	// When 0 is specified, k8s.io/apiserver defaults should be used.
	// The ServerRunOptions defaults are 400 for MaxRequestsInFlight
	// and 200 for MaxMutatingRequestsInFlight (set by options.NewServerRunOptions())
	assert.Equal(t, 400, server.Config.Config.MaxRequestsInFlight)
	assert.Equal(t, 200, server.Config.Config.MaxMutatingRequestsInFlight)
}

func TestRateLimitingConfigPartialOverride(t *testing.T) {
	config := newTestConfig()
	config.MaxRequestsInFlight = 500
	config.MaxMutatingRequestsInFlight = 0

	server, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, server)

	// Custom value for MaxRequestsInFlight
	assert.Equal(t, 500, server.Config.Config.MaxRequestsInFlight)
	// Default for MaxMutatingRequestsInFlight
	assert.Equal(t, 200, server.Config.Config.MaxMutatingRequestsInFlight)
}

func TestServerConfigComplete(t *testing.T) {
	config := &Config{}
	config.complete()

	// Verify defaults are set correctly
	assert.Equal(t, 8080, config.HTTPListenPort)
	assert.Equal(t, 8081, config.HTTPSListenPort)
	assert.Equal(t, []string{"watch", "proxy"}, config.LongRunningVerbs)
	assert.NotNil(t, config.Scheme)
	assert.NotNil(t, config.CodecFactory)
	assert.Equal(t, "mink", config.Name)
	assert.NotNil(t, config.DefaultOptions)
}

func TestServerCreation(t *testing.T) {
	config := newTestConfig()
	config.Version = "v1.0.0"

	server, err := New(config)
	require.NoError(t, err)
	require.NotNil(t, server)
	require.NotNil(t, server.GenericAPIServer)
	require.NotNil(t, server.Config)

	// Verify server name and version are set
	assert.Equal(t, "test-server", server.Config.OpenAPIConfig.Info.Title)
	assert.Equal(t, "v1.0.0", server.Config.OpenAPIConfig.Info.Version)
}
