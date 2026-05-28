package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var OneHundredHealingsStatus = MakeOneHundredHealingsStatus()

func MakeOneHundredHealingsStatus() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "100 Healings: Status",
		Nature:      game.Ptr(game.NsYang),
		Jutsu:       game.Ninjutsu,
		Description: "Heals status from user's party.",
	})

	return game.Action{
		ID:              uuid.MustParse("cea0796c-df52-466f-8474-9dc06ec9db6f"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := game.NewTransactionBuilder()

				party := g.GetActorsFilters(context, game.TeamFilter)
				for _, actor := range party {
					if !actor.Statused {
						continue
					}
					transactions.Push(modifiers.ClearStatus(g, actor))
				}

				return transactions.Build()
			},
		},
	}
}
