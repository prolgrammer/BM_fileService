package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type createFolderController struct {
	createFolderUseCase usecases.CreateFolderUseCase
}

func NewCreateFolderController(
	engine *gin.Engine,
	createFolder usecases.CreateFolderUseCase,
	middleware middleware.Middleware,
) {
	cf := &createFolderController{
		createFolderUseCase: createFolder,
	}

	engine.POST("/app/folder", middleware.Authenticate, cf.CreateFolder, middleware.HandleErrors)
}

func (cf *createFolderController) CreateFolder(ctx *gin.Context) {
	fmt.Println("create folder")

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

	err := cf.createFolderUseCase.CreateFolder(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during create folder: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, gin.H{
		"answer": "folder create successful"})
}
