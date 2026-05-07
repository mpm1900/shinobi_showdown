package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var sagesBlessingID = uuid.MustParse("9d937167-570a-4ded-8f50-7eee22c838c6")
var SagesBlessing = game.Modifier{
	ID:          sagesBlessingID,
	Name:        "Sage's Blessing",
	Description: "Doubles the chance of secondary effects.",
	Icon:        "sages_blessing",
	Show:        true,
	GroupID:     &sagesBlessingID,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&sagesBlessingID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				a.ModifierChanceMult += 1 // +100%/x2 increatse
				a.StatusChanceMult += 1
				return a
			},
		),
	},
}
