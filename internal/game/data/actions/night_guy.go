package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var NightGuy = MakeNightGuy()

func MakeNightGuy() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "8th Gate: Night Guy",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Taijutsu,
		Description: "Deals damage equal the amount of HP user has. User dies.",
	})

	return game.Action{
		ID:              uuid.MustParse("d4058ae1-5b26-4477-92a5-346e953172e4"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				s, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				source := s.Resolve(g)
				targets := g.GetTargets(context)

				for _, t := range targets {
					target := t.Resolve(g)
					log := game.AddLogs(
						game.MakeGameLog("$source$ was protected.", context.WithSource(target.ID), 1),
					)
					if tx, protected := target.IsProtected(log); protected {
						transactions = append(transactions, tx)
						continue
					}

					hp := source.Stats[game.StatHP]
					amount := hp - source.Damage
					damage := game.PureDamage(amount, true)
					targetTx := game.MakeTransaction(damage, context.WithTargetIDs([]uuid.UUID{target.ID}))
					transactions = append(transactions, targetTx)
				}

				hp_loss := game.RatioDamage(1.0)
				sourceTx := game.MakeTransaction(hp_loss, context.WithTargetIDs([]uuid.UUID{source.ID}))
				transactions = append(transactions, sourceTx)

				return transactions
			},
		},
	}
}
