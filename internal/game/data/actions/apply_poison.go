package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var ApplyPoison = MakeApplyPoison()

func MakeApplyPoison() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Apply Poison",
		Description: "Poisons target.",
		Nature:      game.Ptr(game.NsYin),
		Accuracy:    game.Ptr(90),
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("e5eb5e94-eea4-45b4-b033-518e696ca9a3"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := game.NewTransactionBuilder()

				for _, target := range g.GetTargets(context) {
					transactions.Push(modifiers.ApplyPoison(config, g, target, context))
				}

				return transactions.Build()
			},
		},
	}
}
