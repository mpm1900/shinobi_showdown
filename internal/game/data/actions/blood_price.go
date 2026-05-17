package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var BloodPrice = MakeBloodPrice()

func MakeBloodPrice() game.Action {
	config := game.ActionConfig{
		Name:        "Blood Price",
		Description: "Deals damage equal to 1.5x the amount damage last dealt to the user by the target.",
		Nature:      game.Ptr(game.NsJashin),
		Jutsu:       game.Fuinjutsu,
		TargetType:  game.TargetPositionID,
	}
	return game.Action{
		ID:              uuid.MustParse("b21ee132-d0b9-4a2d-b96d-d8fdfd4aa7a7"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.TargetLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPrioritySlow3,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				s, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				source := s.Resolve(g)
				targets := g.GetTargets(context)

				for _, target := range targets {
					damage, ok := source.LastReceivedDamage[target.ID]
					if !ok {
						log := game.NewLog("Blood Price failed.")
						transactions = append(transactions, game.MakeTransaction(game.AddLogs(log), game.NewContext()))
						continue
					}
					mut := game.PureDamage(game.Round(float64(damage)*1.5), true)
					mut_ctx := game.NewContext().WithSource(source.ID).WithTargetIDs([]uuid.UUID{target.ID})
					tx := game.MakeTransaction(mut, mut_ctx)
					transactions = append(transactions, tx)
				}

				return transactions
			},
		},
	}
}
