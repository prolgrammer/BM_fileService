package usecases

import (
	"app/controllers/requests"
	"context"
	"fmt"
	"io"
	"os"
)

type LoadFileUseCases interface {
	LoadFile(ctx context.Context, req requests.LoadFile) error
}

type loadFileUseCase struct{}

func NewLoadFileUseCase() LoadFileUseCases {
	return &loadFileUseCase{}
}

func (l loadFileUseCase) LoadFile(ctx context.Context, req requests.LoadFile) error {
	var err error

	for i := range req.Files {
		data, err := io.ReadAll(req.Files[i])
		if err != nil {
			fmt.Printf("Error while read file: %v", err)
			return err
		}

		file, err := os.Create("../testdir")
		if err != nil {
			fmt.Printf("Error while create file: %v", err)
			return err
		}

		_, err = file.Write(data)
		if err != nil {
			fmt.Printf("Error while write file: %v", err)
			return err
		}
	}

	return err
}
