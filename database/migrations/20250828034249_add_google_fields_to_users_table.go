package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250828034249AddGoogleFieldsToUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20250828034249AddGoogleFieldsToUsersTable) Signature() string {
	return "20250828034249_add_google_fields_to_users_table"
}

// Up Run the migrations.
func (r *M20250828034249AddGoogleFieldsToUsersTable) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.Text("google_id").Nullable()
		table.String("auth_provider").Nullable()
	})
}

// Down Reverse the migrations.
func (r *M20250828034249AddGoogleFieldsToUsersTable) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		table.DropColumn("google_id")
		table.DropColumn("auth_provider")

	})
}
