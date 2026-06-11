package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var ApplyPoison = MakeApplyPoison()

func MakeApplyPoison() game.Action {
	return makeSingleAction(
		uuid.MustParse("e5eb5e94-eea4-45b4-b033-518e696ca9a3"),
		makeStatusConfig(game.ActionConfig{
			Name:        "Apply Poison",
			Description: "Poisons target.",
			Nature:      game.Ptr(game.NsYin),
			Accuracy:    game.Ptr(90),
			Cost:        game.Ptr(30),
			Jutsu:       game.Ninjutsu,
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			action_config, _ := game.GetActiveActionConfig(g)

			for _, target := range g.GetTargets(context) {
				transactions.Concat(modifiers.ChancePoison(action_config, g, context, target, *action_config.Accuracy))
			}

			return transactions.Build()
		},
	)
}
