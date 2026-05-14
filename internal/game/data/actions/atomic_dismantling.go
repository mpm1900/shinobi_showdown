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
	config := game.ActionConfig{
		Name:        "Atomic Dismantling",
		Description: "Charges up, boosts chakra attack. Then attacks the next turn.",
		LogSuccess:  game.Ptr("$source$ is charging up."),
		Nature:      game.Ptr(game.NsParticle),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(130),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}
	return game.Action{
		ID:              uuid.MustParse("1d3c687f-d7a7-4f17-8f07-689764d9d263"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.ChakraAttackUpSource)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				attack := mutations.QueueAction(AtomicDismantling.ID, context)
				transactions = append(transactions, game.MakeTransaction(attack, context))

				return transactions
			},
		},
	}
}

func MakeAtomicDismantling() game.Action {
	ID := uuid.MustParse("d716d826-eb56-4022-a6aa-709091b5a4f0")
	config := game.ActionConfig{
		Name:        "Atomic Dismantling",
		Description: "Charges up, boosts chakra attack. Then attacks the next turn.",
		Nature:      game.Ptr(game.NsParticle),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(130),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
