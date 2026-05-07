package modifiers

import (
	"maps"
	"shinobi_showdown/internal/game"
	"slices"

	"github.com/google/uuid"
)

var poisonedID = uuid.New()

var PoisonedTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: poisonedID,
	On:         game.OnTurnEnd,
	Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
		active := game.TargetsAreActive(p, g, context)
		return active
	},
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			source, ok := g.GetSource(context)
			if !ok {
				return []game.GameTransaction{}
			}
			mut := game.RatioDamage(0.0625 * float64(source.PoisonedCounter-1))
			return []game.Transaction[game.GameMutation]{
				game.MakeTransaction(mut, context),
			}
		},
	},
}

var Poisoned game.Modifier = game.Modifier{
	ID:          poisonedID,
	GroupID:     &poisonedID,
	Name:        "Poisoned",
	Description: "On turn end: lose 1/16th * T HP.",
	Icon:        "poisoned",
	Show:        true,
	Status:      true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&poisonedID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.TargetFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				keys := maps.Keys(actor.AppliedModifiers)
				if slices.Contains(slices.Collect(keys), Guts.ID) {
					actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 1.5)
				}

				return actor
			},
		),
	},
	Triggers: []game.Trigger{
		PoisonedTrigger,
	},
}
