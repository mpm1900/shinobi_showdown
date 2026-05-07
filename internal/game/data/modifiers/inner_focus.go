package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var innerFocusID = uuid.MustParse("20e9e34d-f9cd-48fb-bd12-eb408c58126c")

var InnerFocus = game.Modifier{
	ID:          innerFocusID,
	GroupID:     &innerFocusID,
	Icon:        "inner_focus",
	Name:        "Inner Focus",
	Description: "Cannot be intimidated or stunned.",
	Show:        true,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&innerFocusID,
			game.MutPriorityImmunity,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				a.PushImmunities(*Stunned.GroupID, *Intimidate.GroupID)
				return a
			},
		),
	},
}
