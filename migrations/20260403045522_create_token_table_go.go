package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTokenTableGo, downCreateTokenTableGo)
}

func upCreateTokenTableGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.ExecContext(ctx,
		`
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	
		CREATE TYPE token_type AS ENUM ('bearer', 'refresh', 'reset', 'verify');

		CREATE TABLE user_tokens (
			id SERIAL PRIMARY KEY,
			user_id UUID NOT NULL,
			token TEXT NOT NULL,
			token_type token_type NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			revoked BOOLEAN NOT NULL DEFAULT FALSE,
			
			CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		`)
	return err
}

func downCreateTokenTableGo(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.ExecContext(ctx,
		`
		DROP TABLE IF EXISTS user_tokens;
		DROP TYPE IF EXISTS token_type;
		`,
	)
	return err
}
