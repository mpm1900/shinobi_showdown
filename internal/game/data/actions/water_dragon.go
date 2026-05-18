package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var WaterDragon = MakeWaterDragon()

func MakeWaterDragon() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Water Dragon",
		Description: "",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(80),
		Power:       game.Ptr(110),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("4e6e34e1-6dca-4739-b283-f5bd26f66628"),
		Config: config,
	})
}
