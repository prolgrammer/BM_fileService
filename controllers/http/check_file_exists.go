package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type checkFileExistsController struct {
	checkFileExistsUseCase usecases.CheckFileExistsUseCase
}

func NewCheckFileExistsUseCase(
	engine *gin.Engine,
	checkFileExistsUseCase usecases.CheckFileExistsUseCase,
	middleware middleware.Middleware,
) {
	cf := &checkFileExistsController{
		checkFileExistsUseCase: checkFileExistsUseCase,
	}

	engine.GET("/app/file/exists", middleware.Authenticate, cf.CheckFileExists, middleware.HandleErrors)
}

// CheckFileExists godoc
// @Summary      Проверка существования файла
// @Description  Проверяет, существует ли файл для указанного аккаунта
// @Tags         Files
// @Accept       json
// @Produce      json
// @Param        request body requests.File true "Название файла, папка, категория"
// @Param        Authorization header string true "Токен доступа"
// @Success      200 {object} map[string]bool "Результат проверки {\"exists\": true/false}"
// @Failure      400 {object} string "Некорректный формат запроса"
// @Failure      401 {object} string "Ошибка аутентификации"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /app/file/exists [get]
func (cf *checkFileExistsController) CheckFileExists(ctx *gin.Context) {
	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.File
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w: %w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	exists, err := cf.checkFileExistsUseCase.CheckFileExists(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("failed to check file existence: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}
	ctx.JSON(200, gin.H{
		"exists": exists,
	})
}
