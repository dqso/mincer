package entity

type Horn interface {
	OnPlayerWasted(playerID uint64, playerClass Class, killerID uint64, killerClass Class)
	OnPlayerAttacked(id uint64, directionAim float64)
	SpawnPlayer(player Player)
	SetPlayerStats(id uint64, stats PlayerStats)
	SetPlayerHP(id uint64, hp int32)
	SetPlayerPosition(id uint64, position Point)
	SetPlayerWeapon(id uint64, w Weapon)
	SetProjectilePosition(id uint64, position Point)
}
