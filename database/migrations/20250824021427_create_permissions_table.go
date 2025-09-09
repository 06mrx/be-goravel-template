package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250824021427CreatePermissionsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250824021427CreatePermissionsTable) Signature() string {
	return "20250824021427_create_permissions_table"
}

// Up Run the migrations.
func (r *M20250824021427CreatePermissionsTable) Up() error {
	if err := facades.DB().Statement(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return err
	}

	if err := facades.DB().Statement(`
		CREATE TABLE IF NOT EXISTS permissions (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
		);
	`); err != nil {
		return err
	}
	if err := facades.Schema().Table("permissions", func(table schema.Blueprint) {
		table.String("name")
		table.String("created_by").Nullable()
		table.String("updated_by").Nullable()
		table.TimestampsTz()
		table.SoftDeletesTz()
	}); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250824021427CreatePermissionsTable) Down() error {
	return facades.Schema().DropIfExists("permissions")
}
