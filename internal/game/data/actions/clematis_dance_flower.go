package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ClematisDanceFlower = MakeClematisDanceFlower()

func MakeClematisDanceFlower() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Clematis Dance: Flower",
		Description: "User's Defense and Chakra Defense are lowered",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("8af50cf3-49f3-4529-94f7-465ffa144f53"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			source, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			mutation := mutations.AddModifiers(false, modifiers.DefenseDownSource, modifiers.ChakraDefenseDownSource)
			transaction := game.MakeTransaction(mutation, game.MakeContextForActor(source))
			transactions.PushOne(transaction)

			return transactions.Build()
		},
	})
}
