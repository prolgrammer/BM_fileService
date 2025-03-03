package app

import (
	"app/config"
	http2 "app/controllers/http"
	m "app/infrastructure/minio"
	"app/infrastructure/mongo"
	"app/internal/repositories"
	"app/internal/usecases"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prolgrammer/BM_package/middleware"
	"net/http"
)

var (
	cfg         *config.Config
	minioClient *m.Client
	mongoClient *mongo.Client

	categoryRepository repositories.CategoryRepository
	folderRepository   repositories.FolderRepository

	createCategoryUseCase      usecases.CreateCategoryUseCase
	getCategoryUseCase         usecases.GetCategoryUseCase
	getAllCategoriesUseCase    usecases.GetAllCategoryUseCase
	deleteCategoryUseCase      usecases.DeleteCategoryUseCase
	checkCategoryExistsUseCase usecases.CheckCategoryExistUseCase

	createFolderUseCase      usecases.CreateFolderUseCase
	getFolderUseCase         usecases.SelectFolderUseCase
	getAllFoldersUseCase     usecases.SelectFoldersUseCase
	deleteFolderUseCase      usecases.DeleteFolderUseCase
	checkFolderExistsUseCase usecases.CheckFolderExistUseCase

	loadFileUseCase usecases.LoadFileUseCases
)

func Run() {
	var err error
	cfg, err = config.New()
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	err = initPackages(cfg)
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	minioClient, err = m.NewClient( //TODO организовать получше
		cfg.Minio.Endpoint,
		cfg.Minio.AccessKeyID,
		cfg.Minio.SecretAccessKey,
		cfg.Minio.BucketName,
		false)

	if err != nil {
		panic(fmt.Errorf("failed to initialize minio client: %v", err))
	}

	initRepositories()
	initUseCases()
	runServer()
}

func initPackages(cfg *config.Config) error {
	var err error

	fmt.Printf("Starting mongo client\n")

	mongoClient, err = mongo.NewClient(cfg.Mongo)
	if err != nil {
		return err
	}

	err = mongoClient.MigrateUp()
	if err != nil {
		if !errors.Is(err, mongo.ErrNoChange) {
			fmt.Printf("failed to migrate up mongo client: %v\n", err)
			return err
		}
		fmt.Printf("mongo has the latest version. nothing to migrate\n")
	}

	return nil
}

func initRepositories() {
	collection := mongoClient.Database.Collection("categories")

	categoryRepository = repositories.NewCategoryDataRepository(collection)
	folderRepository = repositories.NewFolderMongoRepository(collection)
}

func initUseCases() {
	createCategoryUseCase = usecases.NewCreateCategoryUseCase(categoryRepository)
	getCategoryUseCase = usecases.NewGetCategory(categoryRepository)
	getAllCategoriesUseCase = usecases.NewGetAllCategory(categoryRepository)
	deleteCategoryUseCase = usecases.NewDeleteCategoryUseCase(categoryRepository)
	checkCategoryExistsUseCase = usecases.NewCheckCategoryExist(categoryRepository)

	createFolderUseCase = usecases.NewCreateFolderUseCase(categoryRepository, folderRepository)
	getFolderUseCase = usecases.NewSelectFolderUseCase(categoryRepository, folderRepository)
	getAllFoldersUseCase = usecases.NewSelectFoldersUseCase(categoryRepository, folderRepository)
	deleteFolderUseCase = usecases.NewDeleteFolderUseCase(categoryRepository, folderRepository)
	checkFolderExistsUseCase = usecases.NewCheckFolderExistUseCase(categoryRepository, folderRepository)

	loadFileUseCase = usecases.NewLoadFileUseCase(minioClient)
}

func runServer() {
	router := gin.New()
	mw := middleware.NewMiddleware("someSuperStrongKey")

	http2.NewCreateCategoryController(router, createCategoryUseCase, mw)
	http2.NewGetCategoryController(router, getCategoryUseCase, mw)
	http2.NewGetAllCategoriesUseCases(router, getAllCategoriesUseCase, mw)
	http2.NewDeleteCategoryController(router, deleteCategoryUseCase, mw)
	http2.NewCheckCategoryExistController(router, checkCategoryExistsUseCase, mw)

	http2.NewCreateFolderController(router, createFolderUseCase, mw)
	http2.NewGetFolderController(router, getFolderUseCase, mw)
	http2.NewGetAllFoldersController(router, getAllFoldersUseCase, mw)
	http2.NewDeleteFolderController(router, deleteFolderUseCase, mw)
	http2.NewCheckFolderExistController(router, checkFolderExistsUseCase, mw)

	http2.NewLoadFileController(router, loadFileUseCase, mw)

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	fmt.Printf("starting server at %s\n", address)

	err := http.ListenAndServe(address, router)
	if err != nil {
		panic(err)
	}
}
