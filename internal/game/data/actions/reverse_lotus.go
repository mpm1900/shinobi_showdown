package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var ReverseLotus = MakeReverseLotus()

func MakeReverseLotus() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Reverse Lotus",
		Description: "If this attack fails or misses, the user loses half of their HP.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(50),
		Power:       game.Ptr(110),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("ca6c2a81-46bc-4c68-a986-39146bd5a52c"),
		Config: config,
		OnFailure: func(g game.Game, context, _ game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			ctx := game.MakeContextForActor(source)
			mut := game.RatioDamage(0.5)
			tx := game.MakeTransaction(mut, ctx)
			transactions = append(transactions, tx)

			return transactions
		},
	})
}
