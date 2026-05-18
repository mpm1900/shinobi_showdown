package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var FlyingRaijinStrike = MakeFlyingRaijin()

func MakeFlyingRaijin() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Flying Raijin: Strike",
		Description: "+2 priority. High critical hit chance.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(20),
		Jutsu:       game.Fuinjutsu,
	})

	config.CritChance = game.Ptr(getCriticalStage(1))

	return makeAttack(AttackConfig{
		ID:       uuid.MustParse("1a54031e-0ae6-49ed-b8b5-931c692417bf"),
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP2),
	})
}
