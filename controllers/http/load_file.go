package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prolgrammer/BM_authService/pkg/middleware"
	e "github.com/prolgrammer/BM_package/errors"
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

	}
}
