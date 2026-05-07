package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Coercion = MakeCoercion()

func MakeCoercion() game.Action {
	config := game.ActionConfig{
		Name:        "Coercion",
		Nature:      game.Ptr(game.NsYin),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Genjutsu,
		Description: "Forces the target to use only their last used action.",
	}
	return game.Action{
		ID:              uuid.MustParse("06840403-52cc-4e8a-95eb-318cf012e634"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPrioritySlow,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				targets := g.GetTargets(context)
				for _, target := range targets {
					mut_ctx := game.MakeContextForActor(target)
					if target.LastUsedActionTX != nil {
						mutation := mutations.AddModifiers(true, modifiers.Coerced(target.LastUsedActionTX.Mutation.ID))
						transaction := game.MakeTransaction(mutation, mut_ctx)
						transactions = append(transactions, transaction)
					}
				}

				return transactions
			},
		},
	}
}
