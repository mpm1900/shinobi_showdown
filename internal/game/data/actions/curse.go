package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Curse = MakeCurse()

func MakeCurse() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:   "Curse",
		Nature: game.Ptr(game.NsJashin),
		Jutsu:  game.Fuinjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("79b92672-ecc5-5b29-b8f2-035879182bda"),
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
					amount := hp / 2

					hp_loss := game.PureDamageWith(amount, false, func(a game.Actor) game.Actor {
						if !a.Alive {
							a.Alive = true
							a.Damage = hp - 1
						}
						return a
					})
					damage := game.PureDamage(amount, false)

					sourceTx := game.MakeTransaction(hp_loss, game.NewContext().WithSource(source.ID).WithTargetIDs([]uuid.UUID{source.ID}))
					targetTx := game.MakeTransaction(damage, game.NewContext().WithSource(source.ID).WithTargetIDs([]uuid.UUID{target.ID}))

					transactions = append(transactions, sourceTx, targetTx)
				}

				return transactions
			},
		},
	}
}
