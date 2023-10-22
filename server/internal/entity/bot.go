package entity

type Bot interface {
	Player
	GetPlayer() Player
	BotID() uint32
}
