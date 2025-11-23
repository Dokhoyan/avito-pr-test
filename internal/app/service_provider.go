package app

import (
	"context"
	"log"
	"time"

	"github.com/Dokhoyan/avito-pr-test/internal/config"
	prImpl "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/pr"
	statsImpl "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/stats"
	teamImpl "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/team"
	userImpl "github.com/Dokhoyan/avito-pr-test/internal/http-server/handlers/user"
	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	postgresRepo "github.com/Dokhoyan/avito-pr-test/internal/repository/postgres"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	prserv "github.com/Dokhoyan/avito-pr-test/internal/service/pr"
	statsserv "github.com/Dokhoyan/avito-pr-test/internal/service/stats"
	teamserv "github.com/Dokhoyan/avito-pr-test/internal/service/team"
	userserv "github.com/Dokhoyan/avito-pr-test/internal/service/user"
	"github.com/Dokhoyan/common/pkg/closer"
	txmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	httpConfig config.HTTPConfig

	dbPool    *pgxpool.Pool
	getter    *txmpgx.CtxGetter
	repo      *postgresRepo.Repository
	txManager *manager.Manager

	teamService  service.TeamService
	userService  service.UserService
	prService    service.PRService
	statsService service.StatsService

	userImpl  *userImpl.Implementation
	teamImpl  *teamImpl.Implementation
	prImpl    *prImpl.Implementation
	statsImpl *statsImpl.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			logger.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			logger.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}
	return s.httpConfig
}

func (s *serviceProvider) DBPool(ctx context.Context) *pgxpool.Pool {
	if s.dbPool == nil {
		poolConfig, err := pgxpool.ParseConfig(s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to parse connection string: %v", err)
		}

		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			log.Fatalf("failed to create connection pool: %v", err)
		}

		timeout := poolConfig.ConnConfig.ConnectTimeout
		if timeout == 0 {
			timeout = 5 * time.Second
		}

		pingCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		err = pool.Ping(pingCtx)
		if err != nil {
			log.Fatalf("failed to ping database: %v", err)
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.dbPool = pool
	}
	return s.dbPool
}

func (s *serviceProvider) TxManager(ctx context.Context) *manager.Manager {
	if s.txManager == nil {
		s.txManager = manager.Must(
			txmpgx.NewDefaultFactory(s.DBPool(ctx)),
		)
	}
	return s.txManager
}

func (s *serviceProvider) Getter(ctx context.Context) *txmpgx.CtxGetter {
	if s.getter == nil {
		s.getter = txmpgx.DefaultCtxGetter
	}
	return s.getter
}

func (s *serviceProvider) Repository(ctx context.Context) *postgresRepo.Repository {
	if s.repo == nil {
		s.repo = postgresRepo.New(s.DBPool(ctx), s.Getter(ctx))
	}
	return s.repo
}

func (s *serviceProvider) TeamService(ctx context.Context) service.TeamService {
	if s.teamService == nil {
		repo := s.Repository(ctx)
		s.teamService = teamserv.NewService(
			repo,
			repo,
			s.TxManager(ctx),
		)
	}
	return s.teamService
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		repo := s.Repository(ctx)
		s.userService = userserv.NewService(
			repo,
			s.TxManager(ctx),
		)
	}
	return s.userService
}

func (s *serviceProvider) PRService(ctx context.Context) service.PRService {
	if s.prService == nil {
		repo := s.Repository(ctx)
		s.prService = prserv.NewService(
			repo,
			repo,
			repo,
			s.TxManager(ctx),
		)
	}
	return s.prService
}

func (s *serviceProvider) StatsService(ctx context.Context) service.StatsService {
	if s.statsService == nil {
		repo := s.Repository(ctx)
		s.statsService = statsserv.NewService(
			repo,
			s.TxManager(ctx),
		)
	}
	return s.statsService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *userImpl.Implementation {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) TeamImpl(ctx context.Context) *teamImpl.Implementation {
	if s.teamImpl == nil {
		s.teamImpl = teamImpl.NewImplementation(s.TeamService(ctx))
	}

	return s.teamImpl
}

func (s *serviceProvider) PRImpl(ctx context.Context) *prImpl.Implementation {
	if s.prImpl == nil {
		s.prImpl = prImpl.NewImplementation(s.PRService(ctx))
	}

	return s.prImpl
}

func (s *serviceProvider) StatsImpl(ctx context.Context) *statsImpl.Implementation {
	if s.statsImpl == nil {
		s.statsImpl = statsImpl.NewImplementation(s.StatsService(ctx))
	}

	return s.statsImpl
}
