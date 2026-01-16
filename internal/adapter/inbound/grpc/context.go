package grpc

import (
	"context"
	"errors"

	"github.com/0xsj/overwatch-pkg/types"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userDIDKey   contextKey = "user_did"
	tenantIDKey  contextKey = "tenant_id"
	sessionIDKey contextKey = "session_id"
)

var (
	ErrNoUserIDInContext    = errors.New("no user_id in context")
	ErrNoUserDIDInContext   = errors.New("no user_did in context")
	ErrNoTenantIDInContext  = errors.New("no tenant_id in context")
	ErrNoSessionIDInContext = errors.New("no session_id in context")
)

func WithUserID(ctx context.Context, userID types.ID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func WithUserDID(ctx context.Context, did string) context.Context {
	return context.WithValue(ctx, userDIDKey, did)
}

func WithTenantID(ctx context.Context, tenantID types.ID) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func WithSessionID(ctx context.Context, sessionID types.ID) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

func getUserIDFromContext(ctx context.Context) (types.ID, error) {
	val := ctx.Value(userIDKey)
	if val == nil {
		return "", ErrNoUserIDInContext
	}
	userID, ok := val.(types.ID)
	if !ok {
		return "", ErrNoUserIDInContext
	}
	return userID, nil
}

func getUserDIDFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(userDIDKey)
	if val == nil {
		return "", ErrNoUserDIDInContext
	}
	did, ok := val.(string)
	if !ok {
		return "", ErrNoUserDIDInContext
	}
	return did, nil
}

func getTenantIDFromContext(ctx context.Context) (types.ID, error) {
	val := ctx.Value(tenantIDKey)
	if val == nil {
		return "", ErrNoTenantIDInContext
	}
	tenantID, ok := val.(types.ID)
	if !ok {
		return "", ErrNoTenantIDInContext
	}
	return tenantID, nil
}

func getOptionalTenantIDFromContext(ctx context.Context) types.Optional[types.ID] {
	tenantID, err := getTenantIDFromContext(ctx)
	if err != nil {
		return types.None[types.ID]()
	}
	return types.Some(tenantID)
}

func getSessionIDFromContext(ctx context.Context) (types.ID, error) {
	val := ctx.Value(sessionIDKey)
	if val == nil {
		return "", ErrNoSessionIDInContext
	}
	sessionID, ok := val.(types.ID)
	if !ok {
		return "", ErrNoSessionIDInContext
	}
	return sessionID, nil
}

func GetUserIDFromContext(ctx context.Context) (types.ID, error) {
	return getUserIDFromContext(ctx)
}

func GetUserDIDFromContext(ctx context.Context) (string, error) {
	return getUserDIDFromContext(ctx)
}

func GetTenantIDFromContext(ctx context.Context) (types.ID, error) {
	return getTenantIDFromContext(ctx)
}

func GetSessionIDFromContext(ctx context.Context) (types.ID, error) {
	return getSessionIDFromContext(ctx)
}
