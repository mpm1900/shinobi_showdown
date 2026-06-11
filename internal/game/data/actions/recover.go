package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Recover = MakeRecover()

func MakeRecover() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Recover",
		Description: "Heals the user for up to half of their HP.",
		Nature:      game.Ptr(game.NsYang),
		Cost:        game.Ptr(30),
		Jutsu:       game.Senjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("c0756ddc-2611-5eef-82cc-c2bc03f9f01c"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.ActiveFilter, game.TeamFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				action_config, _ := game.GetActiveActionConfig(g)
				heal := game.NewHeal(action_config, 0.5)
				transactions = append(
					transactions,
					game.MakeTransaction(heal, context),
				)

				return transactions
			},
		},
	}
}
