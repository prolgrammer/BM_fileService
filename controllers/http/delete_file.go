package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type deleteFileController struct {
	deleteFileUseCase usecases.DeleteFileUseCase
}

func NewDeleteFileController(
	engine *gin.Engine,
	deleteFileUseCase usecases.DeleteFileUseCase,
	middleware middleware.Middleware,
) {
	df := &deleteFileController{
		deleteFileUseCase: deleteFileUseCase,
	}

	engine.POST("/app/file/delete", middleware.Authenticate, df.DeleteFile, middleware.HandleErrors)
}

// DeleteFile godoc
// @Summary Удаление файлы
// @Description Удаляет файл
// @Tags Files
// @Accept json
// @Produce json
// @Param request body requests.File true "Название файла, папки и категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} string "Результат удаление"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/file/delete [post]
func (df *deleteFileController) DeleteFile(ctx *gin.Context) {
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

	err := df.deleteFileUseCase.DeleteFile(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("failed to delete file: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "File deleted successfully",
	})
}
