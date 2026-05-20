package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var hardHeadedID = uuid.MustParse("71608f49-fb43-4ee6-a50e-b00c27a33f86")
var HardHeaded = game.Modifier{
	ID:          hardHeadedID,
	GroupID:     &hardHeadedID,
	Icon:        "hard_headed",
	Name:        "Hard Headed",
	Description: "Takes no recoil damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&hardHeadedID,
			game.MutPrioritySet,
			game.ComposeAF(game.SourceFilter),
			func(g game.Game, a game.Actor, ctx game.Context) game.Actor {
				a.RecoilMultiplier = 0.0
				return a
			},
		),
	},
}
