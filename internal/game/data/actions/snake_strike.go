package actions

import (
	"math/rand/v2"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var SnakeStrike = MakeSnakeStrike()

func MakeSnakeStrike() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Snake Strike",
		Description: "30% chance to paralyze, poison, or put target to sleep.",
		Nature:      game.Ptr(game.NsYin),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("62587c38-1644-4910-a2c0-c44a6b27c576"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				roll := rand.IntN(3)
				switch roll {
				case 0:
					transactions.Push(modifiers.ChanceParalysis(action_config, g, context, target, 30))
				case 1:
					transactions.Push(modifiers.ChancePoison(action_config, g, context, target, 30))
				case 2:
					transactions.Push(modifiers.ChanceSleep(action_config, g, context, target, 30))
				}
			}

			return transactions.Build()
		},
	})
}
