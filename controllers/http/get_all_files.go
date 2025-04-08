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
