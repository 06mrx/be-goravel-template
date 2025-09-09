package controllers

import (
	"goravel_api/app/models"
	"goravel_api/app/utils"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

// Index menampilkan semua data users dari database
func (r *UserController) Index(ctx http.Context) http.Response {

	page := ctx.Request().QueryInt("page", 1)
	per_page := ctx.Request().QueryInt("per_page", 10)

	// 2. Get the search query from the query string
	searchQuery := ctx.Request().Query("search", "")

	// 3. Start a new query
	query := facades.Orm().Query()
	if searchQuery != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%").OrWhere("email ILIKE ?", "%"+searchQuery+"%")
	}
	var users []models.User
	var total int64

	err := query.Paginate(page, per_page, &users, &total)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil data pengguna",
			"error":   err.Error(),
		})
	}

	// Query semua data users
	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Data pengguna berhasil diambil",
		"data": http.Json{
			"total":        total,
			"per_page":     per_page,
			"current_page": page,
			"last_page":    (total + int64(per_page) - 1) / int64(per_page),
			"data":         users,
		},
	})
}

// Show menampilkan detail user berdasarkan ID
func (r *UserController) Show(ctx http.Context) http.Response {
	var user models.User
	userID := ctx.Request().Input("id")
	if userID == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "User ID tidak ada",
		})
	}

	// Query user berdasarkan ID
	if err := facades.Orm().Query().Where("id", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Pengguna tidak ditemukan",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "Data user berhasil diambil",
		"data":    user,
	})
}

// Store membuat user baru
// Store membuat user baru dengan Request Validation
func (r *UserController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name":                  "required|string|min_len:3",
		"email":                 "required|email|unique:users,email",
		"password":              "required|string|min_len:8",
		"password_confirmation": "required|string|min_len:8|same:password",
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

	hashedPassword, err := facades.Hash().Make(password)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Hashing password gagal",
			"error":   err.Error(),
		})
	}

	userId, _ := utils.GetUserIDFromToken(ctx, ctx.Request().Header("Authorization"))

	// Set data user dari input
	user.Name = name
	user.Email = email
	user.Password = hashedPassword
	user.CreatedBy = userId
	user.UpdatedBy = userId

	// Simpan user ke database
	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menyimpan pengguna",
			"error":   err.Error(),
		})
	}

	// Remove password dari response
	user.Password = ""

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "Pengguna berhasil disimpan",
		"data":    user,
	})
}

// Update memperbarui data user
func (r *UserController) Update(ctx http.Context) http.Response {
	var user models.User
	userID := ctx.Request().Input("id")

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

	// Ambil ID dari parameter route

	if userID == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "User ID tidak ada",
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
	userId, _ := utils.GetUserIDFromToken(ctx, ctx.Request().Header("Authorization"))
	user.UpdatedBy = userId
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

// Destroy menghapus user
func (r *UserController) Destroy(ctx http.Context) http.Response {
	// Ambil ID dari parameter route
	userID := ctx.Request().Input("id")
	if userID == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "User ID tidak ada",
		})
	}

	// Cek apakah user exists
	var user models.User
	if err := facades.Orm().Query().Where("id", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "User tidak ditemukan",
			"error":   err.Error(),
		})
	}

	// Hapus user
	result, err := facades.Orm().Query().Where("id", userID).Delete(&models.User{})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Pengguna gagal dihapus",
			"error":   err.Error(),
		})
	}

	// Cek apakah ada row yang terhapus
	if result.RowsAffected == 0 {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Tidak ada pengguna yang dihapus",
		})
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "Berhasil menghapus pengguna",
	})
}

// AssignRole: Menetapkan peran pada pengguna tertentu.
func (r *UserController) AssignRole(ctx http.Context) http.Response {
	// 1. Validasi input: Pastikan ada user_id dan role_id atau role_name.
	validator, err := ctx.Request().Validate(map[string]string{
		"user_id": "required|string",
		"role_id": "required|string",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
			"message": validator.Errors().All(),
		})
	}

	userID := ctx.Request().Input("user_id")
	roleID := ctx.Request().Input("role_id")
	// 2. Cari user dan role di database.
	var user models.User
	if err := facades.Orm().Query().Where("id", userID).First(&user); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "User not found",
		})
	}

	var role models.Role
	if err := facades.Orm().Query().Where("id", roleID).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Role not found",
		})
	}
	// fmt.Println("Role found:", role)
	// fmt.Println("User found:", user)

	// 3. Tambahkan peran ke user menggunakan relasi many-to-many.
	// Association().Append() akan menambahkan entri ke tabel pivot (role_user)
	if err := facades.Orm().Query().Model(&user).Association("Roles").Append(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menetapkan peran",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Peran berhasil ditetapkan",
		"data":    user,
	})
}

func (r *UserController) ReplaceRole(ctx http.Context) http.Response {
	// 1. Validasi input: Pastikan ada user_id dan role_id.
	validator, err := ctx.Request().Validate(map[string]string{
		"user_id": "required|string",
		"role_id": "required|string",
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{"message": validator.Errors().All()})
	}

	userID := ctx.Request().Input("user_id")
	roleID := ctx.Request().Input("role_id")

	// 2. Cari user dan role di database.
	var user models.User
	if err := facades.Orm().Query().Where("id", userID).First(&user); err != nil {
		// 💡 Perbaiki di sini: Tangani kasus 'user not found'
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "User not found."})
	}

	var newRole models.Role
	if err := facades.Orm().Query().Where("id", roleID).First(&newRole); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Role not found"})
	}

	// 3. Ganti peran user menggunakan relasi many-to-many.
	// Replace() akan menghapus semua peran yang ada dan menggantinya dengan peran yang baru.
	if err := facades.Orm().Query().Model(&user).Association("Roles").Replace(&newRole); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengganti peran",
			"error":   err.Error(),
		})
	}

	// GORM will now remove any old roles from the pivot table and add the new one.
	return ctx.Response().Json(http.StatusOK, http.Json{"message": "Peran berhasil diganti."})
}
