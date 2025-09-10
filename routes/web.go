package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
	"github.com/mssola/user_agent"
)

func Web() {
	// facades.Route().Static("/", "./public")
	facades.Route().Static("public", "./storage/app").Name("storage")

	facades.Route().Get("/", func(ctx http.Context) http.Response {
		userAgentStr := ctx.Request().Header("User-Agent")

		// Buat parser baru
		ua := user_agent.New(userAgentStr)

		// Ambil informasi yang relevan
		ipAddress := ctx.Request().Ip()
		browserName, browserVersion := ua.Browser()
		osName := ua.OS()
		platform := ua.Platform()
		isBot := ua.Bot()
		// engine := ua.Engine(name, version expected) // Belum digunakan untuk verifikasi mesin
		osinfo := ua.OSInfo()
		mobile := ua.Mobile()
		model := ua.Model()
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version":  support.Version,
			"ip":       ipAddress,
			"ua":       userAgentStr,
			"platform": platform,
			"os":       osName,
			"browser":  browserName + " " + browserVersion,
			"is_bot":   isBot,
			"os_info":  osinfo,
			"mobile":   mobile,
			"model":    model,
		})
	})
}
