// app/utils/auth_helpers.go
package utils

import (
	"fmt"
	"goravel_api/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// GetUserIDFromToken mengurai dan memvalidasi token untuk mendapatkan ID pengguna.
func GetUserIDFromToken(ctx http.Context, token string) (string, error) {
	_, err := facades.Auth(ctx).Parse(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}
	userID, err := facades.Auth(ctx).ID()
	if err != nil {
		return "", fmt.Errorf("failed to get user ID: %w", err)
	}
	return userID, nil
}

// get user from token
func GetUserFromToken(ctx http.Context, token string) (*models.User, error) {
	// Parse token
	_, err := facades.Auth(ctx).Parse(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Konversi manual: ambil ID dari authUser, lalu cari di DB
	var user models.User
	userId, _ := GetUserIDFromToken(ctx, token)
	if err := facades.Orm().Query().Where("id", userId).First(&user); err != nil {
		return nil, fmt.Errorf("failed to load user model: %w", err)
	}

	return &user, nil
}

// cek otorisasi / permission
func CheckPermission(ctx http.Context, permissionName string) bool {
	token := ctx.Request().Header("Authorization")
	user, err := GetUserFromToken(ctx, token)
	if err != nil {
		// fmt.Printf("Error getting user from token: %v\n", err)
		return false
	}

	if facades.Gate().Denies(permissionName, map[string]any{
		"user": user,
	}) {
		// fmt.Printf("User %s does not have permission %s\n", user.Email, permissionName)
		return false
	}
	return true

}
