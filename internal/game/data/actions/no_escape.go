package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var NoEscape = MakeNoEscape()

func MakeNoEscape() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "No Escape",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Genjutsu,
		Description: "Target cannot escape.",
	})

	return game.Action{
		ID:              uuid.MustParse("8d567e50-0a59-4d5c-8e20-6da2698c05e9"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}

				for _, target := range g.GetTargets(context) {
					mut_ctx := game.MakeContextForActor(target)
					mut_ctx.ParentActorID = context.SourceActorID
					mut := mutations.AddModifiers(true, modifiers.SwitchLocked)
					tx := game.MakeTransaction(mut, mut_ctx)
					transactions = append(transactions, tx)
				}

				return transactions
			},
		},
	}
}
