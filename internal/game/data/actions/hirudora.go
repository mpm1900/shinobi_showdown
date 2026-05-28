package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Hirudora = MakeHirudora()

func MakeHirudora() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "7th Gate: Hirudora",
		Description: "Lowers user's Attack by 2 stages.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(130),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("36c88bee-1eb0-4c55-bdc6-b704221ea846"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			mutation := mutations.AddModifiers(false, modifiers.AttackDown2Source)
			transaction := game.MakeTransaction(mutation, game.MakeContextForActor(source))
			transactions = append(transactions, transaction)

			return transactions
		},
	})
}
