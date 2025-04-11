package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
	"net/http"
)

type getFileController struct {
	getFileUseCase usecases.GetFileUseCase
}

func NewGetFileController(
	engine *gin.Engine,
	getFileUseCase usecases.GetFileUseCase,
	middleware middleware.Middleware) {
	gf := &getFileController{
		getFileUseCase: getFileUseCase,
	}

	engine.GET("/app/file", middleware.Authenticate, gf.GetFile, middleware.HandleErrors)
}

// GetFile godoc
// @Summary Получение конкретного файла
// @Description Возвращает файл конкретной категории и папки
// @Tags Files
// @Accept json
// @Produce json
// @Param request body requests.File true "Название файла, папки, категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} responses.File "Конкретный файл данного аккаунта"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/file [get]
func (gf *getFileController) GetFile(ctx *gin.Context) {
	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.File
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w:%w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	file, err := gf.getFileUseCase.GetFile(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("failed to get file:%w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"file": file})
}
