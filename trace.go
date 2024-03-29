package main

import (
	"context"
	"fmt"

	"go.opencensus.io/trace"
)

func startSpan(ctx context.Context, name string) (context.Context, *trace.Span) {
	return trace.StartSpan(ctx, fmt.Sprintf("/storage_runner/go/%s", name))
}
