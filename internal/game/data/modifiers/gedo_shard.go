package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var gedoShardID = uuid.MustParse("7dd1029f-aac8-5802-9fae-3873e507b57d")

var GedoShardTrigger game.Trigger = game.Trigger{
	ID:         uuid.MustParse("8ddae3f7-dffd-5a57-b3fc-eb3b50e01f96"),
	ModifierID: gedoShardID,
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
			mut := game.RatioDamage(0.1)
			return []game.Transaction[game.GameMutation]{
				game.MakeTransaction(mut, context),
			}
		},
	},
}

var GedoShard game.Modifier = game.Modifier{
	ID:          gedoShardID,
	GroupID:     &gedoShardID,
	Icon:        "gedo_shard",
	Name:        "Gedo Shard",
	Description: "Deal 30% more damage. On turn end: lose 1/10th HP.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&gedoShardID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.DamageMultipliers[game.Attack] += 0.3
				actor.DamageMultipliers[game.ChakraAttack] += 0.3
				return actor
			},
		),
	},
	Triggers: []game.Trigger{
		GedoShardTrigger,
	},
}
