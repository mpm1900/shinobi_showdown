package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Strike = MakeStrike()

func MakeStrike() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Strike",
		Description: "",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(65),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("418352fd-8d2a-42af-ba01-304a2ef44cda"),
		Config: config,
	})
}
