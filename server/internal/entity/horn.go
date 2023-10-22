package entity

type Horn interface {
	OnPlayerWasted(id uint64, killer uint64)
	OnPlayerAttacked(id uint64, directionAim float64)
	SpawnPlayer(player Player)
	SetPlayerStats(id uint64, stats PlayerStats)
	SetPlayerHP(id uint64, hp int32)
	SetPlayerPosition(id uint64, position Point)
	SetPlayerWeapon(id uint64, w Weapon)
	SetProjectilePosition(id uint64, position Point)
}
