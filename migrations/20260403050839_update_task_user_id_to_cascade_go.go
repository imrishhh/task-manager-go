package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upUpdateTaskUserIDToCascadeGo, downUpdateTaskUserIDToCascadeGo)
}

func upUpdateTaskUserIDToCascadeGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx,
		`
		ALTER TABLE tasks
		DROP CONSTRAINT fk_user_id;

		ALTER TABLE tasks
		ADD CONSTRAINT fk_user
		FOREIGN KEY(user_id)
		REFERENCES users(id)
		ON DELETE CASCADE
		`,
	)
	return err
}

func downUpdateTaskUserIDToCascadeGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
