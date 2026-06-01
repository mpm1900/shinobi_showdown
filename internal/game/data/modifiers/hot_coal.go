package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var hotCoalID = uuid.MustParse("aeccc179-4ee1-420c-ac2e-07e5442d7ec0")
var HotCoalTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: hotCoalID,
	On:         game.OnTurnEnd,
	Check:      game.Match__SourceIsNotStatused,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p, g game.Game, context game.Context) []game.GameTransaction{
			transactions := game.NewTransactionBuilder()
			source, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			transactions.Push(ApplyBurn(game.ActionConfig{}, g, source, game.NewContext()))
			return transactions.Build()
		},
	},
}

var HotCoal = game.Modifier{
	ID:          hotCoalID,
	GroupID:     &hotCoalID,
	Name:        "Hot Coal",
	Description: "On turn end: burns holder.",
	Icon:        "hot_coal",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&hotCoalID),
	},
	Triggers: []game.Trigger{
		HotCoalTrigger,
	},
}
