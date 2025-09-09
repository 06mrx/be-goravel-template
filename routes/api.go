package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel_api/app/http/controllers"
	AuthController "goravel_api/app/http/controllers/Auth"

	"goravel_api/app/http/middleware"

	htppMiddleware "github.com/goravel/framework/http/middleware"
)

func Api() {
	userController := controllers.NewUserController()
	registercontroller := AuthController.UserController{}
	roleController := controllers.RoleController{}
	permissionController := controllers.PermissionController{}
	articleController := controllers.ArticleController{}
	auditController := controllers.AuditController{}

	// facades.Route().Post("api/register", registercontroller.Register)
	// facades.Route().Post("api/login", registercontroller.Login)
	// facades.Route().Post("api/logout", registercontroller.Logout)
	//route group with prefix
	facades.Route().Prefix("api").Middleware(htppMiddleware.Throttle("global")).Group(func(router route.Router) {
		//Auth
		router.Post("register", registercontroller.Register)
		router.Post("register-using-google", registercontroller.RegisterUsingGoogle)
		router.Post("login", registercontroller.Login)
		router.Post("login-using-google", registercontroller.VerifyAndLoginGoogle)
		router.Post("logout", registercontroller.Logout)
		router.Post("refresh-token", registercontroller.RefreshToken)
		router.Get("profile", registercontroller.Profile)
		router.Middleware(middleware.Auth()).Put("update-profile", registercontroller.UpdateProfile)

		router.Middleware(middleware.Auth()).Get("users", userController.Index)
		router.Middleware(middleware.Auth()).Get("users/{id}", userController.Show)
		router.Middleware(middleware.Auth()).Post("users", userController.Store)
		router.Middleware(middleware.Auth()).Put("users/{id}", userController.Update)
		router.Middleware(middleware.Auth()).Delete("users/{id}", userController.Destroy)

		//Role
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("roles", roleController.Index)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Post("roles", roleController.Store)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("roles/{id}", roleController.Show)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Put("roles/{id}", roleController.Update)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Delete("roles/{id}", roleController.Destroy)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Post("users/assign-role", userController.AssignRole)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Post("users/replace-role", userController.ReplaceRole)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Post("users/attach-permission", roleController.AttachPermissions)

		//Permission
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("permissions", permissionController.Index)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Post("permissions", permissionController.Store)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("permissions/{id}", permissionController.Show)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Put("permissions/{id}", permissionController.Update)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Delete("permissions/{id}", permissionController.Destroy)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("list-permissions", permissionController.List)
		//Artikel
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("articles", articleController.Index)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Post("articles", articleController.Store)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("articles/{id}", articleController.Show)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Put("articles/{id}", articleController.Update)
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Delete("articles/{id}", articleController.Destroy)
		//Audit
		router.Middleware(middleware.Auth()).Middleware(middleware.Role("admin")).Get("audits", auditController.Index)
	})

	// facades.Route().Prefix("api").Get("users", userController.Index)
	// facades.Route().Get("/users/{id}", userController.Show)
	facades.Route().Prefix("api").Middleware(middleware.Role("admin")).Middleware(middleware.Auth()).Get("hello-world", controllers.NewHelloWorldController().Index)
	facades.Route().Prefix("api").Get("open-hello-world", controllers.NewHelloWorldController().Index)
	facades.Route().Prefix("api").Get("whoami", controllers.NewHelloWorldController().WhoAmI)
	// facades.Route().Prefix("users").Get("/", userController.Show)

}
