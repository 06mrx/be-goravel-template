package http

import (
	"goravel_api/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
)

type Kernel struct{}

func (k *Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
		// Tambahkan global middleware lain di sini jika perlu
		middleware.User(),
	}
}
