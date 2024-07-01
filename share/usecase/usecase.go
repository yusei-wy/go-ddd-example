package usecase

import (
	userDomain "go_ddd_example/feature/user/domain"
	userInfra "go_ddd_example/feature/user/infrastructure"
	userUseCase "go_ddd_example/feature/user/usecase"

	"github.com/jmoiron/sqlx"
)

type UseCaseFacade struct {
	UserUseCase userUseCase.UserUseCase
}

func NewUseCaseFacade(
	db *sqlx.DB,
) UseCaseFacade {
	userRepository := userInfra.NewPsQlUserRepository(db)
	userService := userDomain.NewUserServiceImpl(userRepository)
	userUseCase := userUseCase.NewUserUseCaseImpl(userRepository, userService)

	return UseCaseFacade{
		UserUseCase: userUseCase,
	}
}
