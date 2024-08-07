package ctx

import (
	"context"
)

type traceKey struct{}
type userKey struct{}
type permKey struct{}

func SetTraceID(ctx context.Context, in string) context.Context {
	return context.WithValue(ctx, traceKey{}, in)
}

func GetTraceID(ctx context.Context) string {
	res, _ := ctx.Value(traceKey{}).(string)
	return res
}

func SetUsername(ctx context.Context, in string) context.Context {
	return context.WithValue(ctx, userKey{}, in)
}

func GetUsername(ctx context.Context) string {
	res, _ := ctx.Value(userKey{}).(string)
	return res
}

func SetPermissions(ctx context.Context, in map[string][]string) context.Context {
	if in == nil {
		in = map[string][]string{}
	}
	return context.WithValue(ctx, permKey{}, in)
}

func GetPermissions(ctx context.Context) map[string][]string {
	res, _ := ctx.Value(permKey{}).(map[string][]string)
	return res
}
