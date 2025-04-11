package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type getAllFilesController struct {
	getAllFilesUseCase usecases.GetAllFilesUseCase
}

func NewGetAllFilesUseCase(
	engine *gin.Engine,
	getAllFilesUseCase usecases.GetAllFilesUseCase,
	middleware middleware.Middleware,
) {

	gaf := getAllFilesController{
		getAllFilesUseCase: getAllFilesUseCase,
	}

	engine.GET("/app/files", middleware.Authenticate, gaf.GetAllFiles, middleware.HandleErrors)
}

// GetAllFiles godoc
// @Summary Получение всех файлов
// @Description Возвращает все файлы конкретной категории и папки
// @Tags Files
// @Accept json
// @Produce json
// @Param request body requests.File true "Название файла, папки, категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} []responses.File "Файлы в категории и папке данного аккаунта"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/files [get]
func (gaf *getAllFilesController) GetAllFiles(ctx *gin.Context) {
	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.Folder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w: %w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	fmt.Println(req)
	files, err := gaf.getAllFilesUseCase.GetAllFiles(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("failed to get files: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, gin.H{
		"files": files,
	})
}
