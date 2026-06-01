package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var poisonInfusedID = uuid.MustParse("27591068-a257-5554-b654-60d8e46e30f9")
var PoisonInfusedTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: poisonInfusedID,
	On:         game.OnDamageReceive,
	Check:      game.Match__TargetActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: 0,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			return ChancePoison(game.ActionConfig{}, g, game.MakeContextForActor(source), source, 50)
		},
	},
}

var PoisonInfused game.Modifier = game.Modifier{
	ID:          poisonInfusedID,
	GroupID:     &poisonInfusedID,
	Name:        "Poison Infused",
	Icon:        "poison_infused",
	Description: "On damage recieved: source has a 50% chance of being poisoned.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&poisonInfusedID),
	},
	Triggers: []game.Trigger{
		PoisonInfusedTrigger,
	},
}
