package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var FlyingRaijinStrike = MakeFlyingRaijin()

func MakeFlyingRaijin() game.Action {
	ID := uuid.MustParse("1a54031e-0ae6-49ed-b8b5-931c692417bf")
	config := game.ActionConfig{
		Name:        "Flying Raijin: Strike",
		Description: "+2 priority. High critical hit chance.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(1),
		Jutsu:       game.Fuinjutsu,
		CritChance:  game.Ptr(getCriticalStage(1)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:       ID,
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP2),
	})
}
