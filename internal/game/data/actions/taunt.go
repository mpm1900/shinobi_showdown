package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Taunt = MakeTaunt()

func MakeTaunt() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Taunt",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Ninjutsu,
		Description: "Forces target to use only attacking moves.",
	})

	return game.Action{
		ID:              uuid.MustParse("c62f29ad-2f3e-5e5e-b045-bb0ed58837bc"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				targets := g.GetTargets(context)
				for _, target := range targets {
					context.ParentActorID = &target.ID
					mutation := mutations.AddModifiers(false, modifiers.Taunted)
					transaction := game.MakeTransaction(mutation, context)
					transactions = append(transactions, transaction)
				}

				return transactions
			},
		},
	}
}
