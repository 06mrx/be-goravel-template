package controllers

import (
	"goravel_api/app/models"
	"goravel_api/app/utils"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type AuditController struct {
	// Dependent services
}

func NewAuditController() *AuditController {
	return &AuditController{
		// Inject services
	}
}

func (r *AuditController) Index(ctx http.Context) http.Response {
	allowed := utils.CheckPermission(ctx, "view-audits")
	if !allowed {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized action"})
	}

	page := ctx.Request().QueryInt("page", 1)
	per_page := ctx.Request().QueryInt("per_page", 10)
	searchQuery := ctx.Request().Query("search", "")
	sort_order := ctx.Request().Query("sort_order", "ASC")
	query := facades.Orm().Query().OrderBy("created_at", sort_order)

	if searchQuery != "" {
		query = query.Where("event ILIKE ?", "%"+searchQuery+"%").OrWhere("auditable_type ILIKE ?", "%"+searchQuery+"%")
	}

	var audits []models.Audit
	var total int64
	err := query.Paginate(page, per_page, &audits, &total)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil audit",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Audit berhasil diambil",
		"data": http.Json{
			"total":        total,
			"per_page":     per_page,
			"current_page": page,
			"last_page":    (total + int64(per_page) - 1) / int64(per_page),
			"data":         audits,
		},
	})
}
