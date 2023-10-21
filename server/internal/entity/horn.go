package entity

type Horn interface {
	SpawnPlayer(player Player)
	SetPlayerStats(id uint64, stats PlayerStats)
	SetPlayerHP(id uint64, hp int32)
	SetPlayerPosition(id uint64, position Point)
}
