package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Disable = MakeDisable()

func MakeDisable() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Disable",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Genjutsu,
		Description: "Disables the target's last used action.",
	})

	return game.Action{
		ID:              uuid.MustParse("5cf69985-6785-56a6-b879-e02cb6207960"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
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
						mutation := mutations.AddModifiers(true, modifiers.Disabled(target.LastUsedActionTX.Mutation.ID))
						transaction := game.MakeTransaction(mutation, mut_ctx)
						transactions = append(transactions, transaction)
					}
				}

				return transactions
			},
		},
	}
}
