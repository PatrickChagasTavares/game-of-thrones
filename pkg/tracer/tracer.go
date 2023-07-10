package tracer

import (
	"context"
)

type (
	ITracer interface {
		Close()
		GetProviderName() string
	}

	Provider interface {
		Shutdown(ctx context.Context) error
		GetServiceName() string
		GetName() string
	}

	Tracer struct {
		provider Provider
	}
)
