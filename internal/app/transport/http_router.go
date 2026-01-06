package transport

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/ijalalfrz/coupon-system/internal/app/config"
	"github.com/ijalalfrz/coupon-system/internal/app/dto"
	"github.com/ijalalfrz/coupon-system/internal/app/endpoints"
	httptransport "github.com/ijalalfrz/coupon-system/internal/pkg/transport/http"
)

// MakeHTTPRouter builds the HTTP router with all the service endpoints.
func MakeHTTPRouter(
	cfg *config.Config,
	endpts endpoints.Endpoints,
) *chi.Mux {
	// Initialize Router
	router := chi.NewRouter()

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	router.Route("/api/coupons", func(router chi.Router) {
		router.Use(
			httptransport.CORSMiddleware(),
			httptransport.Recoverer(slog.Default()),
			render.SetContentType(render.ContentTypeJSON),
		)

		router.Post("/", httptransport.MakeHandlerFunc(
			endpts.CouponEndpoint.CreateCoupon,
			httptransport.DecodeRequest[dto.CreateCouponRequest],
			httptransport.CreatedResponse,
		))

		router.Post("/claim", httptransport.MakeHandlerFunc(
			endpts.CouponClaimEndpoint.ClaimCoupon,
			httptransport.DecodeRequest[dto.ClaimCouponRequest],
			httptransport.CreatedResponse,
		))

		router.Get("/{name}", httptransport.MakeHandlerFunc(
			endpts.CouponEndpoint.GetCouponByCouponName,
			httptransport.DecodeRequest[dto.GetCouponRequest],
			httptransport.ResponseWithBody,
		))
	})

	return router
}
