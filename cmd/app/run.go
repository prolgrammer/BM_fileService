package app

import (
	"app/config"
	http2 "app/controllers/http"
	m "app/infrastructure/minio"
	"app/internal/usecases"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prolgrammer/BM_package/middleware"
	"net/http"
)

var (
	cfg             *config.Config
	loadFileUseCase usecases.LoadFileUseCases
	minioClient     *m.Client
)

func Run() {
	var err error
	cfg, err = config.New()
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	minioClient, err = m.NewClient(
		cfg.Minio.Endpoint,
		cfg.Minio.AccessKeyID,
		cfg.Minio.SecretAccessKey,
		cfg.Minio.BucketName,
		false)

	if err != nil {
		panic(fmt.Errorf("failed to initialize minio client: %v", err))
	}

	initUseCases()
	runServer()
}

func runServer() {
	router := gin.New()
	mw := middleware.NewMiddleware("someSuperStrongKey")

	http2.NewLoadFileController(router, loadFileUseCase, mw)

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	fmt.Printf("starting server at %s\n", address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}

func initUseCases() {
	loadFileUseCase = usecases.NewLoadFileUseCase(minioClient)
}
