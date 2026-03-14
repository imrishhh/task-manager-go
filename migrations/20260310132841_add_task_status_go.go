package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTaskStatusGo, downAddTaskStatusGo)
}

func upAddTaskStatusGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE tasks ADD COLUMN status VARCHAR(20) DEFAULT 'pending'
		CHECK (status IN ('pending', 'active', 'completed'));
	`)
	return err
}

func downAddTaskStatusGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE tasks DROP COLUMN status;
	`)
	return err
}
