package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250823035147CreateRolesTable struct{}

// Signature The unique signature for the migration.
func (r *M20250823035147CreateRolesTable) Signature() string {
	return "20250823035147_create_roles_table"
}

// Up Run the migrations.
func (r *M20250823035147CreateRolesTable) Up() error {
	if err := facades.DB().Statement(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return err
	}

	if err := facades.DB().Statement(`
		CREATE TABLE IF NOT EXISTS roles (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
		);
	`); err != nil {
		return err
	}

	if err := facades.Schema().Table("roles", func(table schema.Blueprint) {
		table.String("name")
		table.TimestampsTz()

		table.Unique("name")
	}); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250823035147CreateRolesTable) Down() error {
	return facades.Schema().DropIfExists("roles")
}
