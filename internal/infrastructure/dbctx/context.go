package dbctx

import (
	"context"

	"github.com/uptrace/bun"
)

type dbKey struct{}

// WithDB injects a bun.IDB (can be *bun.DB or bun.Tx) into the context.
func WithDB(ctx context.Context, db bun.IDB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

// FromCtx extracts bun.IDB from context, falling back to the default DB if not found.
func FromCtx(ctx context.Context, defaultDB *bun.DB) bun.IDB {
	if db, ok := ctx.Value(dbKey{}).(bun.IDB); ok {
		return db
	}
	return defaultDB
}
