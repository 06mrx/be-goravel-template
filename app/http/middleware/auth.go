package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Auth() http.Middleware {
	return func(ctx http.Context) {
		// 1. Dapatkan token dari header Authorization
		token := ctx.Request().Header("Authorization")
		if token == "" {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": "Unauthenticated",
			})
		}

		// 2. Verifikasi token menggunakan facades.Auth()
		if _, err := facades.Auth(ctx).Parse(token); err != nil {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": "Unauthenticated",
			})
		}

		// 3. Pastikan user login valid
		user := facades.Auth(ctx).User(token)
		if user == nil {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"message": "Unauthenticated",
			})
		}

		// 4. Jika semua valid, lanjutkan ke handler berikutnya
		ctx.Request().Next()
	}
}
