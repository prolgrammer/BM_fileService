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

	engine.GET("/app/categories", middleware.Authenticate, gac.GetAllCategories, middleware.HandleErrors)
}

// GetAllCategories godoc
// @Summary Получение всех категорий
// @Description Возвращает все категории и лежащие в них папки
// @Tags Categories
// @Accept json
// @Produce json
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} []responses.Category "Категории данного аккаунта"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/categories [get]
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
