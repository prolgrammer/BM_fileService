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

type CheckCategoryExistController struct {
	checkCategoryExistUseCase usecases.CheckCategoryExistUseCase
}

func NewCheckCategoryExistController(
	engine *gin.Engine,
	checkCategoryExistUseCase usecases.CheckCategoryExistUseCase,
	middleware middleware.Middleware,
) {

	checkCategoryExistController := &CheckCategoryExistController{
		checkCategoryExistUseCase: checkCategoryExistUseCase,
	}

	engine.GET("/app/category/exist", middleware.Authenticate, checkCategoryExistController.CheckCategoryExists, middleware.HandleErrors)
}

// CheckCategoryExists godoc
// @Summary      Проверка существования категории
// @Description  Проверяет, существует ли категория для указанного аккаунта
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        request body requests.Category true "Название категории"
// @Param        Authorization header string true "Токен доступа"
// @Success      200 {object} bool "Результат проверки (true - существует, false - не существует)"
// @Failure      400 {object} string "Некорректный формат запроса"
// @Failure      401 {object} string "Ошибка аутентификации"
// @Failure      500 {object} string "Внутренняя ошибка сервера"
// @Router       /app/category/exist [get]
func (c *CheckCategoryExistController) CheckCategoryExists(ctx *gin.Context) {
	fmt.Println("check category exist")

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

	response, err := c.checkCategoryExistUseCase.CheckCategoryExist(ctx, accountId.(string), req)
	if err != nil {
		wrappedError := fmt.Errorf("there was a problem during check category exist: %w", err)
		middleware.AddGinError(ctx, wrappedError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
