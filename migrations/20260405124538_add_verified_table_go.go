package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddVerifiedTableGo, downAddVerifiedTableGo)
}

func upAddVerifiedTableGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE users ADD COLUMN verified BOOLEAN NOT NULL DEFAULT FALSE;
		`)
	return err
}

func downAddVerifiedTableGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE users DROP COLUMN verified;
		`)
	return err
}
