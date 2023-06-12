package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/majidbl/discount/config"
	discountV1 "github.com/majidbl/discount/internal/discount/delivery/http/v1"
	discountNats "github.com/majidbl/discount/internal/discount/delivery/nats"
	discountRepo "github.com/majidbl/discount/internal/discount/repository"
	discountUC "github.com/majidbl/discount/internal/discount/usecase"
	giftChargeV1 "github.com/majidbl/discount/internal/giftcharge/delivery/http/v1"
	giftChargeRepo "github.com/majidbl/discount/internal/giftcharge/repository"
	giftChargeUC "github.com/majidbl/discount/internal/giftcharge/usecase"
	"github.com/majidbl/discount/internal/interceptors"
	"github.com/majidbl/discount/internal/middlewares"
	reportV1 "github.com/majidbl/discount/internal/report/delivery/http/v1"
	reportNats "github.com/majidbl/discount/internal/report/delivery/nats"
	reportRepo "github.com/majidbl/discount/internal/report/repository"
	reportUC "github.com/majidbl/discount/internal/report/usecase"
	"github.com/majidbl/discount/pkg/grpc_client"
	"github.com/majidbl/discount/pkg/logger"
)

const (
	maxHeaderBytes  = 1 << 20
	gzipLevel       = 5
	stackSize       = 1 << 10 // 1 KB
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
)

type server struct {
	log      logger.Logger
	cfg      *config.Config
	natsConn stan.Conn
	dbx      *pgxpool.Pool
	tracer   opentracing.Tracer
	echo     *echo.Echo
	redis    *redis.Client
}

// NewServer constructor
func NewServer(
	log logger.Logger,
	cfg *config.Config,
	natsConn stan.Conn,
	db *pgxpool.Pool,
	tracer opentracing.Tracer,
	redis *redis.Client,
) *server {
	return &server{
		log:      log,
		cfg:      cfg,
		natsConn: natsConn,
		dbx:      db,
		tracer:   tracer,
		redis:    redis,
		echo:     echo.New(),
	}
}

// Run start application
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	giftChargePgRepo := giftChargeRepo.NewGiftChargePGRepository(s.dbx)
	giftChargeRedisRepo := giftChargeRepo.NewGiftChargeRedisRepository(s.redis)
	giftChargeUseCase := giftChargeUC.NewGiftChargeUseCase(s.log, giftChargePgRepo, giftChargeRedisRepo)

	publisher := reportNats.NewPublisher(s.natsConn)
	reportPgRepo := reportRepo.NewReportPGRepository(s.dbx)
	reportRedisRepo := reportRepo.NewReportRedisRepository(s.redis)
	reportUseCase := reportUC.NewReportUseCase(s.log, reportPgRepo, publisher, reportRedisRepo)

	interceptorManager := interceptors.NewInterceptorManager(s.log, nil)
	walletGrpcClient, _, err := grpc_client.NewWalletGrpcClient(
		ctx,
		":5007",
		interceptorManager)

	discountPublisher := discountNats.NewPublisher(s.natsConn)
	discountRedisRepo := discountRepo.NewDiscountRedisRepository(s.redis)
	discountUseCase := discountUC.NewDiscountUseCase(
		s.log, reportPgRepo,
		walletGrpcClient,
		discountPublisher,
		discountRedisRepo,
		giftChargePgRepo,
		giftChargeRedisRepo,
	)

	mw := middlewares.NewMiddlewareManager(s.log, s.cfg)
	validate := validator.New()
	v1 := s.echo.Group("/api/v1")
	v1.Use(mw.Metrics)

	giftChargeHandlers := giftChargeV1.NewGiftChargeHandlers(
		v1.Group("/giftCharge"),
		giftChargeUseCase,
		s.log,
		validate)
	giftChargeHandlers.MapRoutes()

	reportHandlers := reportV1.NewReportHandlers(
		v1.Group("/report"),
		reportUseCase,
		s.log,
		validate)
	reportHandlers.MapRoutes()

	discountHandlers := discountV1.NewDiscountHandlers(
		v1.Group("/discount"),
		discountUseCase,
		s.log,
		validate)
	discountHandlers.MapRoutes()

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.HTTP.Port)
		s.runHttpServer()
	}()

	metricsServer := echo.New()
	go func() {
		metricsServer.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Metrics.Port)
		if err := metricsServer.Start(s.cfg.Metrics.Port); err != nil {
			s.log.Error(err)
			cancel()
		}
	}()

	l, err := net.Listen("tcp", s.cfg.GRPC.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.log.Errorf("ctx.Done: %v", done)
	}

	if err = s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	if err = metricsServer.Shutdown(ctx); err != nil {
		s.log.Errorf("metricsServer.Shutdown: %v", err)
	}

	return nil
}
