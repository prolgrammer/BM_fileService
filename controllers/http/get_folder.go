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

type getFolderController struct {
	getFolderUseCase usecases.SelectFolderUseCase
}

func NewGetFolderController(
	engine *gin.Engine,
	getFolderUseCase usecases.SelectFolderUseCase,
	middleware middleware.Middleware,
) {

	gf := &getFolderController{
		getFolderUseCase: getFolderUseCase,
	}

	engine.GET("/app/folder/get", middleware.Authenticate, gf.GetFolder, middleware.HandleErrors)
}

func (gf *getFolderController) GetFolder(ctx *gin.Context) {
	fmt.Println("get folder")

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

	response, err := gf.getFolderUseCase.SelectFolder(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during get folder: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
