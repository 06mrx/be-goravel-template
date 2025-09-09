package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250824082334CreateAuditsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250824082334CreateAuditsTable) Signature() string {
	return "20250824082334_create_audits_table"
}

// Up Run the migrations.
func (r *M20250824082334CreateAuditsTable) Up() error {
	if !facades.Schema().HasTable("audits") {
		return facades.Schema().Create("audits", func(table schema.Blueprint) {
			table.BigIncrements("id")
			table.String("auditable_type")
			table.String("auditable_id")
			table.String("event")
			table.String("user_type").Nullable()
			table.String("user_id").Nullable()
			table.Json("old_values").Nullable()
			table.Json("new_values").Nullable()
			// table.TimestampsTz()
			table.DateTime("created_at").Default("now()")
			table.DateTime("updated_at").Default("now()")
			table.Index("auditable_type", "auditable_id", "user_id")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250824082334CreateAuditsTable) Down() error {
	return facades.Schema().DropIfExists("audits")
}
