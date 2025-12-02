package main

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/google/uuid"
)

// TracingCallback 请求追踪
type TracingCallback struct {
	callbacks.HandlerBuilder
}

func (cb *TracingCallback) OnStart(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	// 生成或获取 trace ID
	traceID, ok := ctx.Value("trace_id").(string)
	if !ok {
		traceID = uuid.New().String()
		ctx = context.WithValue(ctx, "trace_id", traceID)
	}

	// 生成 span ID
	spanID := uuid.New().String()
	ctx = context.WithValue(ctx, "span_id_"+info.Name, spanID)

	fmt.Printf("[Trace:%s] [Span:%s] 开始 %s\\n",
		traceID[:8], spanID[:8], info.Name)

	// 记录开始时间
	ctx = context.WithValue(ctx, "start_"+spanID, time.Now())

	return ctx
}

func (cb *TracingCallback) OnEnd(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	traceID, _ := ctx.Value("trace_id").(string)
	spanID, _ := ctx.Value("span_id_" + info.Name).(string)
	startTime, _ := ctx.Value("start_" + spanID).(time.Time)

	duration := time.Since(startTime)

	fmt.Printf("[Trace:%s] [Span:%s] 完成 %s (耗时: %v)\\n",
		traceID[:8], spanID[:8], info.Name, duration)

	return ctx
}

func (cb *TracingCallback) OnError(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
	traceID, _ := ctx.Value("trace_id").(string)
	spanID, _ := ctx.Value("span_id_" + info.Name).(string)

	fmt.Printf("[Trace:%s] [Span:%s] 错误 %s: %v\\n",
		traceID[:8], spanID[:8], info.Name, err)

	return ctx
}
