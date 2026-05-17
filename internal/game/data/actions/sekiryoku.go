package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Sekiryoku = MakeSekiryoku()

func MakeSekiryoku() game.Action {
	ID := uuid.MustParse("ebb162b4-ead0-5601-afea-4948f147604c")

	config := game.ActionConfig{
		Name:        "Sekiryoku",
		Description: "Forces target to switch out.",
		Nature:      game.Ptr(game.NsYinYang),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
	}

	return game.Action{
		ID:              ID,
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				targets := g.GetTargets(context)
				for _, target := range targets {
					switch_mux := game.RemovePositions
					switch_ctx := game.NewContext()
					switch_ctx.TargetActorIDs = append(switch_ctx.TargetActorIDs, target.ID)
					switch_tx := game.MakeTransaction(switch_mux, switch_ctx)
					transactions = append(transactions, switch_tx)
				}

				return transactions
			},
		},
	}
}
