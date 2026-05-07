package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var ichirakuRamenID = uuid.MustParse("4186b062-e930-51cd-b838-12c618b790c0")

var IchirakuRamenTrigger game.Trigger = game.Trigger{
	ID:         uuid.MustParse("7d6c4561-b6e8-5107-8e8b-3d69b18cb93f"),
	ModifierID: ichirakuRamenID,
	On:         game.OnTurnEnd,
	Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
		return true
	},
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			context.TargetPositionIDs = []uuid.UUID{}
			context.TargetActorIDs = []uuid.UUID{*context.SourceActorID}
			mut := game.RatioHeal(0.0625)
			return []game.Transaction[game.GameMutation]{
				game.MakeTransaction(mut, context),
			}
		},
	},
}

var IchirakuRamen game.Modifier = game.Modifier{
	ID:          ichirakuRamenID,
	GroupID:     &ichirakuRamenID,
	Icon:        "ichiraku_ramen",
	Name:        "Ichiraku Ramen",
	Description: "End of turn: heal 1/16th HP.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&ichirakuRamenID),
	},
	Triggers: []game.Trigger{
		IchirakuRamenTrigger,
	},
}
