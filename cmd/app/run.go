package app

import (
	"app/config"
	http2 "app/controllers/http"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prolgrammer/BM_authService/pkg/middleware"
)

var (
	cfg             *config.Config
	loadFileUseCase usecases.LoadFileUseCases
)

func Run() {
	var err error
	cfg, err = config.New()
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}
	initUseCases()
	runServer()
}

func runServer() {
	router := gin.New()
	mw := middleware.NewMiddleware()

	http2.NewLoadFileController(router, loadFileUseCase, mw)
}

func initUseCases() {
	loadFileUseCase = usecases.NewLoadFileUseCase()
}
