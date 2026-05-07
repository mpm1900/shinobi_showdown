package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var burdenOfPainID = uuid.MustParse("606b4892-56fa-58c4-af72-63f7ea3ef669")

var BurdenOfPain game.Modifier = game.Modifier{
	ID:                burdenOfPainID,
	GroupID:           &burdenOfPainID,
	Icon:              "std_attack",
	Name:              "Burden of Pain",
	Description:       "Attack x0.75.",
	ParentDescription: "Other shinobi: Attack x0.75.",
	Show:              true,
	Duration:          game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&burdenOfPainID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.OtherFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 0.75)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
