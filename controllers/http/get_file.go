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

	engine.GET("app/file", middleware.Authenticate, gf.GetFile, middleware.HandleErrors)
}

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
