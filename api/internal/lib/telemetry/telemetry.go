package telemetry

import (
	"context"
	"runtime/debug"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-launcher-go/launcher"
	"github.com/pubgolf/pubgolf/api/internal/lib/config"
	"go.opentelemetry.io/otel/baggage"
)

// Init configures OTel and Honeycomb reporting.
func Init(cfg *config.App) (func(), error) {
	return launcher.ConfigureOpenTelemetry(
		launcher.WithSpanProcessor(honeycomb.NewBaggageSpanProcessor()),
		launcher.WithServiceName("pubgolf-api-server"),
		launcher.WithServiceVersion(gitVersion()),
		launcher.WithMetricsEnabled(false),
		honeycomb.WithMetricsDataset("pubgolf-api-server-metrics"),
		honeycomb.WithApiKey(cfg.HoneycombKey),
	)
}

// AddRecursiveAttribute adds an attribute to a span and all of its children.
func AddRecursiveAttribute(ctx *context.Context, key, value string) {
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
