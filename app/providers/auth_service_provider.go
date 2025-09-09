package providers

import (
	"context"
	"goravel_api/app/models"

	"github.com/goravel/framework/auth/access"
	contractsAccess "github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

type AuthServiceProvider struct{}

func (receiver *AuthServiceProvider) Register(app foundation.Application) {}

func (receiver *AuthServiceProvider) Boot(app foundation.Application) {
	// fmt.Print("Booting AuthServiceProvider...\n")
	var permissions []models.Permission
	if err := facades.Orm().Query().Find(&permissions); err != nil {
		facades.Log().Errorf("Gagal memuat izin dari database: %v", err)
		return
	}

	// facades.Log().Infof("Memuat %d izin dari database", len(permissions))
	// fmt.Printf("Memuat %d izin dari database\n", len(permissions))

	for _, permission := range permissions {
		p := permission.Name
		facades.Gate().Define(p, func(ctx context.Context, arguments map[string]any) contractsAccess.Response {
			// fmt.Printf("Arguments %s:\n", arguments)
			user, ok := arguments["user"].(*models.User)

			if !ok || user == nil {
				// fmt.Print("User not found in context or type assertion failed\n")
				return access.NewDenyResponse("User not found")
			}

			if user.HasPermissionTo(p) {
				// fmt.Printf("User %s memiliki izin %s\n", user.Email, p)
				return access.NewAllowResponse()
			}
			// fmt.Printf("User %s tidak memiliki izin %s\n", user.Email, p)
			return access.NewDenyResponse("Unauthorized")
		})
	}
}
