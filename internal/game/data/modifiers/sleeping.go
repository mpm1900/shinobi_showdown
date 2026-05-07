package modifiers

import (
	"maps"
	"shinobi_showdown/internal/game"
	"slices"

	"github.com/google/uuid"
)

var sleepingID = uuid.MustParse("b8778c2e-12e3-410f-b857-b304f2f7dd1e")

var Sleeping game.Modifier = game.Modifier{
	ID:          sleepingID,
	GroupID:     &sleepingID,
	Name:        "Sleeping",
	Description: "Sleeping",
	Icon:        "sleeping",
	Show:        true,
	Status:      true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&sleepingID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				keys := maps.Keys(actor.AppliedModifiers)
				if slices.Contains(slices.Collect(keys), Guts.ID) {
					actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 1.5)
				}

				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
