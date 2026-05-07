package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"
)

func UseStaminaCost(amount int) game.GameMutation {
	return game.GameMutation{
		Delta: func(p game.Game, g game.Game, context game.Context) game.Game {
			g = mutations.UseStaminaSource(amount).Delta(p, g, context)
			g = mutations.CheckStaminaExhaustion(Sleeping, mutations.Sleep).Delta(p, g, context)
			return g
		},
	}
}
