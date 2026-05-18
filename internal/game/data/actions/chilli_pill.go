package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ChilliPill = MakeChilliPill()

func MakeChilliPill() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Chilli Pill",
		Nature:      game.Ptr(game.NsYang),
		Jutsu:       game.Fuinjutsu,
		Description: "User pays half of their HP to raise Attack to +6.",
	})

	return game.Action{
		ID:              uuid.MustParse("49092b2d-84b3-47cf-a4a6-ba5fc7d5ff52"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.ComposeGF(game.SourceIsAlive, game.SourceHasHpRatio(0.5)),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}
				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				ctx := game.MakeContextForActor(source)
				cost_mut := game.RatioDamage(0.5)
				cost_tx := game.MakeTransaction(cost_mut, ctx)
				mod_mut := mutations.AddModifiers(false, modifiers.MaxAttackSource)
				mod_tx := game.MakeTransaction(mod_mut, ctx)
				transactions = append(transactions, cost_tx, mod_tx)

				return transactions
			},
		},
	}
}
