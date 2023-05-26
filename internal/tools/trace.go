package tools

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// initJaeger 将jaeger tracer设置为全局tracer
func InitJaeger(jaegerURL string, serverName string) error {
	// 创建 Jaeger exporter
	fmt.Println(jaegerURL)
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(exp),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serverName),
			attribute.String("exporter", "jaeger"),
			attribute.Float64("float", 312.23),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
	//cfg := jaegercfg.Configuration{
	//	// 将采样频率设置为1，每一个span都记录，方便查看测试结果
	//	Sampler: &jaegercfg.SamplerConfig{
	//		Type:  jaeger.SamplerTypeConst,
	//		Param: 1,
	//	},
	//	Reporter: &jaegercfg.ReporterConfig{
	//		LogSpans: true,
	//		// 将span发往jaeger-collector的服务地址
	//		//CollectorEndpoint: "http://localhost:14268/api/traces",
	//		CollectorEndpoint: c.GetJaegerUrl(),
	//	},
	//}
	//closer, err := cfg.InitGlobalTracer(c.Server.GetName, jaegercfg.Logger(jaeger.StdLogger))
	//if err != nil {
	//	panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	//}
	//return closer
}
