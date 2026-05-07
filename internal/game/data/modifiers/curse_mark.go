package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var cmoChakraID = uuid.MustParse("f0dcf295-6321-5694-a55e-b3bf8afcc508")
var cmoSpeedID = uuid.MustParse("f5655358-ad5a-57ab-bf0b-83695b8b5185")
var cmoStrengthID = uuid.MustParse("4d1ea730-f18b-5e72-843c-fe5845fc9acf")

var CurseMarkOfChakra game.Modifier = game.Modifier{
	ID:          cmoChakraID,
	GroupID:     &cmoChakraID,
	Icon:        "cmo_chakra",
	Name:        "Curse Mark of Chakra",
	Description: "Chakra Attack x1.5; one action only.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&cmoChakraID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.ActionLocked = true
				actor.Stats[game.StatChakraAttack] = game.Round(float64(actor.Stats[game.StatChakraAttack]) * 1.5)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}

var CurseMarkOfSpeed game.Modifier = game.Modifier{
	ID:          cmoSpeedID,
	GroupID:     &cmoSpeedID,
	Icon:        "cmo_speed",
	Name:        "Curse Mark of Speed",
	Description: "Speed x1.5; one action only.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&cmoSpeedID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.ActionLocked = true
				actor.Stats[game.StatSpeed] = game.Round(float64(actor.Stats[game.StatSpeed]) * 1.5)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}

var CurseMarkOfStrength game.Modifier = game.Modifier{
	ID:          cmoStrengthID,
	GroupID:     &cmoStrengthID,
	Icon:        "cmo_strength",
	Name:        "Curse Mark of Strength",
	Description: "Attack x1.5; one action only.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&cmoStrengthID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.ActionLocked = true
				actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 1.5)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
