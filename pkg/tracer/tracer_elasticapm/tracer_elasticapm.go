package tracer_elastic

import (
	"context"
	"log"

	"go.elastic.co/apm/module/apmotel/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/PatrickChagastavares/game-of-thrones/pkg/tracer"
)

type (
	elasticIpml struct {
		provider    trace.TracerProvider
		serviceName string
	}
)

func NewExporter(serviceName string) tracer.Provider {
	provider, err := apmotel.NewTracerProvider()
	if err != nil {
		log.Fatal("we can't initialize elastic apm", err)
	}

	// Setup global tracer provider.
	otel.SetTracerProvider(provider)

	// Setup global text map propagator.
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	return newProvider(provider, serviceName)
}

func newProvider(provider trace.TracerProvider, serviceName string) tracer.Provider {
	return &elasticIpml{provider: provider, serviceName: serviceName}
}

func (tr *elasticIpml) Shutdown(ctx context.Context) error {
	return nil
}

func (tr *elasticIpml) GetName() string {
	return "elastic"
}

func (tr *elasticIpml) GetServiceName() string {
	return tr.serviceName
}
