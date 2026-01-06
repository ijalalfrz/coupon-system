package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ijalalfrz/coupon-system/internal/app/config"
	"github.com/ijalalfrz/coupon-system/internal/app/endpoints"
	"github.com/ijalalfrz/coupon-system/internal/app/repository"
	"github.com/ijalalfrz/coupon-system/internal/app/service"
	"github.com/ijalalfrz/coupon-system/internal/app/transport"
	"github.com/ijalalfrz/coupon-system/internal/pkg/db"
	"github.com/ijalalfrz/coupon-system/internal/pkg/logger"
)

func main() {

	cfg := config.MustInitConfig(".env")
	logger.InitStructuredLogger(cfg.LogLevel)
	runApp(cfg)
}

func runApp(cfg config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	slog.InfoContext(ctx, "starting...", slog.String("log_level", string(cfg.LogLevel)))

	var waitGroup sync.WaitGroup
	// Starts the server in a go routine
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		startHTTPServer(ctx, cfg)
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigChannel:
		cancel()
		slog.InfoContext(ctx, "received OS signal. Exiting...", slog.String("signal", sig.String()))
	case <-ctx.Done():
		slog.ErrorContext(ctx, "failed to start HTTP server")
	}

	waitGroup.Wait()
	slog.InfoContext(ctx, "All service closed...")
}

func startHTTPServer(ctx context.Context, cfg config.Config) {
	endpts := makeEndpoints(ctx, &cfg)
	router := transport.MakeHTTPRouter(&cfg, endpts)
	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		WriteTimeout: cfg.HTTP.Timeout,
		ReadTimeout:  cfg.HTTP.Timeout,
	}

	slog.Info("running HTTP server...", slog.Int("port", cfg.HTTP.Port))

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "failed to start HTTP server", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to shutdown HTTP server", slog.String("error", err.Error()))
	}

	slog.InfoContext(ctx, "HTTP server shutdown gracefully")
}

func makeEndpoints(ctx context.Context, cfg *config.Config) endpoints.Endpoints {
	// init db
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		slog.ErrorContext(ctx, "failed to init db", slog.String("error", err.Error()))
	}

	// init repository
	couponRepository := repository.NewCouponRepository(dbConn)
	couponClaimRepository := repository.NewCouponClaimRepository(dbConn)

	// init service endpoint
	return endpoints.Endpoints{
		CouponEndpoint:      makeCouponEndpoint(couponRepository, couponClaimRepository),
		CouponClaimEndpoint: makeCouponClaimEndpoint(couponRepository, couponClaimRepository),
	}
}

// makeCouponEndpoint creates a coupon endpoint
func makeCouponEndpoint(
	couponRepository *repository.CouponRepository,
	couponClaimRepository *repository.CouponClaimRepository,
) endpoints.CouponEndpoint {

	service := service.NewCouponService(
		couponRepository,
		couponRepository,
		couponClaimRepository,
	)

	return endpoints.MakeCouponEndpoint(service)
}

// makeCouponClaimEndpoint creates a coupon claim endpoint
func makeCouponClaimEndpoint(
	couponRepository *repository.CouponRepository,
	couponClaimRepository *repository.CouponClaimRepository,
) endpoints.CouponClaimEndpoint {

	service := service.NewCouponClaimService(
		couponRepository,
		couponClaimRepository,
	)

	return endpoints.MakeCouponClaimEndpoint(service)
}
