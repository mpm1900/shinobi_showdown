package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var sompID = uuid.MustParse("9bbeb2b1-ac4c-4e5b-bc62-2f631d3a5bd6")
var SealOfMindProtectionTrigger game.Trigger = game.Trigger{
	ID:         uuid.MustParse("5e8ddef1-6892-4422-bb39-e8148cf33ba4"),
	ModifierID: sompID,
	On:         game.OnModifierAdd,
	Check: game.ComposeTF(
		game.Match__TargetActor_SourceActor,
		game.Modifier__IsOneOf(
			*Taunted.GroupID,
			*Disabled(uuid.Nil).GroupID,
			*Coerced(uuid.Nil).GroupID,
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

var SealOfMindProtection game.Modifier = game.Modifier{
	ID:          sompID,
	GroupID:     &sompID,
	Icon:        "seal_up",
	Name:        "Seal of Mind Protection",
	Description: "On mental debuff: remove it and break this seal.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&sompID),
	},
	Triggers: []game.Trigger{
		SealOfMindProtectionTrigger,
	},
}
