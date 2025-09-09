package middleware

import (
	"fmt"
	"goravel_api/app/utils"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/mssola/user_agent"
)

func User() http.Middleware {
	return func(ctx http.Context) {

		requestPath := ctx.Request().Url()
		// Ambil header Referer dan User-Agent
		referer := ctx.Request().Header("Referer")
		userAgentStr := ctx.Request().Header("User-Agent")
		appHost := facades.Config().Env("APP_HOST", "127").(string)
		appName := facades.Config().GetString("app.name")
		appEnv := facades.Config().Env("APP_ENV", "local").(string)

		if appEnv != "local" {
			// Periksa apakah Referer ada dan sesuai dengan domain yang diharapkan
			// Ganti "your-domain.com" dengan domain frontend Anda
			if requestPath != "/" && (referer == "" || !strings.Contains(referer, appHost)) {

				ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
					"message": "Akses ditolak.",
				})

			}

			// Periksa apakah User-Agent mengindikasikan bot atau aplikasi non-browser
			ua := user_agent.New(userAgentStr)
			if ua.Bot() || strings.Contains(userAgentStr, "Postman") || strings.Contains(userAgentStr, "curl") || strings.Contains(userAgentStr, "insomnia") || userAgentStr == "" {
				ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
					"message": "Akses ditolak!",
				})
			}
		}

		user_id, _ := utils.GetUserIDFromToken(ctx, ctx.Request().Header("Authorization"))

		ctx.WithValue("x-user-id", user_id)
		fmt.Println("global idleware " + appHost)
		fmt.Println("Nama Aplikasi:", appName)
		fmt.Println("Reerer:", referer)
		fmt.Println("URL:", requestPath)
		fmt.Println("ENV:", appEnv)

		ctx.Request().Next()
	}
}
