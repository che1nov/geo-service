package main

import (
	"example/internal/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"log"

	"example/internal/controller"
	"example/internal/repository"
	"example/internal/service"

	_ "example/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Address API
// @version         1.0
// @description     API для поиска и геокодирования адресов с использованием DaData API.
// @termsOfService  http://swagger.io/terms/
// @host      localhost:8080
// @BasePath  /api
//
// @securityDefinitions.apikey  ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {
	// Инициализация JWT авторизации
	tokenAuth := jwtauth.New("HS256", []byte("193df3dae9f5da653f15242535ef418b201b3ac1"), nil)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository()   // Репозиторий пользователей
	dadataRepo := repository.NewDataRepository() // Репозиторий для вызова DaData API

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, tokenAuth)
	addressService := service.NewAddressService(dadataRepo)

	// Инициализация контроллеров
	authController := controller.NewAuthController(authService)
	defaultResponder := controller.DefaultResponder{}
	addressController := controller.NewAddressController(addressService, defaultResponder)

	// Настройка роутера
	r := chi.NewRouter()

	// Публичные эндпоинты
	r.Post("/api/register", authController.Register)
	r.Post("/api/login", authController.Login)

	// Защищённые эндпоинты
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/api/address/search", addressController.Search)
		r.Post("/api/address/geocode", addressController.Geocode)
	})

	// Маршрут для Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Создаем сервер на порту 8080
	srv := server.NewServer(":8080", r)

	// Запускаем сервер и ждем его остановки
	if err := srv.Serve(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
