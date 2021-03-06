package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Using constants in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users table",
		Script: `
		CREATE TABLE IF NOT EXISTS users (
			user_id UUID PRIMARY KEY,
			email TEXT NOT NULL,
			user_name TEXT NOT NULL,
			avatar TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NULL,
			deleted_at TIMESTAMP DEFAULT NULL
		)`,
	},
	{
		Version:     2,
		Description: "Add Password hash column.",
		Script:      `ALTER TABLE users ADD Column password_hash TEXT;`,
	},
	{
		Version:     3,
		Description: "Add unique index on email.",
		Script:      `CREATE UNIQUE INDEX email_idx ON users(email);`,
	},
	{
		Version:     4,
		Description: "ADD UNIQ INDEX on user_name",
		Script:      "CREATE UNIQUE INDEX user_name_idx ON users(user_name);",
	},
	{
		Version:     5,
		Description: "Create function trigger which returns time.NOW()",
		Script: `
		CREATE OR REPLACE FUNCTION updated_at_refresh() 
				RETURNS TRIGGER AS $$ 
			BEGIN NEW.updated_at = NOW(); 
			RETURN NEW; 
			END;
		$$ LANGUAGE 'plpgsql'`,
	},
	{
		Version:     6,
		Description: "Apply trigger on update operations on users table",
		Script:      `
		CREATE TRIGGER users_updated_at BEFORE UPDATE ON users 
		FOR EACH ROW EXECUTE PROCEDURE updated_at_refresh()`,
	},
}
