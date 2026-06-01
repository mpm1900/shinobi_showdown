package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var speedBoostID = uuid.MustParse("557f72c6-142b-5195-a24e-09851e14223b")
var SpeedBoostTrigger game.Trigger = game.Trigger{
	ID:         uuid.MustParse("a6cd5955-5996-5fa1-844d-54a1cfdf243e"),
	ModifierID: speedBoostID,
	On:         game.OnTurnEnd,
	Check: func(g1, g2 game.Game, ctx game.Context, t game.Transaction[game.Modifier]) bool {
		return true
	},
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			mutation := mutations.AddModifiers(false, SpeedUpSource)
			transactions = append(transactions, game.MakeTransaction(mutation, context))
			return transactions
		},
	},
}

var SpeedBoost game.Modifier = game.Modifier{
	ID:          speedBoostID,
	GroupID:     &speedBoostID,
	Icon:        "speed_up",
	Name:        "Speed Boost",
	Description: "On turn end: gain Speed Up.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&speedBoostID),
	},
	Triggers: []game.Trigger{
		SpeedBoostTrigger,
	},
}
