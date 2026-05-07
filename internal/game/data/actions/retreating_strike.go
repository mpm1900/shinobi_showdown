package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var RetreatingStrike = MakeRetreatingStrike()

func MakeRetreatingStrike() game.Action {
	ID := uuid.MustParse("a6c3bd65-c750-4260-bfc2-bcada542c663")

	config := game.ActionConfig{
		Name:        "Retreating Strike",
		Description: "User then switches out.",
		Nature:      game.Ptr(game.NsTai),
		Cost:        game.Ptr(0),
		TargetCount: game.Ptr(1),
		Jutsu:       game.Taijutsu,
		Power:       game.Ptr(40),
		Accuracy:    game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			switch_mux := game.RemovePositions
			switch_ctx := game.NewContext()
			switch_ctx.TargetActorIDs = append(switch_ctx.TargetActorIDs, *context.SourceActorID)
			switch_tx := game.MakeTransaction(switch_mux, switch_ctx)
			transactions = append(transactions, switch_tx)

			return transactions
		},
	})
}
