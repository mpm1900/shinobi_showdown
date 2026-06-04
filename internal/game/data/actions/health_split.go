package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var HealthSplit = MakeHealthSplit()

func MakeHealthSplit() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Health Split",
		Nature:      game.Ptr(game.NsYang),
		Jutsu:       game.Ninjutsu,
		Description: "Averages the damage between user and target.",
	})

	return game.Action{
		ID:              uuid.MustParse("4566bbac-23e4-4464-a71f-05367d43acf2"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := game.NewTransactionBuilder()

				s, ok := g.GetSource(context)
				if !ok {
					return transactions.Build()
				}

				total_damage := s.Damage
				targets := g.GetTargets(context)
				for _, t := range targets {
					total_damage += t.Damage
				}

				avg_damage := game.Round(float64(total_damage) / float64(len(targets)+1))

				s_ctx := game.MakeContextForActor(s)
				transactions.PushOne(game.MakeTransaction(game.SetDamage(avg_damage), s_ctx))
				for _, t := range targets {
					t_ctx := game.MakeContextForActor(t)
					transactions.PushOne(game.MakeTransaction(game.SetDamage(avg_damage), t_ctx))
				}

				return transactions.Build()
			},
		},
	}
}
