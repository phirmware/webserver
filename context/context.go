package context

import (
	"context"

	"lenslocked.com/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)

// WithUser is a middleware function that sets a current context with a user
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User looks up a user in the context
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
