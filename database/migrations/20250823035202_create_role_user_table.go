package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250823035202CreateRoleUserTable struct{}

// Signature The unique signature for the migration.
func (r *M20250823035202CreateRoleUserTable) Signature() string {
	return "20250823035202_create_role_user_table"
}

// Up Run the migrations.
func (r *M20250823035202CreateRoleUserTable) Up() error {
	if err := facades.DB().Statement(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return err
	}

	if err := facades.DB().Statement(`
		CREATE TABLE IF NOT EXISTS role_user (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
		);
	`); err != nil {
		return err
	}
	if err := facades.Schema().Table("role_user", func(table schema.Blueprint) {
		table.Uuid("user_id")
		table.Uuid("role_id")
		table.Foreign("user_id").References("id").On("users").CascadeOnDelete()

		table.Foreign("role_id").References("id").On("roles").CascadeOnDelete()
		table.TimestampsTz()

		table.Index("user_id")
		table.Index("role_id")
	}); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250823035202CreateRoleUserTable) Down() error {
	return facades.Schema().DropIfExists("role_user")
}
