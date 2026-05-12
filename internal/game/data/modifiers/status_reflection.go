package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var statusReflectionID = uuid.MustParse("6e1f834f-96a1-4f85-b04f-02d293592ad3")
var StatusReflectionTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: statusReflectionID,
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

			modifier_tx, ok := g.GetModifierTxByID(*context.ModifierID)
			if !ok || modifier_tx.Context.SourceActorID == nil {
				return transactions
			}

			reflectDebuff := mutations.AddModifiers(false, modifier_tx.Mutation)
			removeDebuff := mutations.RemoveModifierTxByID(*context.ModifierID)
			source, ok := g.GetActorByID(*modifier_tx.Context.SourceActorID)
			if !ok {
				return transactions
			}

			targets := g.GetTargets(context)
			for _, target := range targets {
				if target.ID != source.ID {
					remove_ctx := game.MakeContextForActor(target)
					transactions = append(transactions, game.MakeTransaction(removeDebuff, remove_ctx))
				}
			}

			if len(transactions) > 0 {
				reflect_ctx := game.MakeContextForActor(source)
				transactions = append(transactions, game.MakeTransaction(reflectDebuff, reflect_ctx))
			}

			return transactions
		},
	},
}

var StatusReflection = game.Modifier{
	ID:          statusReflectionID,
	GroupID:     &statusReflectionID,
	Icon:        "status_reflection",
	Name:        "Debuff Reflection",
	Description: "On stat drop: reflect the effect back at the source.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&statusReflectionID),
	},
	Triggers: []game.Trigger{
		StatusReflectionTrigger,
	},
}
