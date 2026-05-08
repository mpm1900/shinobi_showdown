package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var waterPrisonID = uuid.MustParse("98ed6d8b-f8c1-4018-b40c-f6d7757fa78c")

var WaterPrison = game.Modifier{
	ID:          waterPrisonID,
	GroupID:     &waterPrisonID,
	Name:        "Water Prison",
	Icon:        "water_prison",
	Description: "Cannot switch out. On turn end: lose 1/8th HP.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&waterPrisonID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.ParentFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.SwitchLocked = true
				return actor
			},
		),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: waterPrisonID,
			On:         game.OnTurnEnd,
			Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
				active := game.TargetsAreActive(p, g, context)
				return active
			},
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					mut := game.RatioDamage(0.125)
					return []game.Transaction[game.GameMutation]{
						game.MakeTransaction(mut, context),
					}
				},
			},
		},
	},
}
