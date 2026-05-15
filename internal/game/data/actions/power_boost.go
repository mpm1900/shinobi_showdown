package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var PowerBoost = MakePowerBoost()

func MakePowerBoost() game.Action {
	config := game.ActionConfig{
		Name:        "Power Boost",
		Nature:      game.Ptr(game.NsSage),
		Jutsu:       game.Ninjutsu,
		Description: "Powers up target's attacks this turn.",
	}
	return game.Action{
		ID:              uuid.MustParse("fa0a4e99-9b26-5962-9ed0-fc88a6e73cb5"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.ComposeAF(game.TeamFilter, game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.TargetLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityP5,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.PowerBoosted)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
