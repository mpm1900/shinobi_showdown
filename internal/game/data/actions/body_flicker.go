package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var BodyFlicker = MakeBodyFlicker()

func MakeBodyFlicker() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Body Flicker",
		Description: "User then switches out.",
		Nature:      game.Ptr(game.NsWind),
		Cost:        game.Ptr(0),
		Jutsu:       game.Taijutsu,
		Power:       game.Ptr(70),
		Accuracy:    game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("f052f07c-bb06-4f44-8b26-ec2f17401446"),
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
