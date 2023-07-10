package tracer

import (
	"context"
	"log"
)

func New(provicer Provider) ITracer {
	return Tracer{provider: provicer}
}

func (t Tracer) Close() {
	err := t.provider.Shutdown(context.Background())
	if err != nil {
		log.Println("failed to shutdown -> ", err)
	}
}

func (t Tracer) GetProviderName() string {
	return t.provider.GetName()
}
