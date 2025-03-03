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
