package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var HeavyPunch = MakeHeavyPunch()

func MakeHeavyPunch() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Heavy Punch",
		Description: "30% chance to paralyze the target.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("420bad58-1238-4124-909e-09ef76d743e8"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceParalysis(config, g, context, target, 30)...)
			}

			return transactions
		},
	})
}
