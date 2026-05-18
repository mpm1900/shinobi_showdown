package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var CreateRain = MakeCreateRain()

func MakeCreateRain() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Create Rain",
		Nature:      game.Ptr(game.NsWater),
		Jutsu:       game.Ninjutsu,
		Description: "Creates rain.",
	})

	return game.Action{
		ID:              uuid.MustParse("2e05db18-ac05-48fd-9cf4-b50fc6e5dbc3"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				return modifiers.ApplyRain(g, context)
			},
		},
	}
}
