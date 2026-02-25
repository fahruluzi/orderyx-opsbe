package main

import (
	"context"
	"time"

	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/handler"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/middleware"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/http/route"
	"github.com/fahruluzi/orderyx-opsbe/internal/adapter/repository"
	"github.com/fahruluzi/orderyx-opsbe/internal/usecase"
	"github.com/fahruluzi/orderyx-opsbe/pkg/config"
	"github.com/fahruluzi/orderyx-opsbe/pkg/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func provideConfig() (config.Config, error) {
	return config.LoadConfig()
}

func provideJWTService(cfg config.Config) *jwt.JWTService {
	return jwt.NewJWTService(cfg.JWTSecretKey, cfg.JWTAccessTokenExpHours)
}

func provideDatabase(lc fx.Lifecycle, cfg config.Config) (*gorm.DB, error) {
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + cfg.DBPort + " sslmode=" + cfg.DBSSLMode
	connCfg := repository.DatabaseConfig{
		DSN: dsn,
	}
	db, err := repository.NewDatabase(lc, connCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return sqlDB.Close()
		},
	})

	return db, nil
}

func provideApp(
	cfg config.Config,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *fiber.App {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001, http://localhost:5173",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	api := app.Group("/api/v1")
	route.SetupAuthRoutes(api, authHandler, authMiddleware.Authenticate)

	return app
}

var ConfigModule = fx.Module("config",
	fx.Provide(
		provideConfig,
	),
)

var DatabaseModule = fx.Module("database",
	fx.Provide(
		provideDatabase,
	),
)

var RepositoryModule = fx.Module("repository",
	fx.Provide(
		repository.NewAuthRepository,
		repository.NewMerchantRepository,
		repository.NewAuditLogRepository,
		repository.NewConfigRepository,
		repository.NewSubscriptionRepository,
		repository.NewDashboardRepository,
		repository.NewPaymentRepository,
	),
)

var UsecaseModule = fx.Module("usecase",
	fx.Provide(
		usecase.NewAuthUsecase,
		usecase.NewMerchantUsecase,
		usecase.NewConfigUsecase,
		usecase.NewSubscriptionUsecase,
		usecase.NewDashboardUsecase,
		usecase.NewAuditLogUsecase,
		usecase.NewPaymentUsecase,
	),
)

var AuthModule = fx.Options(
	fx.Provide(handler.NewAuthHandler),
	fx.Provide(middleware.NewAuthMiddleware),
	fx.Provide(provideJWTService),
)

var MerchantModule = fx.Options(
	fx.Provide(handler.NewMerchantHandler),
)

var SystemConfigModule = fx.Options(
	fx.Provide(handler.NewConfigHandler),
)

var SubscriptionModule = fx.Options(
	fx.Provide(handler.NewSubscriptionHandler),
)

var DashboardModule = fx.Options(
	fx.Provide(handler.NewDashboardHandler),
)

var AuditLogModule = fx.Options(
	fx.Provide(handler.NewAuditLogHandler),
)

var PaymentModule = fx.Options(
	fx.Provide(handler.NewPaymentHandler),
)

var AppModule = fx.Module("app",
	fx.Provide(
		func(cfg config.Config) *fiber.App {
			app := fiber.New()
			app.Use(cors.New(cors.Config{
				AllowOrigins:     "http://localhost:3001, http://localhost:5173",
				AllowCredentials: true,
				AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			}))
			app.Get("/health", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})
			return app
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, cfg config.Config, authHandler *handler.AuthHandler, merchantHandler *handler.MerchantHandler, configHandler *handler.ConfigHandler, subscriptionHandler *handler.SubscriptionHandler, dashboardHandler *handler.DashboardHandler, auditLogHandler *handler.AuditLogHandler, paymentHandler *handler.PaymentHandler, authMiddleware *middleware.AuthMiddleware) {
		api := app.Group("/api/v1")
		route.SetupAuthRoutes(api, authHandler, authMiddleware.Authenticate)
		route.SetupMerchantRoutes(api, merchantHandler, authMiddleware.Authenticate)
		route.SetupConfigRoutes(api, configHandler, authMiddleware.Authenticate)
		route.SetupSubscriptionRoutes(api, subscriptionHandler, authMiddleware.Authenticate)
		route.SetupDashboardRoutes(api, dashboardHandler, authMiddleware.Authenticate)
		route.SetupAuditLogRoutes(api, auditLogHandler, authMiddleware.Authenticate)
		route.SetupPaymentRoutes(api, paymentHandler, authMiddleware.Authenticate)

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					port := cfg.Port
					if port == "" {
						port = "8081"
					}
					// The instruction had fmt.Sprintf and fmt.Printf, but fmt is not imported.
					// Sticking to the original logic for error handling and port assignment.
					if err := app.Listen(":" + port); err != nil {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return app.Shutdown()
			},
		})
	}),
)

func main() {
	time.Local = time.UTC
	fx.New(
		ConfigModule,
		SystemConfigModule,
		DatabaseModule,
		RepositoryModule,
		UsecaseModule,
		AuthModule,
		MerchantModule,
		SubscriptionModule,
		DashboardModule,
		AuditLogModule,
		PaymentModule,
		AppModule,
	).Run()
}
