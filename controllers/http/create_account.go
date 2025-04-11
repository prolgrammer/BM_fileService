package http

import (
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
	"net/http"
)

type createAccountController struct {
	createAccountUseCase usecases.CreateAccountUseCase
}

func NewCreateAccountController(
	engine *gin.Engine,
	createAccountUseCase usecases.CreateAccountUseCase,
	middleware middleware.Middleware) {
	ac := &createAccountController{
		createAccountUseCase: createAccountUseCase,
	}

	engine.POST("/app/account", middleware.Authenticate, ac.CreateAccount, middleware.HandleErrors)
}

// CreateAccount godoc
// @Summary Создание аккаунта пользователя
// @Description Создает аккаунт для пользователя, с последующем созданием категорий
// @Tags Account
// @Accept json
// @Produce json
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} string "Результат создания"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/account [post]
func (ac *createAccountController) CreateAccount(ctx *gin.Context) {
	fmt.Println("create account")

	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	err := ac.createAccountUseCase.CreateAccount(ctx, accountId.(string))
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during create account: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"answer": "account created successfully"})
}
