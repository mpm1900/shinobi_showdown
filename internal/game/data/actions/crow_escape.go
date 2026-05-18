package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var MirageCrow = MakeMirageCrow()

func MakeMirageCrow() game.Action {
	config := game.ActionConfig{
		Name:        "Crow Escape",
		Description: "Lowers the target's Chakra Attack by 2 stages. User then switches out.",
		Nature:      game.Ptr(game.NsYin),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(30),
		Jutsu:       game.Genjutsu,
	}

	return game.Action{
		ID:              uuid.MustParse("355753f0-5332-5ee4-b438-899d1a71c184"),
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
					mut_ctx := context
					mut_ctx.ParentActorID = &target.ID
					mut_ctx.TargetActorIDs = []uuid.UUID{target.ID}
					mutation := mutations.AddModifiers(false, modifiers.ChakraAttackDown2Target)
					transaction := game.MakeTransaction(mutation, mut_ctx)
					transactions = append(transactions, transaction)
				}

				switch_mux := game.RemovePositions
				switch_ctx := game.NewContext()
				switch_ctx.TargetActorIDs = append(switch_ctx.TargetActorIDs, *context.SourceActorID)
				switch_tx := game.MakeTransaction(switch_mux, switch_ctx)
				transactions = append(transactions, switch_tx)

				return transactions
			},
		},
	}
}
