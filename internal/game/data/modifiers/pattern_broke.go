package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var patternBrokeID = uuid.MustParse("0ce5b5e0-5865-473b-a9cb-4b574f775ef8")

var PatternBroke = game.Modifier{

	ID:          uuid.New(),
	Name:        "Pattern Broke",
	Description: "Cannot repeat actions.",
	Icon:        "pattern_broke",
	Show:        true,
	GroupID:     &patternBrokeID,
	Duration:    5,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&patternBrokeID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.TargetFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				if a.LastUsedActionTX == nil {
					return a
				}

				for i, _ := range a.Actions {
					if a.Actions[i].ID == a.LastUsedActionTX.Mutation.ID {
						a.Actions[i].State.Disabled = true
					}
				}
				return a
			},
		),
	},
}
