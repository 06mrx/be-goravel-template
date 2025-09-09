package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/http/limit"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {

}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {
	facades.RateLimiter().For("api_limiter", func(ctx http.Context) http.Limit {
		// Aturan: 60 permintaan per menit per IP
		return limit.PerMinute(1)
		// return facades.RateLimiter().PerMinute(60).
		//     By(ctx.Request().Ip()).
		//     Response(func(ctx http.Context) {
		//         ctx.Response().Status(http.StatusTooManyRequests).Json(http.Json{
		//             "message": "Terlalu banyak permintaan. Silakan coba lagi nanti.",
		//         })
		//     })
	})

}
