package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/mssola/user_agent"
)

type HelloWorldController struct {
	// Dependent services
}

func NewHelloWorldController() *HelloWorldController {
	return &HelloWorldController{
		// Inject services
	}
}

func (r *HelloWorldController) Index(ctx http.Context) http.Response {
	// return json helo world
	//get user_id from context in middleware
	user_id := ctx.Value("x-user-id")
	// fmt.Println("User ID from context:", user_id)
	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Hello World from Goravel 🚀",
		"user_id": user_id,
	})

}

func (c *HelloWorldController) WhoAmI(ctx http.Context) http.Response {
	// Ambil string User-Agent dari header permintaan
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
	uaa := ua.UA() // Dapatkan string User-Agent yang sudah diparse
	// ua.
	// Buat respons dalam bentuk map JSON
	response := http.Json{
		"ip":          ipAddress,
		"user_agent":  userAgentStr,
		"platform":    platform,
		"os":          osName,
		"browser":     browserName + " " + browserVersion,
		"is_bot":      isBot,
		"os_info":     osinfo,
		"mobile":      mobile,
		"model":       model,
		"ua_parsed":   uaa,
		"message":     "Hello from WhoAmI endpoint!",
		"description": "This endpoint provides information about your request.",
	}

	return ctx.Response().Json(http.StatusOK, response)
}
