package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/osamikoyo/yoconf/casher"
	"github.com/osamikoyo/yoconf/config"
	"github.com/osamikoyo/yoconf/core"
	"github.com/osamikoyo/yoconf/grpcserver"
	"github.com/osamikoyo/yoconf/handler"
	"github.com/osamikoyo/yoconf/httpserver"
	"github.com/osamikoyo/yoconf/logger"
	"github.com/osamikoyo/yoconf/pb"
	"github.com/osamikoyo/yoconf/retrier"
	"github.com/osamikoyo/yoconf/storage"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	configPath := "config.yaml"

	for i, arg := range os.Args {
		if arg == "--config" {
			configPath = os.Args[i+1]
		}
	}

	logCfg := logger.Config{
		AppName:   "yoconf",
		LogFile:   "logs/yoconf.log",
		LogLevel:  "debug",
		AddCaller: false,
	}

	logger.Init(logCfg)
	logger := logger.Get()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		logger.Fatal("failed get config",
			zap.String("path", configPath),
			zap.Error(err))

		return
	}

	DBconn, err := retrier.Connect(5, func() (*gorm.DB, error) {
		return gorm.Open(sqlite.Open(cfg.DBPath))
	})
	if err != nil {
		logger.Fatal("failed connect to db",
			zap.String("path", cfg.DBPath),
			zap.Error(err))

		return
	}

	redisConn, err := retrier.Connect(5, func() (*redis.Client, error) {
		config := &redis.Options{
			DB:   0,
			Addr: cfg.RedisURL,
		}

		client := redis.NewClient(config)

		return client, client.Ping(context.Background()).Err()
	})
	if err != nil {
		logger.Fatal("failed connect to redis",
			zap.String("url", cfg.RedisURL),
			zap.Error(err))

		return
	}

	casher := casher.NewCasher(redisConn, logger)
	storage := storage.NewStorage(DBconn, logger, cfg)

	core := core.NewCore(casher, storage, logger, 30*time.Second)

	handler := handler.NewHandler(core)
	grpcserver := grpcserver.NewGRPCServer(core)
	httpserver := httpserver.NewHTTPServer(echo.New(), logger, cfg, handler)

	coreserver := grpc.NewServer()
	pb.RegisterYoConfServer(coreserver, grpcserver)

	go func() {
		<-ctx.Done()
		httpserver.Close(ctx)
		coreserver.GracefulStop()
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf("%s%d", cfg.Addr, cfg.GrpcPort))
	if err != nil {
		logger.Fatal("failed listen", zap.Error(err))
	}

	go func() {
		if err = coreserver.Serve(lis); err != nil {
			logger.Fatal("failed start coreserver", zap.Error(err))

			return
		}
	}()

	logger.Info("coreserver started")

	logger.Info("starting http server...")
	if err = httpserver.Run(ctx); err != nil {
		logger.Fatal("failed run http server", zap.Error(err))
	}
}
