package tracer

import (
	"io"
	"log"
	"os"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func Create() (io.Closer, error) {
	tracer, closer, err := config.Configuration{
		ServiceName: "GraphQL",
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: getEnv("ADDR_TRACE_SERVER", "127.0.0.1:5775"),
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  getEnv("ADDR_TRACE_SERVER", "127.0.0.1:5775"),
		},
	}.NewTracer(config.Logger( /*jaeger.StdLogger*/ jaeger.NullLogger))

	opentracing.SetGlobalTracer(tracer)
	return closer, err
}

func getEnv(env string, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Printf("For %s applied default value: %s\n", env, defaultValue)
		return defaultValue
	}
	return value
}
