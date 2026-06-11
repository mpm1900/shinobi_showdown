package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var AtomicDismantling = MakeAtomicDismantling()
var AtomicDismantlingCharge = MakeAtomicDismantlingCharge()

func MakeAtomicDismantlingCharge() game.Action {
	return makeSingleAction(
		uuid.MustParse("1d3c687f-d7a7-4f17-8f07-689764d9d263"),
		makeAttackConfig(game.ActionConfig{
			Name:        "Atomic Dismantling",
			Description: "Charges up, boosts chakra attack. Then attacks the next turn.",
			LogSuccess:  game.Ptr("$source$ is charging up."),
			Nature:      game.Ptr(game.NsParticle),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(130),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(0),
			Jutsu:       game.Ninjutsu,
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.ChakraAttackUpSource)
			transaction := game.MakeTransaction(mutation, context)
			transactions.Push(transaction)

			attack := mutations.QueueAction(AtomicDismantling.ID, context)
			transactions.Push(game.MakeTransaction(attack, context))

			return transactions.Build()
		},
	)
}

func MakeAtomicDismantling() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("d716d826-eb56-4022-a6aa-709091b5a4f0"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Atomic Dismantling",
			Description: "Charges up, boosts chakra attack. Then attacks the next turn.",
			Nature:      game.Ptr(game.NsParticle),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(130),
			Stat:        game.Ptr(game.StatChakraAttack),
			Jutsu:       game.Ninjutsu,
		}),
	})
}
