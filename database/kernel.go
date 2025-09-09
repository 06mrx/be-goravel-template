package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"

	"goravel_api/database/migrations"
	"goravel_api/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateUsersTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20250823035147CreateRolesTable{},
		&migrations.M20250823035202CreateRoleUserTable{},
		&migrations.M20250824021427CreatePermissionsTable{},
		&migrations.M20250824021941CreatePermissionRoleTable{},
		&migrations.M20250824082334CreateAuditsTable{},
		&migrations.M20250828034249AddGoogleFieldsToUsersTable{},
		&migrations.M20250907104100CreateArticlesTable{},
	}
}

func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
