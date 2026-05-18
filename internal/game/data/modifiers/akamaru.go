package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var akamaruID = uuid.MustParse("cf31744b-d2aa-4ced-a39e-f2e57022d2a8")
var Akamaru = game.Modifier{
	ID:          akamaruID,
	GroupID:     &akamaruID,
	Icon:        "akamaru",
	Name:        "Akamaru",
	Description: "Is good boy.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&akamaruID),
	},
	Triggers: []game.Trigger{},
}

var mansBestFriendID = uuid.MustParse("b5a62630-a563-41cb-87f3-52a678987f56")
var MansBestFriend = game.Modifier{
	ID:          mansBestFriendID,
	GroupID:     &mansBestFriendID,
	Icon:        "akamaru",
	Name:        "Man's Best Friend",
	Description: "Doubled attack with Akamaru.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&mansBestFriendID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				if actor.Item == nil || actor.Item.ID != akamaruID {
					return actor
				}

				actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 2)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
