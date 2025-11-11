package grpc

import "context"
import "github.com/axmdl1/MedicalDataExchange/user-service/internal/service"

type ctxKey int

const claimsKey ctxKey = 1

func withClaims(ctx context.Context, c *service.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}
func claimsFromCtx(ctx context.Context) (*service.Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*service.Claims)
	return c, ok
}
