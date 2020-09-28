package tracer

import (
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func Create(conf Configuration) (opentracing.Tracer, io.Closer, error) {
	tracer, closer, err := config.Configuration{
		ServiceName: conf.GetServiceName(),
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: conf.GetAddress(),
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  conf.GetAddress(),
		},
	}.NewTracer(config.Logger( /*jaeger.StdLogger*/ jaeger.NullLogger))

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, err
}

type Configuration interface {
	GetServiceName() string
	GetAddress() string
}
