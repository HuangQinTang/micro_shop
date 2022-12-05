package common

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"time"
)

// NewTracer 创建链路追踪实例
// https://www.jaegertracing.io/docs/1.22/sampling/
// https://www.ssfiction.com/pythons/1069811.html
func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName, //服务名
		// 采样器，以下设置表示全部采集
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		// 配置客户端上报追踪信息
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: 1 * time.Second, //强制清空、推送队列时间，对于流量不高的程序，队列可能长时间不能满，那么设置这个时间，超时可以自动推送一次。对于高并发的情况，一般队列很快就会满的，满了后也会自动推送。默认为1秒。
			LogSpans:            true,            //日志推送
			LocalAgentHostPort:  addr,            //推送地址
		},
	}
	return cfg.NewTracer()
}

// WithTrace 返回trace_id
func WithTrace(ctx context.Context) string {
	var jTraceId jaeger.TraceID
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			mySpan := tracer.StartSpan("my info", opentracing.ChildOf(parentCtx))
			// 提取出一个jaeger的traceid
			if sc, ok := mySpan.Context().(jaeger.SpanContext); ok {
				jTraceId = sc.TraceID()
			}
			defer mySpan.Finish()
		}
	}
	return fmt.Sprint(jTraceId)
}
