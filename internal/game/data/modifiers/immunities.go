package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var addImmunitiesID = uuid.MustParse("b4680929-95b5-599d-8a11-d4c08a87512f")

func AddImmunities(ids ...uuid.UUID) game.Modifier {
	return game.Modifier{
		ID:      uuid.MustParse("451cca1b-8b8a-5ee7-8e05-320d88e1a343"),
		GroupID: &addImmunitiesID,
		Name:    "Immunity",
		Show:    false,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				&addImmunitiesID,
				game.MutPriorityImmunity,
				game.ComposeAF(game.SourceFilter, game.ActiveFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					a.PushImmunities(ids...)
					return a
				},
			),
		},
	}
}
