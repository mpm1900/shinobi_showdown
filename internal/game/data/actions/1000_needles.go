package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var OneThousandNeedles = MakeOneThousandNeedles()

func MakeOneThousandNeedles() game.Action {
	return makeSingleAction(
		uuid.MustParse("58c829b9-aa81-4a44-84c7-73cf08501e48"),
		makeAttackConfig(game.ActionConfig{
			Name:        "1000 Needles",
			Description: "Hits 2-5 times. High critical chance.",
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(25),
			Stat:        game.Ptr(game.StatChakraAttack),
			Nature:      game.Ptr(game.NsIce),
			Cost:        game.Ptr(50),
			Jutsu:       game.Ninjutsu,
			CritStage:   game.Ptr(1),
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			action_config, _ := game.GetActiveActionConfig(g)
			damage_config := game.NewDamageConfig(game.RandomDamageFactor())
			damage_config = makeRepeats(damage_config, 2, 5, g, context)
			transactions.Concat(game.ResolveDamageCore(action_config, damage_config, g, context))

			return transactions.Build()
		},
	)
}
