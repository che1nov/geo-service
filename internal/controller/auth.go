package controller

import (
	"encoding/json"
	"net/http"

	"example/internal/models"
	"example/internal/service"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Регистрирует нового пользователя. Принимает имя пользователя и пароль, сохраняет в in-memory хранилище.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        registerRequest  body      models.RegisterRequest  true  "Данные для регистрации"
// @Success      200              {string}  string  "OK"
// @Failure      400              {object}  map[string]string  "Неверный запрос или пользователь уже существует"
// @Router       /api/register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := c.authService.Register(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Выполняет вход пользователя и возвращает JWT-токен при успешной аутентификации.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRequest  body      models.LoginRequest  true  "Данные для аутентификации"
// @Success      200           {object}  models.LoginResponse  "JWT-токен"
// @Failure      400           {object}  map[string]string  "Неверный запрос"
// @Failure      401           {object}  map[string]string  "Пользователь не найден или неверный пароль"
// @Router       /api/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	token, err := c.authService.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}
