package http

import (
	"app/controllers/requests"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
)

type getAllFoldersController struct {
	getAllFoldersUseCase usecases.SelectFoldersUseCase
}

func NewGetAllFoldersController(
	engine *gin.Engine,
	getAllFoldersUseCase usecases.SelectFoldersUseCase,
	middleware middleware.Middleware) {

	gaf := &getAllFoldersController{
		getAllFoldersUseCase: getAllFoldersUseCase,
	}

	engine.GET("/app/folders", middleware.Authenticate, gaf.GetAllFolders, middleware.HandleErrors)
}

func (gaf *getAllFoldersController) GetAllFolders(ctx *gin.Context) {
	fmt.Println("get all folders")

	accountId, exists := ctx.Get("account_id")
	if !exists {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.Category
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w:%w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	response, err := gaf.getAllFoldersUseCase.SelectFolders(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during get all folders: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(200, response)
}
