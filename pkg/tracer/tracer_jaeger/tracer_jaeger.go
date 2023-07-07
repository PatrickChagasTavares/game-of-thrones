package tracerjaeger

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
	"go.opentelemetry.io/otel"
	jaegerLib "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semanticConvention "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type (
	jaegerImpl struct {
		provider    *traceSDK.TracerProvider
		serviceName string
	}

	Options struct {
		ServiceName string `json:"service_name" mapstructure:"service_name"`
		EndpointURL string `json:"endpoint_url" mapstructure:"endpoint_url"`
		Username    string `json:"username" mapstructure:"username"`
		Password    string `json:"password" mapstructure:"password"`
		// Sample rate is expressed as 1/X where x is the population size.
		RateSampling int `json:"rate_sampling" mapstructure:"rate_sampling"`
	}
)

func NewExporter(opts Options) tracer.Provider {
	exporter, err := jaegerLib.New(jaegerLib.WithCollectorEndpoint(
		jaegerLib.WithEndpoint(opts.EndpointURL),
	))
	if err != nil {
		log.Fatal("we can't initialize jaeger tracer", err)
	}

	tracerOptions := []traceSDK.TracerProviderOption{
		traceSDK.WithBatcher(
			exporter,
			traceSDK.WithMaxExportBatchSize(traceSDK.DefaultMaxExportBatchSize),
			traceSDK.WithBatchTimeout(traceSDK.DefaultScheduleDelay*time.Microsecond),
		),
		traceSDK.WithResource(resource.NewWithAttributes(
			semanticConvention.SchemaURL,
			semanticConvention.ServiceName(opts.ServiceName),
			semanticConvention.TelemetrySDKLanguageGo,
		)),
	}

	if opts.RateSampling > 1 {
		fractionOfTraffic := 1 / float64(opts.RateSampling)
		percentageTraffic := fractionOfTraffic * 100

		fmt.Printf("Sampling  %f percentage of traffic", percentageTraffic)

		tracerOptions = append(tracerOptions, traceSDK.WithSampler(traceSDK.TraceIDRatioBased(fractionOfTraffic)))
	}
	tracerProvider := traceSDK.NewTracerProvider(tracerOptions...)

	// Setup global tracer provider.
	otel.SetTracerProvider(tracerProvider)

	// Setup global text map propagator.
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	return newProvider(tracerProvider, opts.ServiceName)
}

func newProvider(provider *traceSDK.TracerProvider, serviceName string) tracer.Provider {
	return &jaegerImpl{provider: provider, serviceName: serviceName}
}

func (tr *jaegerImpl) Shutdown(ctx context.Context) error {
	return tr.provider.Shutdown(ctx)
}

func (tr *jaegerImpl) GetName() string {
	return "jaeger"
}

func (tr *jaegerImpl) GetServiceName() string {
	return tr.serviceName
}
