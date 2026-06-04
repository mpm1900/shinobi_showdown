package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var GreatFireball = MakeGreatFireball()

func MakeGreatFireball() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("57ddb3d7-0853-4a64-885b-18d93286c806"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Great Fireball",
			Description: "20% chance to burn target.",
			Nature:      game.Ptr(game.NsFire),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(95),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(60),
			Jutsu:       game.Ninjutsu,
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Push(modifiers.ChanceBurn(action_config, g, context, target, 20))
			}

			return transactions.Build()
		},
	})
}
