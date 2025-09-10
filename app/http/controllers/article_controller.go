package controllers

import (
	"fmt"
	stdhttp "net/http" // Mengganti nama impor untuk menghindari bentrokan

	"strings"
	// "time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	// "github.com/pkg/errors"

	"goravel_api/app/models"
	"goravel_api/app/utils"
)

type ArticleController struct {
	// Dependent services
}

func NewArticleController() *ArticleController {
	return &ArticleController{}
}

// slugify creates a URL-friendly slug from a string.
func slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, "--", "-")
	return s
}

func (r *ArticleController) Index(ctx http.Context) http.Response {
	allowed := utils.CheckPermission(ctx, "view-articles")
	if !allowed {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized"})
	}

	page := ctx.Request().QueryInt("page", 1)
	per_page := ctx.Request().QueryInt("per_page", 10)
	searchQuery := ctx.Request().Query("search", "")
	// userIDStr, _ := facades.Auth(ctx).ID()
	// query := facades.Orm().Query().With("User").Where("user_id", userIDStr)
	query := facades.Orm().Query().With("User")

	if searchQuery != "" {
		query = query.Where("title ILIKE ?", "%"+searchQuery+"%").OrWhere("content ILIKE ?", "%"+searchQuery+"%")
	}

	var articles []models.Article
	var total int64
	err := query.Paginate(page, per_page, &articles, &total)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil artikel",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Artikel berhasil diambil",
		"data": http.Json{
			"total":        total,
			"per_page":     per_page,
			"current_page": page,
			"last_page":    (total + int64(per_page) - 1) / int64(per_page),
			"data":         articles,
		},
	})
}

func (r *ArticleController) Store(ctx http.Context) http.Response {
	// Memeriksa otorisasi
	allowed := utils.CheckPermission(ctx, "store-articles")
	if !allowed {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized"})
	}

	validator, _ := ctx.Request().Validate(map[string]string{
		"title":   "required|string|min_len:3|unique:articles,title",
		"content": "required|string|min_len:3",
		"status":  "required|string",
		"image":   "required|file", // Validasi file akan dilakukan secara manual
	})
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{"message": validator.Errors().All()})
	}
	allowedExt := []string{"jpg", "png", "jpeg", "webp"}
	file, err := ctx.Request().File("image")
	maxUploadSize := int64(1 * 1024 * 1024)

	if err != nil {
		return ctx.Response().Json(stdhttp.StatusBadRequest, http.Json{"message": err.Error()})
	}

	isAllowed := utils.CheckFileTypeAndSize(allowedExt, maxUploadSize, &file)
	if !isAllowed {
		return ctx.Response().Json(stdhttp.StatusBadRequest, http.Json{
			"message": http.Json{
				"image": http.Json{
					"error": fmt.Sprintf(
						"file harus bertipe salah satu dari: %s. Maksimal ukuran file yang diizinkan: %d MB",
						strings.Join(allowedExt, ", "), (maxUploadSize / 1024 / 1024),
					),
				},
			},
		})

	}

	// Dapatkan ID pengguna dari token
	userIDStr, err := facades.Auth(ctx).ID()
	if err != nil {
		return ctx.Response().Json(http.StatusUnauthorized, http.Json{"message": "User ID tidak ditemukan"})
	}

	// Buat slug dari judul
	title := ctx.Request().Input("title")
	slug := slugify(title)

	// fmt.Println(slug)
	// Buat model artikel
	article := models.Article{
		Title:     title,
		Slug:      slug,
		Content:   ctx.Request().Input("content"),
		Status:    ctx.Request().Input("status"),
		CreatedBy: userIDStr,
		UpdatedBy: userIDStr,
		UserID:    userIDStr,
	}

	// fmt.Println("article ", article)
	// Dapatkan file yang diunggah

	path, err := facades.Storage().PutFile("articles/images", file)
	article.ImageUrl = path
	fmt.Println("path", path)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengunggah gambar.",
			"error":   err.Error(),
		})
	}

	if err := facades.Orm().Query().Create(&article); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal membuat artikel",
			"error":   err.Error(),
		})
	}
	fmt.Println("article saved ", article)
	return ctx.Response().Json(http.StatusCreated, http.Json{
		"message": "Artikel berhasil dibuat",
		"data":    article,
	})
}

func (r *ArticleController) Show(ctx http.Context) http.Response {
	allowed := utils.CheckPermission(ctx, "view-articles")
	if !allowed {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized"})
	}

	articleIDStr := ctx.Request().Route("id")
	var article models.Article
	// userIDStr, _ := facades.Auth(ctx).ID()
	// err := facades.Orm().Query().With("User").Where("id", articleIDStr).Where("user_id", userIDStr).First(&article)
	err := facades.Orm().Query().With("User").Where("id", articleIDStr).First(&article)
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return ctx.Response().Json(http.StatusNotFound, http.Json{"message": "Artikel tidak ditemukan."})
	// }
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal mengambil artikel.",
			"error":   err.Error(),
		})
	}
	// article.ImageUrl = facades.Storage().Url(article.ImageUrl)
	// fmt.Println("article ImageUrl", facades.Storage().Url(article.ImageUrl))

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Artikel berhasil diambil.",
		"data":    article,
	})
}

func (r *ArticleController) Update(ctx http.Context) http.Response {
	// Memeriksa otorisasi
	allowed := utils.CheckPermission(ctx, "update-articles")

	if !allowed {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized"})
	}

	articleIDStr := ctx.Request().Route("id")
	// userIDStr, _ := facades.Auth(ctx).ID()
	var article models.Article
	err := utils.FindModelByID(articleIDStr, &article)
	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Failed to find article",
			"error":   err.Error(),
		})
	}

	validator, _ := ctx.Request().Validate(map[string]string{
		"title":   "required|string|min_len:3|unique:articles,title," + articleIDStr,
		"content": "required|string|min_len:3",
		"status":  "required|string",
	})
	allowedExt := []string{"jpg", "png", "jpeg", "webp"}
	file, err := ctx.Request().File("image")
	maxUploadSize := int64(1 * 1024 * 1024)
	// extension, _ := file.Extension()
	// size, _ := file.Size()
	if err != nil {
		return ctx.Response().Json(stdhttp.StatusBadRequest, http.Json{"message": err.Error()})
	}
	isAllowed := utils.CheckFileTypeAndSize(allowedExt, maxUploadSize, &file)
	if !isAllowed {
		return ctx.Response().Json(stdhttp.StatusBadRequest, http.Json{
			"message": http.Json{
				"image": http.Json{
					"error": fmt.Sprintf(
						"file harus bertipe salah satu dari: %s. Maksimal ukuran file yang diizinkan: %d MB",
						strings.Join(allowedExt, ", "), (maxUploadSize / 1024 / 1024),
					),
				},
			},
		})
	}
	if err != nil {
		return ctx.Response().Json(stdhttp.StatusBadRequest, http.Json{"message": err.Error()})
	}
	if validator.Fails() {
		return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{"message": validator.Errors().All()})
	}

	// Dapatkan ID pengguna
	userIDStr, _ := facades.Auth(ctx).ID()
	title := ctx.Request().Input("title")
	slug := slugify(title)
	if slug != article.Slug {
		var existingArticle models.Article
		err := facades.Orm().Query().Where("slug", slug).First(&existingArticle)
		if err != nil {
			return ctx.Response().Json(http.StatusConflict, http.Json{"message": "Slug sudah ada."})
		}
	}

	article.Title = title
	article.Slug = slug
	article.Content = ctx.Request().Input("content")
	article.Status = ctx.Request().Input("status")
	article.UpdatedBy = userIDStr // Menetapkan UpdatedBy

	if err != nil {
		fmt.Println("gagal dapat file")
	} else {
		oldPath := strings.TrimPrefix(article.ImageUrl, "storage/")
		facades.Storage().Delete(oldPath)

		path, err := facades.Storage().PutFile("articles/images", file)
		article.ImageUrl = path
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, http.Json{
				"message": "Gagal mengunggah gambar.",
				"error":   err.Error(),
			})
		}
	}

	// Simpan perubahan
	if err := facades.Orm().Query().Save(&article); err != nil { //error disini
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal memperbarui artikel.",
			"error":   err.Error(),
		})
	}

	fmt.Println("after save")
	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Artikel berhasil diperbarui.",
		"data":    article,
	})
}

func (r *ArticleController) Destroy(ctx http.Context) http.Response {
	// Memeriksa otorisasi
	allowed := utils.CheckPermission(ctx, "delete-articles")
	if !allowed {
		return ctx.Response().Json(http.StatusForbidden, http.Json{"message": "Unauthorized"})
	}

	articleIDStr := ctx.Request().Route("id")
	// articleID,_ := uuid.Parse(articleIDStr)
	var article models.Article
	err := utils.FindModelByID(articleIDStr, &article)

	if err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Artikel tidak ditemukan",
			"error":   err.Error(),
		})

	}

	// Hapus file gambar
	if article.ImageUrl != "" {
		path := strings.TrimPrefix(article.ImageUrl, "storage/")
		facades.Storage().Delete(path)
	}

	// Hapus artikel dari database (soft delete)
	if _, err := facades.Orm().Query().Where("id", articleIDStr).Delete(&article); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Gagal menghapus artikel.",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Artikel berhasil dihapus.",
	})
}
