package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250824021941CreatePermissionRoleTable struct{}

// Signature The unique signature for the migration.
func (r *M20250824021941CreatePermissionRoleTable) Signature() string {
	return "20250824021941_create_permission_role_table"
}

// Up Run the migrations.
func (r *M20250824021941CreatePermissionRoleTable) Up() error {
	if err := facades.DB().Statement(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return err
	}

	if err := facades.DB().Statement(`
		CREATE TABLE IF NOT EXISTS permission_role (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
		);
	`); err != nil {
		return err
	}
	if err := facades.Schema().Table("permission_role", func(table schema.Blueprint) {
		table.Uuid("permission_id")
		table.Uuid("role_id")
		table.Timestamps()
		table.Foreign("permission_id").References("id").On("permissions").CascadeOnDelete()
		table.Foreign("role_id").References("id").On("roles").CascadeOnDelete()

		// table.TimestampsTz()
		table.SoftDeletesTz()
	}); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250824021941CreatePermissionRoleTable) Down() error {
	return facades.Schema().DropIfExists("permission_role")
}
