package tx

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"

	"minigo/internal/infrastructure/dbctx"
)

// Manager provides transaction management using Bun ORM
type Manager struct {
	DB *bun.DB
}

// NewManager creates a new transaction manager
func NewManager(db *bun.DB) *Manager {
	return &Manager{DB: db}
}

// InTx executes a function within a database transaction
func (m *Manager) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.DB.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// Inject the transaction into the context
		ctx = dbctx.WithDB(ctx, tx)
		return fn(ctx)
	})
}
