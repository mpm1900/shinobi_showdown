package mutations

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

func QueueAction(actionID uuid.UUID, context game.Context) game.GameMutation {
	return game.GameMutation{
		Delta: func(p game.Game, g game.Game, context game.Context) game.Game {
			if context.SourceActorID == nil {
				return g
			}

			g.QueuedActions[*context.SourceActorID] = game.MakeTransaction(actionID, context)
			return g
		},
	}
}
