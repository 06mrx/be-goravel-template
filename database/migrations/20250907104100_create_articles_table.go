package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250907104100CreateArticlesTable struct{}

// Signature The unique signature for the migration.
func (r *M20250907104100CreateArticlesTable) Signature() string {
	return "20250907104100_create_articles_table"
}

// Up Run the migrations.
func (r *M20250907104100CreateArticlesTable) Up() error {
	if err := facades.DB().Statement(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return err
	}

	if err := facades.DB().Statement(`
        CREATE TABLE IF NOT EXISTS articles (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
        );
    `); err != nil {
		return err
	}

	if err := facades.Schema().Table("articles", func(table schema.Blueprint) {
		table.String("title")
		table.String("slug")
		table.String("image_url").Nullable()
		table.Text("content")
		table.Uuid("user_id")
		table.String("status")
		table.Uuid("created_by")
		table.Uuid("updated_by").Nullable()
		table.TimestampsTz()
		table.SoftDeletesTz()

		table.Foreign("created_by").References("id").On("users")
		table.Foreign("updated_by").References("id").On("users").CascadeOnDelete()
	}); err != nil {
		return err
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250907104100CreateArticlesTable) Down() error {
	return facades.Schema().DropIfExists("articles")
}
