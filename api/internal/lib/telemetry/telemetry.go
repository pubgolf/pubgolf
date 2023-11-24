package telemetry

import (
	"context"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

// ServiceName provides the attribute value for "service.name".
const ServiceName = "pubgolf-api-server"

// Init configures OTel and Honeycomb reporting.
func Init(cfg *config.App) (func(), error) {
	return otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithSpanProcessor(honeycomb.NewBaggageSpanProcessor()),
		otelconfig.WithServiceName(ServiceName),
		otelconfig.WithServiceVersion(gitVersion()),
		honeycomb.WithApiKey(cfg.HoneycombKey),
	)
}

// AddRecursiveAttribute adds an attribute to a span and all of its children.
func AddRecursiveAttribute(ctx *context.Context, key, value string) {
	// Set attribute on current span.
	span := trace.SpanFromContext(*ctx)
	span.SetAttributes(attribute.String(key, value))

	// Add to baggage so child spans will receive the attribute as well.
	bag := baggage.FromContext(*ctx)
	multiSpanAttribute, _ := baggage.NewMember(key, value)
	bag, _ = bag.SetMember(multiSpanAttribute)
	*ctx = baggage.ContextWithBaggage(*ctx, bag)
}

// gitVersion returns the git sha of the compiled server binary.
func gitVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return "unknown"
}

// FnSpan annotates a function with a span for tracing, automatically inferring the name.
var FnSpan = AutoSpan("func")

// AutoSpan generates an auto instrumentation for a function. For example, if you define
//
// var mySpan = telemetry.AutoSpan("my")
//
// then you would call it by adding `defer mySpan(&ctx)()` to the start of the function definition you want to annotate, resulting in a span named `my.NameOfAnnotatedFunc`.
func AutoSpan(prefix string) func(ctx *context.Context) func() {
	return func(ctx *context.Context) func() {
		name := prefix + ".Unknown"
		if pc, _, _, ok := runtime.Caller(1); ok {
			name = prefix + "." + strings.Split(filepath.Base(runtime.FuncForPC(pc).Name()), ".")[2]
		}
		newCtx, span := otel.Tracer("").Start(*ctx, name)
		*ctx = newCtx
		return func() { span.End() }
	}
}
