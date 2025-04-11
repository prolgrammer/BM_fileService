package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type createFileController struct {
	createFileUseCases usecases.CreateFileUseCases
}

func NewCreateFileController(
	engine *gin.Engine,
	createFileUseCases usecases.CreateFileUseCases,
	middleware middleware.Middleware) {

	cf := &createFileController{
		createFileUseCases: createFileUseCases,
	}

	engine.POST("/app/file", middleware.Authenticate, cf.CreateFile, middleware.HandleErrors)
}

// CreateFile godoc
// @Summary Добавление файла
// @Description Создание и добавление файла в базу данных
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param request body requests.CreateFile true "Файл, версия файлы, название папки, название категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} string "Результат создания"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/file [post]
func (cf *createFileController) CreateFile(ctx *gin.Context) {
	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.CreateFile
	if err := ctx.ShouldBind(&req); err != nil {
		wrappedError := fmt.Errorf("%w: %w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}
	fmt.Println("req:", req)

	err := cf.createFileUseCases.CreateFile(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("failed to upload files: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Files uploaded successfully",
	})
}
