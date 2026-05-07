package modifiers

import (
	"maps"
	"shinobi_showdown/internal/game"
	"slices"

	"github.com/google/uuid"
)

var burnedID = uuid.MustParse("497240d3-0d1a-5194-b475-62e418ad94eb")

var BurnedTrigger game.Trigger = game.Trigger{
	ID:         uuid.MustParse("6d3a4475-cca7-57e1-b606-3c4515c955ec"),
	ModifierID: burnedID,
	On:         game.OnTurnEnd,
	Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
		active := game.TargetsAreActive(p, g, context)
		return active
	},
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			mut := game.RatioDamage(0.0625)
			return []game.Transaction[game.GameMutation]{
				game.MakeTransaction(mut, context),
			}
		},
	},
}

var Burned game.Modifier = game.Modifier{
	ID:          burnedID,
	GroupID:     &burnedID,
	Name:        "Burned",
	Description: "Burned: Attack x0.5. On turn end: lose 1/16th HP.",
	Icon:        "burned",
	Show:        true,
	Status:      true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&burnedID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.TargetFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				keys := maps.Keys(actor.AppliedModifiers)
				if slices.Contains(slices.Collect(keys), Guts.ID) {
					actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 1.5)
				} else {
					actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 0.5)
				}

				return actor
			},
		),
	},
	Triggers: []game.Trigger{
		BurnedTrigger,
	},
}
