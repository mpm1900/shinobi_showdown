package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var shadowCageID = uuid.MustParse("4494f97f-ac70-4fa5-b599-633aa13536f7")

var ShadowCage game.Modifier = game.Modifier{
	ID:                shadowCageID,
	GroupID:           &shadowCageID,
	Name:              "Shadow Cage",
	Icon:              "switch_locked",
	Description:       "Cannot switch out",
	ParentDescription: "Enemy shinobi cannot switch out",
	Show:              true,
	Duration:          game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&shadowCageID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.OtherTeamFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.SwitchLocked = true
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}

var switchLockedID = uuid.MustParse("be4bea95-9415-47f0-add6-9d9f07095ef2")
var SwitchLocked game.Modifier = game.Modifier{
	ID:          switchLockedID,
	GroupID:     &switchLockedID,
	Name:        "Switch Locked",
	Icon:        "switch_locked",
	Description: "Cannot switch out",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&switchLockedID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.SwitchLocked = true
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
