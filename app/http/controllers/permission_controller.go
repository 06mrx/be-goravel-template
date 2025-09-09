package controllers

import (
	"goravel_api/app/models"
	"goravel_api/app/utils"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type PermissionController struct {
	// Dependent services
}

// Index: Mengambil semua permission.
func (r *PermissionController) Index(ctx http.Context) http.Response {
	// Check for authorization
	// allowed := utils.CheckPermission(ctx, "index-permissions")
	// if !allowed {
	// 	return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized"})
	// }

	// 1. Get pagination parameters from the query string
	page := ctx.Request().QueryInt("page", 1)
	per_page := ctx.Request().QueryInt("per_page", 10)

	// 2. Get the search query from the query string
	searchQuery := ctx.Request().Query("search", "")

	// 3. Start a new query
	query := facades.Orm().Query()

	// 4. Add a WHERE clause if a search query is provided
	if searchQuery != "" {
		// We use `ILIKE` for case-insensitive searching in PostgreSQL,
		// and it will be converted to `LIKE` for other databases.
		// The `%` are wildcards to search for the term anywhere in the name.
		query = query.Where("name ILIKE ?", "%"+searchQuery+"%")
	}

	var permissions []models.Permission
	var total int64

	// 5. Paginate the filtered query
	err := query.Paginate(page, per_page, &permissions, &total)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Error fetching permissions",
			"error":   err.Error(),
		})
	}

	// 6. Build and return the JSON response with the paginated data
	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Permission berhasil diambil",
		"data": http.Json{
			"total":        total,
			"per_page":     per_page,
			"current_page": page,
			"last_page":    (total + int64(per_page) - 1) / int64(per_page),
			"data":         permissions,
		},
	})
}

// fungsi list permission tanpa paginasi
func (r *PermissionController) List(ctx http.Context) http.Response {

	var permissions []models.Permission
	err := facades.Orm().Query().Find(&permissions)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Error fetching permissions",
			"error":   err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Permission berhasil diambil",
		"data":    permissions,
	})
}

// Store: Menyimpan permission baru.
func (r *PermissionController) Store(ctx http.Context) http.Response {
	validator, err := ctx.Request().Validate(map[string]string{
		"name": "required|string|min_len:3|unique:permissions,name",
	})
	userID := ctx.Value("x-user-id").(string)
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

	permission := models.Permission{
		Name:      ctx.Request().Input("name"),
		CreatedBy: userID,
	}

	if err := facades.Orm().Query().Create(&permission); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal membuat permission",
		})
	}

	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "Success creating permission",
		"data":    permission,
	})
}

// fungsi show menampilkan berdasarkan ID
func (r *PermissionController) Show(ctx http.Context) http.Response {
	var permission models.Permission
	id := ctx.Request().Route("id")
	err := utils.FindModelByID(id, &permission)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Failed to find permission",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Success showing permission",
		"data":    permission,
	})
}

// Update: Memperbarui permission yang ada.
func (r *PermissionController) Update(ctx http.Context) http.Response {
	var permission models.Permission
	id := ctx.Request().Route("id")
	err := utils.FindModelByID(id, &permission)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Failed to find permission",
			"error":   err.Error(),
		})
	}

	validator, err := ctx.Request().Validate(map[string]string{
		"name": "required|string|min_len:3|unique:permissions,name," + id,
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
	userId, _ := utils.GetUserIDFromToken(ctx, ctx.Request().Header("Authorization"))
	permission.Name = ctx.Request().Input("name")
	permission.UpdatedBy = userId

	if err := facades.Orm().Query().Save(&permission); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Failed to update permission",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Success updating permission",
		"data":    permission,
	})
}

// Destroy: Menghapus permission.
func (r *PermissionController) Destroy(ctx http.Context) http.Response {
	var permission models.Permission
	id := ctx.Request().Route("id")

	err := utils.FindModelByID(id, &permission)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Permission not found",
			"error":   err.Error(),
		})
	}

	userId, _ := utils.GetUserIDFromToken(ctx, ctx.Request().Header("Authorization"))
	permission.UpdatedBy = userId

	_, err = facades.Orm().Query().Where("id", id).Delete(&permission)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menghapus permission",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "permission berhasil dihapus",
		"data":    id,
	})
}
