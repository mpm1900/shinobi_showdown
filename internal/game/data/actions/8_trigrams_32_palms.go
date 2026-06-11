package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var ThirtyTwoPalms = MakeThirtyTwoPalms()

func MakeThirtyTwoPalms() game.Action {
	return makeSingleAction(
		uuid.MustParse("7e72e33f-3bcf-4144-ad6f-6d1127c9ab92"),
		makeAttackConfig(game.ActionConfig{
			Name:        "8 Trigrams: 32 Palms",
			Description: "Target loses 50% of the remaining HP. Never misses.",
			Nature:      game.Ptr(game.NsTai),
			Jutsu:       game.Taijutsu,
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			s, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			source := s.Resolve(g)
			targets := g.GetTargets(context)

			for _, t := range targets {
				target := t.Resolve(g)
				hp := target.Stats[game.StatHP] - target.Damage
				amount := hp / 2

				mut := game.PureDamage(amount, true)
				ctx := game.MakeContextForActor(t)
				ctx.SourceActorID = &source.ID
				tx := game.MakeTransaction(mut, ctx)
				transactions.Push(tx)
			}

			return transactions.Build()
		},
	)
}
