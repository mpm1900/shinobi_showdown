package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var allyGuardID = uuid.MustParse("1dc88d90-8f98-4a87-bfa8-7ab5a39165fd")
var AllyGuard = game.Modifier{
	ID:          allyGuardID,
	GroupID:     &allyGuardID,
	Name:        "Ally Guard",
	Icon:        "ally_guard",
	Show:        true,
	Description: "Ally shinobi take 25% less damage.",
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&allyGuardID,
			game.MutPriorityDefault,
			game.ComposeAF(game.OtherFilter, game.TeamFilter),
			func(g game.Game, a game.Actor, ctx game.Context) game.Actor {
				a.DamageReduction[game.ChakraAttack] *= 0.75
				a.DamageReduction[game.Attack] *= 0.75
				return a
			},
		),
	},
}
