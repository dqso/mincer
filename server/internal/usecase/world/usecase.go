package usecase_world

import "context"

type Usecase struct {
	ncProducer ncProducer
	players    Players
}

type ncProducer interface {
	Pong(ctx context.Context, toClientID uint64, pong string) error
}

func NewUsecase(ncProducer ncProducer) *Usecase {
	return &Usecase{
		ncProducer: ncProducer,
	}
}
