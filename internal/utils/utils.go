package utils

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	//todo может быть надо заменить
	"go.opencensus.io/trace"
)

func InitJaeger(serviceName string) {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:6831",
		Process: jaeger.Process{
			ServiceName: serviceName,
			Tags: []jaeger.Tag{
				jaeger.StringTag("hostname", "localhost"),
			},
		},
	})
	if err != nil {
		return
	}
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.AlwaysSample(),
	})
}
