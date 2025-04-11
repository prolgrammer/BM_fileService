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

type deleteCategoryController struct {
	deleteCategoryUseCase usecases.DeleteCategoryUseCase
}

func NewDeleteCategoryController(
	engine *gin.Engine,
	deleteCategoryUseCase usecases.DeleteCategoryUseCase,
	middleware middleware.Middleware,
) {
	dc := &deleteCategoryController{
		deleteCategoryUseCase: deleteCategoryUseCase,
	}

	engine.POST("/app/category/delete", middleware.Authenticate, dc.Delete, middleware.HandleErrors)
}

// Delete godoc
// @Summary Удаление категории
// @Description Удаляет категорию, ее папки и файлы
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body requests.Category true "Название категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} string "Результат удаление"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/category/delete [post]
func (dc *deleteCategoryController) Delete(ctx *gin.Context) {
	fmt.Println("delete category")

	accountId, exist := ctx.Get("account_id")
	if !exist {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.Category
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w:%w", err, e.ErrDataBindError)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	err := dc.deleteCategoryUseCase.DeleteCategory(ctx, accountId.(string), req)
	if err != nil {
		middleware.AddGinError(ctx, fmt.Errorf("something went wrong while delete category: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"answer": "category delete successful"})
}
