// Package honeycomb wraps our honeycomb (monitoring + tracing) initialization and provides some helper functions to make instrumentation easier.
package honeycomb

import (
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/client"
	"github.com/honeycombio/beeline-go/trace"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/honeycombio/libhoney-go"

	"github.com/pubgolf/pubgolf/api/internal/lib/config"
)

const (
	serviceName = "pubgolf-app-server"

	datasetDefault = "pubgolf-dev"
	datasetStaging = "pubgolf-staging"
	datasetProd    = "pubgolf"
)

var (
	_client *libhoney.Client
)

// WrapMux provides a top-level event for all inbound HTTP requests, even if they don't have a wrapped handler.
func WrapMux(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(hnynethttp.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.GetSpanFromContext(r.Context())
		if span != nil {
			span.AddField("name", "inbound_http_request")
		}
		mux.ServeHTTP(w, r)
	}))
}

// WrapHandlerFunc instruments an HTTP handler func. It should be applied on the handler itself (not as a middleware or on the mux) to allow introspection of the handler function's name.
func WrapHandlerFunc(hf func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	fnName := runtime.FuncForPC(reflect.ValueOf(hf).Pointer()).Name()
	return hnynethttp.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.GetSpanFromContext(r.Context())
		if span != nil {
			span.AddField("handler_func_name", fnName)
			span.AddField("name", strings.TrimSuffix(filepath.Base(fnName), ".func1"))
		}
		hf(w, r)
	})
}

// Init sets up Honeycomb based on the env-specific configuration and returns a flush function that must be called on program exit. Init should only be called once during the lifecycle of an application.
func Init(cfg *config.App) func() {
	var closer = beeline.Close

	if _client != nil {
		return closer
	}

	beeline.Init(beeline.Config{
		WriteKey:    cfg.HoneycombWriteKey,
		Dataset:     datasetForEnv(cfg.EnvName),
		ServiceName: serviceName,
	})

	_client = client.Get()

	addGlobalFields(cfg, _client)

	return closer
}

// AddNonemptyField adds a global field to the Honeycomb client if the value is non-empty and Honeycomb has been initialized.
func AddNonemptyField(name string, value string) {
	if _client != nil && len(value) > 0 {
		_client.AddField(name, value)
	}
}

// datasetForEnv labels non-production data so it can be analyzed separately.
func datasetForEnv(env config.DeployEnv) string {
	if env == config.DeployEnvProd {
		return datasetProd
	}

	if env == config.DeployEnvStaging {
		return datasetStaging
	}

	return datasetDefault
}

// addGlobalFields configures service-level info for each span.
func addGlobalFields(cfg *config.App, cl *libhoney.Client) {
	startTime := time.Now()
	cl.AddDynamicField("meta.process_uptime_sec", func() interface{} {
		return time.Now().Sub(startTime) / time.Second
	})

	// Metadata related to deployment and versioning.
	cl.AddField("meta.env", cfg.EnvName)

	// Runtime and utilization data.
	cl.AddDynamicField("meta.num_goroutines", func() interface{} { return runtime.NumGoroutine() })
	cl.AddDynamicField("meta.memory_in_use", func() interface{} {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem.Alloc
	})
}
