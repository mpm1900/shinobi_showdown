package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var PatternBreak = MakePatternBreak()

func MakePatternBreak() game.Action {
	config := game.ActionConfig{
		Name:        "Pattern Break",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Ninjutsu,
		Description: "Target cannot repeat actions.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("1f17c177-bf08-451e-a052-c1e681e8d499"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.TargetLengthFilter(1),
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
					mod := modifiers.PatternBroke
					mutation := mutations.AddModifiers(false, mod)
					transaction := game.MakeTransaction(mutation, context)
					transactions = append(transactions, transaction)
				}

				return transactions
			},
		},
	}
}
