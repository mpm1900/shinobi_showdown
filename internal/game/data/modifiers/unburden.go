package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var unburdenID = uuid.MustParse("37dc2127-ac02-4472-af74-4dcd2c25a601")
var UnburdenTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: unburdenID,
	On:         game.OnItemConsume,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			transactions := []game.GameTransaction{}

			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			mut_ctx := game.MakeContextForActor(source)
			mut_ctx.SourceActorID = context.SourceActorID
			mut_ctx.ModifierID = &intimidateID
			mutation := mutations.AddModifiers(false, SpeedUpSource)
			transaction := game.MakeTransaction(mutation, mut_ctx)
			transactions = append(transactions, transaction)

			return transactions
		},
	},
}

var Unburden game.Modifier = game.Modifier{
	ID:          unburdenID,
	GroupID:     &unburdenID,
	Icon:        "unburden",
	Name:        "Unburden",
	Description: "On item consume: raises user's speed.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&unburdenID),
	},
	Triggers: []game.Trigger{
		UnburdenTrigger,
	},
}
