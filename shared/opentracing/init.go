package opentracing

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func Create(name, addr string) (io.Closer, error) {
	tracer, closer, err := config.Configuration{
		ServiceName: name,
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: addr,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  addr,
		},
	}.NewTracer(config.Logger( /*jaeger.StdLogger*/ jaeger.NullLogger))

	opentracing.SetGlobalTracer(tracer)
	return closer, err
}
