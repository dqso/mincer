package usecase_token

import "context"

type Usecase struct {
	config     config
	repository repository
}

type config interface {
	NCPrivateKey() []byte
	NCPort() int
}

type repository interface {
	AcquireClientID(ctx context.Context) (uint64, error)
}

func NewUsecase(config config, repository repository) *Usecase {
	return &Usecase{
		config:     config,
		repository: repository,
	}
}
