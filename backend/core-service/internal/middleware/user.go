package middleware

import "context"

type ctxUserKey struct{}

var userKey ctxUserKey

type UserInfo struct {
	Id       int64
	Email    string
	Role     string
	ClinicID *int64
}

func SetUser(ctx context.Context, u UserInfo) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func GetUserFromContext(ctx context.Context) (UserInfo, bool) {
	u, ok := ctx.Value(userKey).(UserInfo)
	return u, ok
}
