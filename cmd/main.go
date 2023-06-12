package main

import (
	"log"

	"github.com/majidbl/discount/config"
	"github.com/majidbl/discount/internal/server"
	"github.com/majidbl/discount/pkg/jaeger"
	"github.com/majidbl/discount/pkg/logger"
	"github.com/majidbl/discount/pkg/nats"
	"github.com/majidbl/discount/pkg/postgresql"
	"github.com/majidbl/discount/pkg/redis"
	"github.com/opentracing/opentracing-go"
)

// @title Discount microservice
// @version 1.0
// @description Discount microservice
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /api/v1
func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(logger.LogConfig{
		LogLevel: cfg.Logger.Level,
		DevMode:  false,
		Encoder:  "",
	})

	appLogger.InitLogger()
	appLogger.Info("Starting egiftcharges microservice")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.HTTP.Development,
	)
	appLogger.Infof("Success loaded config: %+v", cfg.AppVersion)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		appLogger.Fatalf("NewRedisClient: %+v", err)
	}

	appLogger.Infof("Redis connected: %+v", redisClient.PoolStats())

	natsConn, err := nats.NewNatsConnect(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("NewNatsConnect: %+v", err)
	}
	appLogger.Infof(
		"Nats Connected: Status: %+v IsConnected: %v ConnectedUrl: %v ConnectedServerId: %v",
		natsConn.NatsConn().Status(),
		natsConn.NatsConn().IsConnected(),
		natsConn.NatsConn().ConnectedUrl(),
		natsConn.NatsConn().ConnectedServerId(),
	)

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		appLogger.Fatalf("NewPgxConn: %+v", err)
	}
	appLogger.Infof("PostgreSQL connected: %+v", pgxPool.Stat().TotalConns())

	s := server.NewServer(appLogger, cfg, natsConn, pgxPool, tracer, redisClient)

	appLogger.Fatal(s.Run())
}
