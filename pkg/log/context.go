package log

import (
	"context"

	xctx "github.com/clyso/ceph-api/pkg/ctx"
	"github.com/rs/zerolog"
)

func WithTraceID(ctx context.Context, t string) context.Context {
	if t == "" {
		return ctx
	}
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(TraceID, t)
	})
	return xctx.SetTraceID(ctx, t)
}

func WithUsername(ctx context.Context, t string) context.Context {
	if t == "" {
		return ctx
	}
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(Username, t)
	})
	return xctx.SetUsername(ctx, t)
}
