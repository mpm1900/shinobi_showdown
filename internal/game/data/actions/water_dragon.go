package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var WaterDragon = MakeWaterDragon()

func MakeWaterDragon() game.Action {
	ID := uuid.MustParse("4e6e34e1-6dca-4739-b283-f5bd26f66628")

	config := game.ActionConfig{
		Name:        "Water Dragon",
		Description: "",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(80),
		Power:       game.Ptr(110),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
