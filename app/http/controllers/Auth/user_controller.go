// app/http/controllers/user_controller.go
package controllers

import (
	"errors"
	"goravel_api/app/models"
	"goravel_api/app/utils"

	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type UserController struct {
	// You can declare dependencies here
}

// Register a new user
func (r *UserController) Register(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":                  "required|string|min_len:3",
		"email":                 "required|email|unique:users,email",
		"password":              "required|string|min:6",
		"password_confirmation": "required|string|min:6|same:password",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": validator.Errors().All(),
		})
	}

	var user models.User

	// Get input from request
	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	// Hash password
	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	// Set data user dari input
	user.Name = name
	user.Email = email
	user.Password = hashedPassword

	// Simpan user ke database
	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	// Remove password dari response
	user.Password = ""
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}
	var role models.Role
	if err := facades.Orm().Query().Where("name", "user").First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Role not found",
		})
	}

	if err := facades.Orm().Query().Model(&user).Association("Roles").Append(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menetapkan peran",
			"error":   err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "User registered successfully",
		"data":    user,
		"token":   token, // generate token here
		// "payload": payload,
	})
}

func (r *UserController) RegisterUsingGoogle(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":      "required|string|min_len:3",
		"email":     "required|email|unique:users,email",
		"google_id": "required",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": validator.Errors().All(),
		})
	}

	var user models.User

	// Get input from request
	name := ctx.Request().Input("name")
	email := ctx.Request().Input("email")
	password := "password"
	google_id := ctx.Request().Input(("google_id"))

	// Hash password
	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
	}

	// Set data user dari input
	user.Name = name
	user.Email = email
	user.Password = hashedPassword
	user.GoogleId = google_id
	user.AuthProvider = "google"

	// Simpan user ke database
	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to create user",
			"error":   err.Error(),
		})
	}

	// Remove password dari response
	user.Password = ""
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}
	var role models.Role
	if err := facades.Orm().Query().Where("name", "user").First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Role not found",
		})
	}

	if err := facades.Orm().Query().Model(&user).Association("Roles").Append(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menetapkan peran",
			"error":   err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "User registered successfully",
		"data":    user,
		"token":   token, // generate token here
		// "payload": payload,
	})
}

// Login user
func (r *UserController) Login(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"email":    "required|email",
		"password": "required|string",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": validator.Errors().All(),
		})
	}
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")
	if email == "" || password == "" {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": "Validation failed",
			"errors":  "Email and password are required",
		})
	}
	// Cari user berdasarkan email
	var user models.User
	if err := facades.Orm().Query().With("Roles").Where("email", email).First(&user); err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Invalid email or password",
		})
	}
	// Verifikasi password
	if !facades.Hash().Check(password, user.Password) {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Invalid email or password",
		})
	}
	// Hapus password dari response
	user.Password = ""

	// Generate token

	token, err := facades.Auth(ctx).Login(&user)

	// token, err := facades.Auth(ctx).CreateToken(&user, map[string]any{
	//     "username": user.Username,
	//     "email":    user.Email,
	//     "roles":    "admin", // Contoh data kustom
	// })

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}
	// payload, err := facades.Auth(ctx).Parse(token)
	_, err = facades.Auth(ctx).Parse(token)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to parse token",
			"error":   err.Error(),
		})
	}

	var role string
	if len(user.Roles) > 0 {
		role = user.Roles[0].Name // Asumsikan user.Roles adalah slice of structs dengan field 'Name'
		user.Roles = nil
	} else {
		role = "guest" // Jika tidak ada peran, tetapkan default
	}
	userData := map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"role":  role,
		"name":  user.Name,
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Login successful",
		"data": map[string]any{
			"user":  userData,
			"token": token,
			"role":  role,
		},
	})
}

func (c *UserController) VerifyAndLoginGoogle(ctx http.Context) http.Response {
	// Ambil token ID dari header Authorization (atau dari body)
	// tokenString := strings.Replace(ctx.Request().Header("Authorization"), "Bearer ", "", 1)
	// if tokenString == "" {
	// 	return ctx.Response().Json(http.StatusUnauthorized, http.Json{
	// 		"message": "Token tidak ditemukan.",
	// 	})
	// }

	google_id := ctx.Request().Input("idToken")

	if google_id == "" {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Token Google tidak ditemukan.",
		})
	}

	// Inisialisasi Firebase Admin SDK dengan file Service Account Key
	serviceAccountKeyPath := facades.Config().GetString("FIREBASE_ADMIN_SDK_PATH")
	opt := option.WithCredentialsFile(serviceAccountKeyPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menginisialisasi Firebase Admin SDK.",
		})
	}

	// Verifikasi Firebase ID Token
	client, err := app.Auth(ctx)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": err.Error(),
		})
	}

	token, err := client.VerifyIDToken(ctx, google_id)
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Firebase ID Token tidak valid.",
		})
	}

	// Cari pengguna berdasarkan UID Firebase
	var user models.User
	err = facades.Orm().Query().With("Roles").Where("email", token.Claims["email"].(string)).First(&user)

	if (err != nil && errors.Is(err, gorm.ErrRecordNotFound)) || user.ID == uuid.Nil {
		// Jika pengguna baru, buat akun baru
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Pengguna belum terdaftar",
			"data":    user,
		})
	} else if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Terjadi kesalahan saat mencari pengguna.",
		})
	}

	// Jika pengguna ditemukan, buat token API untuk login
	apiToken, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": err.Error(),
			"user":    user.ID == uuid.Nil,
		})
	}

	// return ctx.Response().Json(http.StatusOK, http.Json{
	// 	"message":      "Login dengan Google berhasil.",
	// 	"access_token": apiToken,
	// })
	var role string
	if len(user.Roles) > 0 {
		role = user.Roles[0].Name // Asumsikan user.Roles adalah slice of structs dengan field 'Name'
		user.Roles = nil
	} else {
		role = "guest" // Jika tidak ada peran, tetapkan default
	}
	userData := map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"role":  role,
		"name":  user.Name,
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Login successful",
		"data": map[string]any{
			"user":  userData,
			"token": apiToken,
			"role":  role,
		},
	})
}

func (r *UserController) RefreshToken(ctx http.Context) http.Response {
	token, err := facades.Auth(ctx).Refresh()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Refresh token successful",
		"token":   token,
	})
}

// show profile
func (r *UserController) Profile(ctx http.Context) http.Response {
	var user models.User
	id, err := facades.Auth(ctx).ID()

	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Unauthenticated",
		})
	}
	// ctx.Request().Session().Get("key")
	err = utils.FindModelByID(id, &user, "Roles")
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "User not found",
			"error":   err.Error(),
		})
	}
	user.Password = ""

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "User found",
		"data":    user,
	})
}

func (r *UserController) UpdateProfile(ctx http.Context) http.Response {
	var user models.User
	// userID := ctx.Request().Input("id")
	userID, _ := facades.Auth(ctx).ID()
	validator, err := ctx.Request().Validate(map[string]string{
		"name":                  "required|string|min_len:3",
		"email":                 "required|email|unique:users,email," + userID,
		"password":              "min_len:8",
		"password_confirmation": "min_len:8|same:password",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": validator.Errors().All(),
		})
	}

	// Cari user berdasarkan ID
	if err := facades.Orm().Query().Where("id", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Pengguna tidak ditemukan",
			"error":   err.Error(),
		})
	}

	// Update field yang diberikan
	if name := ctx.Request().Input("name"); name != "" {
		user.Name = name
	}
	if email := ctx.Request().Input("email"); email != "" {
		user.Email = email
	}
	user.UpdatedBy = userID
	if password := ctx.Request().Input("password"); password != "" {
		hashedPassword, err := facades.Hash().Make(password)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"message": "Hashing password gagal",
				"error":   err.Error(),
			})
		}
		user.Password = hashedPassword
	}

	// Simpan perubahan
	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal memperbarui pengguna",
			"error":   err.Error(),
		})
	}

	// Remove password dari response
	user.Password = ""

	return ctx.Response().Success().Json(http.Json{
		"message": "Pengguna berhasil diperbarui",
		"data":    user,
	})
}

// Logout user
func (r *UserController) Logout(ctx http.Context) http.Response {
	// Dapatkan token dari header Authorization
	token := ctx.Request().Header("Authorization")
	if token == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "Token is required",
		})
	}

	if _, err := facades.Auth(ctx).Parse(token); err != nil {
		// ... unauthorized response ...
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}
	// Logout user (invalidate token)
	if err := facades.Auth(ctx).Logout(); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to logout",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Logout successful",
	})
}
