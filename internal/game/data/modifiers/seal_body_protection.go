package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var sobpID = uuid.MustParse("096e2442-b231-53da-9892-91b0dea908b9")
var SealOfBodyProtectionTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: sobpID,
	On:         game.OnModifierAdd,
	Check: game.ComposeTF(
		game.Match__TargetActor_SourceActor,
		game.Modifier__IsOneOf(
			AttackDownID,
			DefenseDownID,
			ChakraAttackDownID,
			ChakraDefenseDownID,
			SpeedDownID,
			AccuracyDownID,
			EvasionDownID,
		),
	),
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			if context.ModifierID == nil {
				return transactions
			}

			removeDebuff := mutations.RemoveModifierTxByID(*context.ModifierID)
			targets := g.GetTargets(context)
			for _, target := range targets {
				consumeCtx := game.MakeContextForActor(target)
				transactions = append(transactions, game.MakeTransaction(mutations.ConsumeItem, consumeCtx))
				transactions = append(transactions, game.MakeTransaction(removeDebuff, context))
			}

			return transactions
		},
	},
}

var SealOfBodyProtection game.Modifier = game.Modifier{
	ID:          sobpID,
	GroupID:     &sobpID,
	Icon:        "seal_up",
	Name:        "Seal of Body Protection",
	Description: "On stat drop: remove it and break this seal.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&sobpID),
	},
	Triggers: []game.Trigger{
		SealOfBodyProtectionTrigger,
	},
}
