package controllers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel_api/app/models"
	"goravel_api/app/utils"
)

type RoleController struct {
}

// Index: Mengambil semua peran.
func (r *RoleController) Index(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt("page", 1)
	per_page := ctx.Request().QueryInt("per_page", 10)

	// 2. Get the search query from the query string
	searchQuery := ctx.Request().Query("search", "")

	// 3. Start a new query
	query := facades.Orm().Query()

	if searchQuery != "" {
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%")
	}
	var roles []models.Role
	var total int64

	err := query.With("Permissions").Paginate(page, per_page, &roles, &total)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil peran",
			"error":   err.Error(),
		})
	}

	// 3. Buat respons yang mencakup data, total, dan metadata paginasi
	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Peran berhasil diambil",
		"data": http.Json{
			"total":        total,
			"per_page":     per_page,
			"current_page": page,
			"last_page":    (total + int64(per_page) - 1) / int64(per_page),
			"data":         roles,
		},
	})

}

// Store: Menyimpan peran baru.
func (r *RoleController) Store(ctx http.Context) http.Response {
	user, _ := utils.GetUserFromToken(ctx, ctx.Request().Header("Authorization"))

	if facades.Gate().Denies("store-role", map[string]any{
		"user": user,
	}) {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized."})
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name": "required|string|min_len:3|unique:roles,name",
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

	role := models.Role{
		Name: ctx.Request().Input("name"),
	}

	if err := facades.Orm().Query().Create(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal membuat peran",
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "Peran berhasil dibuat",
		"data":    role,
	})
}

// Show: Mengambil satu peran berdasarkan ID.
func (r *RoleController) Show(ctx http.Context) http.Response {
	var role models.Role
	id := ctx.Request().Route("id")

	if err := facades.Orm().Query().With("Permissions").Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Peran tidak ditemukan",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Peran berhasil diambil",
		"data":    role,
	})
}

// Update: Memperbarui peran yang ada.
func (r *RoleController) Update(ctx http.Context) http.Response {
	var role models.Role
	id := ctx.Request().Route("id")

	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Peran tidak ditemukan",
		})
	}
	validator, err := ctx.Request().Validate(map[string]string{
		"name": "required|string|min_len:3|unique:roles,name," + id,
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

	role.Name = ctx.Request().Input("name")

	if err := facades.Orm().Query().Save(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal memperbarui peran",
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Peran berhasil diperbarui",
		"data":    role,
	})
}

// Destroy: Menghapus peran.
func (r *RoleController) Destroy(ctx http.Context) http.Response {
	var role models.Role
	id := ctx.Request().Route("id")

	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Peran tidak ditemukan",
		})
	}

	_, err := facades.Orm().Query().Where("id", id).Delete(&role)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menghapus peran",
			"error":   err.Error(),
		})
	}

	if err := facades.Orm().Query().Model(&role).Association("Permissions").Clear(); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal membersihkan peran yang ada",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Peran berhasil dihapus",
	})
}

// AttachPermissions melampirkan banyak izin ke satu peran, mengikuti pola yang Anda berikan.
func (r *RoleController) AttachPermissions(ctx http.Context) http.Response {
	// 1. Validasi input: Pastikan ada role_id (dari URL) dan permission_ids (dari body).
	validator, err := ctx.Request().Validate(map[string]string{
		"id": "required|uuid", // 💡 Ini adalah ID dari URL
		// "permission_ids": "required|array", // 💡 Ini adalah array dari body
	})
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{"message": validator.Errors().All()})
	}

	roleID := ctx.Request().Input("id")
	permissionIDsString := ctx.Request().InputArray("permission_ids")
	fmt.Println("1")
	fmt.Println(permissionIDsString)
	// Konversi array string ke array uint.
	var permissionIDs []any
	// var permissionIDStrings []any
	for _, idStr := range permissionIDsString {
		// 2. Parse string ke UUID (seperti yang sudah Anda lakukan).
		id, err := uuid.Parse(idStr)
		if err != nil {
			return ctx.Response().Json(http.StatusBadRequest, http.Json{"message": "ID izin tidak valid."})
		}
		permissionIDs = append(permissionIDs, id)
		fmt.Println("2")
		fmt.Println(permissionIDs)
		// 3. Konversi kembali UUID ke string dan tambahkan ke slice baru.
		// permissionIDStrings = append(permissionIDStrings, id.String())
	}

	// 2. Cari role dan permissions di database.
	var role models.Role
	if err := facades.Orm().Query().Where("id", roleID).First(&role); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Peran tidak ditemukan."})
	}

	var permissions []models.Permission
	// Gunakan WhereIn untuk mencari semua izin yang ID-nya ada dalam array.
	if err := facades.Orm().Query().WhereIn("id", permissionIDs).Find(&permissions); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Gagal menemukan beberapa izin."})
	}

	// 3. Gunakan Association().Replace() untuk mengganti relasi yang ada.
	// Ini akan menghapus relasi lama dan menambahkan yang baru.
	if err := facades.Orm().Query().Model(&role).Association("Permissions").Clear(); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal membersihkan peran yang ada",
			"error":   err.Error(),
		})
	}
	if err := facades.Orm().Query().Model(&role).Association("Permissions").Append(&permissions); err != nil {
		fmt.Printf("Gagal melampirkan izin ke peran: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal melampirkan izin.",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{"message": "Izin berhasil dilampirkan."})
}
