package tracer

import (
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func Create() (io.Closer, error) {
	tracer, closer, err := config.Configuration{
		ServiceName: "Mail",
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: "127.0.0.1:5775",
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  "127.0.0.1:5775",
		},
	}.NewTracer(config.Logger(jaeger.StdLogger))

	opentracing.SetGlobalTracer(tracer)
	return closer, err
}
