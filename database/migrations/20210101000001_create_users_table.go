package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20210101000001CreateUsersTable struct{}

func (r *M20210101000001CreateUsersTable) Signature() string {
	return "20210101000001_create_users_table"
}

func (r *M20210101000001CreateUsersTable) Up() error {
	// 1️⃣ Buat tabel users
	if err := facades.DB().Statement(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		return err
	}

	if err := facades.DB().Statement(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4()
		);
	`); err != nil {
		return err
	}

	if err := facades.Schema().Table("users", func(table schema.Blueprint) {
		table.String("name")
		table.String("email")
		table.String("password")
		table.String("created_by").Nullable()
		table.String("updated_by").Nullable()
		table.TimestampsTz()
		table.SoftDeletesTz()
	}); err != nil {
		return err
	}

	return nil

}

func (r *M20210101000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
