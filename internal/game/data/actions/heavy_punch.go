package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var HeavyPunch = MakeHeavyPunch()

func MakeHeavyPunch() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("420bad58-1238-4124-909e-09ef76d743e8"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Heavy Punch",
			Description: "30% chance to paralyze the target.",
			Nature:      game.Ptr(game.NsTai),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(80),
			Stat:        game.Ptr(game.StatAttack),
			Cost:        game.Ptr(0),
			Jutsu:       game.Ninjutsu,
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceParalysis(action_config, g, context, target, 30))
			}

			return transactions.Build()
		},
	})
}
