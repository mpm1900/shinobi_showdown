package actions

import (
	"fmt"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Fireball = MakeFireball()

func MakeFireball() game.Action {
	chance := 10
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("aaf5174b-f386-54b1-84c4-0c062937c770"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Fireball",
			Description: fmt.Sprintf("%d%% chance to burn target.", chance),
			Nature:      game.Ptr(game.NsFire),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(55),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(40),
			Jutsu:       game.Ninjutsu,
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			for _, target := range g.GetTargets(context) {
				transactions.Push(modifiers.ChanceBurn(action_config, g, context, target, chance))
			}

			return transactions.Build()
		},
	})
}
