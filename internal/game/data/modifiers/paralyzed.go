package modifiers

import (
	"maps"
	"shinobi_showdown/internal/game"
	"slices"

	"github.com/google/uuid"
)

var paralysisID = uuid.MustParse("b269dcbc-a6a0-5ea3-9ed5-6e8b5fee8024")

var Paralysis game.Modifier = game.Modifier{
	ID:          paralysisID,
	GroupID:     &paralysisID,
	Name:        "Paralyzed",
	Description: "Paralyzed: Speed x0.25.",
	Icon:        "paralyzed",
	Show:        true,
	Status:      true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&paralysisID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.TargetFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				keys := maps.Keys(actor.AppliedModifiers)
				if slices.Contains(slices.Collect(keys), Guts.ID) {
					actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 1.5)
				}

				actor.Stats[game.StatSpeed] = game.Round(float64(actor.Stats[game.StatSpeed]) * 0.25)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
