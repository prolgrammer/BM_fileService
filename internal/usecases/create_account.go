package usecases

import (
	"app/internal/repositories"
	"context"
	"fmt"
)

type createAccountUseCase struct {
	accountRepository repositories.AccountRepository
}

type CreateAccountUseCase interface {
	CreateAccount(ctx context.Context, accountId string) error
}

func NewCreateAccountUseCase(accountRepository repositories.AccountRepository) CreateAccountUseCase {
	return &createAccountUseCase{
		accountRepository: accountRepository,
	}
}

func (uc *createAccountUseCase) CreateAccount(ctx context.Context, accountId string) error {
	exists, err := uc.accountRepository.CheckAccountExists(ctx, accountId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if exists {
		return ErrAccountAlreadyExists
	}

	return uc.accountRepository.CreateAccount(ctx, accountId)
}
