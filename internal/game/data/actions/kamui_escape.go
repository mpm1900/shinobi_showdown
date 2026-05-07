package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var KamuiEscape = MakeKamuiEscape()

func MakeKamuiEscape() game.Action {
	ID := uuid.MustParse("58924949-3d71-4d01-8f11-f20dfe44df14")

	config := game.ActionConfig{
		Name:        "Kamui: Escape",
		Description: "Lowers the target's Attack by 2 stages. User then switches out.",
		Nature:      game.Ptr(game.NsYin),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(30),
		Jutsu:       game.Genjutsu,
	}

	return game.Action{
		ID:              ID,
		Config:          config,
		TargetType:      game.TargetPositionID,
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
					mutation := mutations.AddModifiers(false, modifiers.AttackDown2Target)
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
