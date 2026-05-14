package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var defiantID = uuid.MustParse("4424a645-5f93-4d47-aa4b-3a44469b4b9b")
var defiantTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: defiantID,
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

			attackUp := mutations.AddModifiers(false, AttackUp2Target)
			transactions = append(transactions, game.MakeTransaction(attackUp, context))

			return transactions
		},
	},
}

var Defiant = game.Modifier{
	ID:          defiantID,
	GroupID:     &defiantID,
	Icon:        "defiant",
	Name:        "Defiant",
	Description: "On stat drop: gain +2 Attack Up",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&defiantID),
	},
	Triggers: []game.Trigger{
		defiantTrigger,
	},
}
