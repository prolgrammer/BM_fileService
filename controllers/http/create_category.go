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

type createCategoryController struct {
	createCategory usecases.CreateCategoryUseCase
}

func NewCreateCategoryController(
	engine *gin.Engine,
	createCategory usecases.CreateCategoryUseCase,
	middleware middleware.Middleware,
) {
	cc := &createCategoryController{
		createCategory: createCategory,
	}

	engine.POST("/app/category", middleware.Authenticate, cc.CreateCategory, middleware.HandleErrors)
}

// CreateCategory godoc
// @Summary Создание категории пользователя
// @Description Создает категорию для пользователя по ID переданного через токен
// @Tags Categories
// @Accept json
// @Produce json
// @Param request body requests.Category true "Название категории"
// @Param Authorization header string true "Токен доступа"
// @Success 200 {object} string "Результат создания"
// @Failure 400 {object} string "Некорректный формат запроса"
// @Failure 401 {object} string "Ошибка аутентификации"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /app/category [post]
func (cc *createCategoryController) CreateCategory(ctx *gin.Context) {
	fmt.Println("create category")

	accountId, exist := ctx.Get("account_id")
	if !exist {
		wrappedError := fmt.Errorf("%w", e.ErrAuthenticated)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	var req requests.Category
	if err := ctx.ShouldBindJSON(&req); err != nil {
		wrappedError := fmt.Errorf("%w: %w", e.ErrDataBindError, err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	err := cc.createCategory.CreateCategory(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during create category: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"answer": "category create successful"})

}
