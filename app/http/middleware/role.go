// app/http/middlewares/RoleMiddleware.go

package middleware

import (
	"strings"

	"goravel_api/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Role(rolesParam string) http.Middleware {
	return func(ctx http.Context) {
		// rolesParam := ctx.Request().Route("role")
		requiredRoles := strings.Split(rolesParam, ",")

		// fmt.Println("Required roles for this route:", requiredRoles)

		token := ctx.Request().Header("Authorization")
		_, err := facades.Auth(ctx).Parse(token)
		if err != nil {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."})
			return
		}

		user := facades.Auth(ctx).User(token)
		if user == nil {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."})
			return
		}
		userID, err := facades.Auth(ctx).ID()
		if err != nil {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{"message": "Unauthenticated."})
			return
		}
		// fmt.Println("Authenticated user ID:", userID)
		// Query user by ID
		var userWithRoles models.User
		if err := facades.Orm().Query().With("Roles").Where("id", userID).First(&userWithRoles); err != nil {
			facades.Log().Errorf("Failed to retrieve user with roles for ID %s: %v", userID, err)
			ctx.Request().AbortWithStatusJson(http.StatusInternalServerError, http.Json{"message": "Failed to load user information."})
			return
		}

		// 4. Periksa apakah pengguna memiliki salah satu peran yang diperlukan.
		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range userWithRoles.Roles {
				if userRole.Name == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		// 5. Jika tidak ada peran yang cocok, tolak akses.
		if !hasRequiredRole {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{"message": "Unauthorized."})
			return
		}

		// 6. Jika semua valid, lanjutkan ke handler berikutnya.
		ctx.Request().Next()
	}
}
