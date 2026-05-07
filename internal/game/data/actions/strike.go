package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Strike = MakeStrike()

func MakeStrike() game.Action {
	ID := uuid.MustParse("418352fd-8d2a-42af-ba01-304a2ef44cda")

	config := game.ActionConfig{
		Name:        "Strike",
		Description: "",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(65),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(40),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
