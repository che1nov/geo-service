package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"net/http"

	"geo-service/internal/config"
	"geo-service/internal/db"
	"geo-service/internal/handlers"
	"geo-service/internal/repository"
	"geo-service/internal/router"
	"geo-service/internal/service"

	_ "net/http/pprof"
)

type App struct {
	logger         *logrus.Logger
	cfg            *config.Config
	db             *db.DB
	redisClient    *redis.Client
	userRepo       repository.UserRepository
	daDataRepo     repository.AddressRepository
	userService    *service.UserService
	addressService *service.AddressService
	server         *http.Server
}

func NewApp(logger *logrus.Logger) *App {
	app := &App{
		logger: logger,
	}
	app.loadConfig()
	app.initRedis() // ! Добавьте эту строку
	app.connectToDatabase()
	app.initRepositories()
	app.initServices()
	app.initServer()

	return app
}

func (a *App) loadConfig() {
	if err := godotenv.Load(); err != nil {
		a.logger.Fatalf("Ошибка загрузки .env: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		a.logger.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	a.cfg = cfg
}

func (a *App) initRedis() {
	a.redisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	a.logger.Info("Redis клиент инициализирован!")
}

func (a *App) connectToDatabase() {
	// Исправлено: Убран dataSourceName — NewDB() не принимает параметров
	var err error
	var dbInstance *db.DB
	for i := 0; i < 10; i++ {
		dbInstance, err = db.NewDB() // ! Удалите параметр dataSourceName
		if err == nil {
			break
		}
		a.logger.Errorf("Ошибка подключения к БД. Повторная попытка через 2 секунды... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		a.logger.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	a.db = dbInstance
	a.logger.Info("База данных готова!")
}

func (a *App) initRepositories() {
	// Исправлено: Используйте a.db.DB для доступа к sql.DB
	a.userRepo = repository.NewDBUserRepository(a.db.DB)
	a.daDataRepo = repository.NewDaDataRepository(
		a.cfg.DaDataAPIKey,
		a.cfg.DaDataURL,
		a.logger,
		a.redisClient,
	)
}

func (a *App) initServices() {
	a.userService = service.NewUserService(a.userRepo)
	a.addressService = service.NewAddressService(a.daDataRepo)
}

func (a *App) initServer() {
	authHandler := handlers.NewAuthHandler(a.userService, a.logger)
	addressHandler := handlers.NewAddressHandler(a.addressService, a.logger)

	r := router.NewChiRouter(addressHandler, authHandler)

	a.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func (a *App) Run() {
	go func() {
		a.logger.Info("Запуск сервера на порту 8080...")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.WithError(err).Fatal("Ошибка запуска сервера")
		}
	}()

	startPprofServer(a.logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Получен сигнал завершения. Останавливаю сервер...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.WithError(err).Fatal("Ошибка при graceful shutdown")
	}

	a.logger.Info("Сервер успешно остановлен.")
}

func startPprofServer(logger *logrus.Logger) {
	go func() {
		logger.Info("Запуск сервера pprof на порту :6060...")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logger.Fatalf("Ошибка запуска сервера pprof: %v", err)
		}
	}()
}
