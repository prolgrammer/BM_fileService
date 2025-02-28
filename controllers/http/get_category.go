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

type getCategoryController struct {
	getCategoryUseCase usecases.GetCategoryUseCase
}

func NewGetCategoryController(
	engine *gin.Engine,
	getCategoryUseCase usecases.GetCategoryUseCase,
	middleware middleware.Middleware,
) {

	gc := &getCategoryController{
		getCategoryUseCase: getCategoryUseCase,
	}

	engine.GET("/app/category/get", middleware.Authenticate, gc.GetCategory, middleware.HandleErrors)
}

func (gc *getCategoryController) GetCategory(ctx *gin.Context) {
	fmt.Println("get category")

	accountId, exist := ctx.Get("account_id")
	if !exist {
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

	response, err := gc.getCategoryUseCase.GetCategory(ctx, accountId.(string), req)
	if err != nil {
		middleware.AddGinError(ctx, fmt.Errorf("something went wrong while get category: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
