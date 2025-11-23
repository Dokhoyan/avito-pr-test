package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Dokhoyan/avito-pr-test/internal/config"
	"github.com/Dokhoyan/avito-pr-test/internal/http-server/middleware"
	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"github.com/Dokhoyan/common/pkg/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	loggerMaxSize    = 10
	loggerMaxBackups = 3
	loggerMaxAge     = 3
	reqLimite        = 100
	reqSecondtime    = 1
)

var configPath string
var logLevel = flag.String("level", "info", "log level for logger")

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		// Закрываем HTTP сервер перед закрытием других ресурсов
		if a.httpServer != nil {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()
			if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
				logger.Error("server shutdown error in defer", zap.Error(err))
			}
		}
		// Закрываем остальные ресурсы (база данных и т.д.)
		closer.CloseAll()
		closer.Wait()
	}()

	serverErrChan := make(chan error, 1)

	go func() {
		err := a.runHTTPServer()
		if err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
			return err
		}
		logger.Info("server stopped")
		return nil
	case err := <-serverErrChan:
		// При ошибке сервера тоже нужно корректно закрыть его
		logger.Error("server error occurred, shutting down", zap.Error(err))
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if shutdownErr := a.httpServer.Shutdown(shutdownCtx); shutdownErr != nil {
			logger.Error("server shutdown error after server error", zap.Error(shutdownErr))
		}
		return err
	}
}

func (a *App) initDeps(ctx context.Context) error {
	flag.Parse()

	inits := []func(context.Context) error{
		a.initLogger,
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()

	a.serviceProvider.UserImpl(ctx).RegisterRoutes(mux)
	a.serviceProvider.TeamImpl(ctx).RegisterRoutes(mux)
	a.serviceProvider.PRImpl(ctx).RegisterRoutes(mux)
	a.serviceProvider.StatsImpl(ctx).RegisterRoutes(mux)

	var handler http.Handler = mux
	handler = middleware.RequestLogger(handler)

	httpConfig := a.serviceProvider.HTTPConfig()
	a.httpServer = &http.Server{
		Addr:              httpConfig.Address(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.getCore(a.getLevel()))

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())
	logger.Info("starting server",
		zap.String("addr", a.serviceProvider.HTTPConfig().Address()))
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    loggerMaxSize,
		MaxBackups: loggerMaxBackups,
		MaxAge:     loggerMaxAge,
		Compress:   false,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func (a *App) getLevel() zap.AtomicLevel {
	var level zapcore.Level

	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level")
	}

	return zap.NewAtomicLevelAt(level)
}
