package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type loadFileController struct {
	loadFileUseCase usecases.LoadFileUseCases
}

func NewLoadFileController(
	handler *gin.Engine,
	loadFileUseCase usecases.LoadFileUseCases,
	middleware middleware.Middleware,
) {
	lf := &loadFileController{
		loadFileUseCase: loadFileUseCase,
	}

	handler.POST("/load/file", lf.LoadFile, middleware.HandleErrors, middleware.Authenticate)
}

func (lf *loadFileController) LoadFile(ctx *gin.Context) {
	fmt.Println("loadFile")
	var req requests.LoadFile
	if err := ctx.ShouldBind(&req); err != nil {
		wrappedError := fmt.Errorf("%w: %w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	err := lf.loadFileUseCase.LoadFile(ctx, req)
	if err != nil {
		wrappedError := fmt.Errorf("There was a problem during load file: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, gin.H{})
}
