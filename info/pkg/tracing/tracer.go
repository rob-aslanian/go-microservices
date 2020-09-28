package tracing

import (
	"context"
	"io"
	"runtime/debug"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// Settings ...
type Settings struct {
	Address     string
	ServiceName string
}

// Tracer ...
type Tracer struct {
	tracer opentracing.Tracer
}

// NewTracer creates new instance of Tracer
func NewTracer(s Settings) (*Tracer, io.Closer, error) {
	tracer, closer, err := config.Configuration{
		ServiceName: s.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: s.Address,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  s.Address,
		},
	}.NewTracer(config.Logger( /*jaeger.StdLogger*/ jaeger.NullLogger))

	opentracing.SetGlobalTracer(tracer)
	return &Tracer{tracer}, closer, err
}

// MakeSpan returns span
func (t *Tracer) MakeSpan(ctx context.Context, spanName string) opentracing.Span {
	var span opentracing.Span
	if opentracing.SpanFromContext(ctx) == nil {
		span = t.tracer.StartSpan(spanName) // <trace-without-root-span>
	} else {
		span = opentracing.SpanFromContext(ctx)
	}

	return span.Tracer().StartSpan(
		spanName,
		opentracing.ChildOf(span.Context()),
	)
}

// GetTracer returns tracer
func (t *Tracer) GetTracer() opentracing.Tracer {
	return t.tracer
}

// LogError sents information about error to tracer server
func (t *Tracer) LogError(span opentracing.Span, err interface{}) {
	span.SetTag("error", true)
	span.LogKV("error.message", err)
	span.LogKV("error.stack", string(debug.Stack()))
}
