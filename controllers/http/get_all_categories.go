package http

import (
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	e "github.com/prolgrammer/BM_package/errors"
	"github.com/prolgrammer/BM_package/middleware"
	"net/http"
)

type getAllCategories struct {
	getAllCategoriesUseCase usecases.GetAllCategoryUseCase
}

func NewGetAllCategoriesUseCases(
	engine *gin.Engine,
	getAllCategoriesUseCase usecases.GetAllCategoryUseCase,
	middleware middleware.Middleware,
) {

	gac := &getAllCategories{
		getAllCategoriesUseCase: getAllCategoriesUseCase,
	}

	engine.GET("/app/category/all/get", middleware.Authenticate, gac.GetAllCategories, middleware.HandleErrors)
}

func (gac *getAllCategories) GetAllCategories(ctx *gin.Context) {
	fmt.Println("get all categories")

	accountId, exist := ctx.Get("account_id")
	if !exist {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	response, err := gac.getAllCategoriesUseCase.GetAllCategory(ctx, accountId.(string))
	if err != nil {
		middleware.AddGinError(ctx, fmt.Errorf("something went wrong: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, response)
}
