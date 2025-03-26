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

type CheckFolderExistController struct {
	checkFolderExistUseCase usecases.CheckFolderExistUseCase
}

func NewCheckFolderExistController(
	engine *gin.Engine,
	checkCategoryExistUseCase usecases.CheckFolderExistUseCase,
	middleware middleware.Middleware,
) {

	checkFolderExistController := &CheckFolderExistController{
		checkFolderExistUseCase: checkCategoryExistUseCase,
	}

	engine.GET("/app/folder/exist", middleware.Authenticate, checkFolderExistController.CheckFolderExist, middleware.HandleErrors)
}

func (c *CheckFolderExistController) CheckFolderExist(ctx *gin.Context) {
	fmt.Println("check folder exist")

	accountId, exist := ctx.Get("account_id")
	if !exist {
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

	response, err := c.checkFolderExistUseCase.CheckFolderExist(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during check category exist: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
