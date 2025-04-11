package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type deleteFolderController struct {
	deleteFolderUseCase usecases.DeleteFolderUseCase
}

func NewDeleteFolderController(
	engine *gin.Engine,
	deleteFolderUseCase usecases.DeleteFolderUseCase,
	middleware middleware.Middleware,
) {
	df := deleteFolderController{
		deleteFolderUseCase: deleteFolderUseCase,
	}

	engine.POST("/app/folder/delete", middleware.Authenticate, df.DeleteFolder, middleware.HandleErrors)
}

// DeleteFolder godoc
// @Summary Удаление папки и лежащие в ней файлы
// @Description Удаляет категорию, ее папки и файлы
// @Tags Folders
// @Accept json
// @Produce json
// @Param request body requests.Folder true "Название папки, категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} string "Результат удаление"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/folder/delete [post]
func (df deleteFolderController) DeleteFolder(ctx *gin.Context) {
	fmt.Println("delete folder")

	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.Folder
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w:%w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	err := df.deleteFolderUseCase.DeleteFolder(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during delete folder: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, gin.H{
		"answer": "folder delete successful"})
}
