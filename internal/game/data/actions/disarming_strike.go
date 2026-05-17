package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var DisarmingStrike = MakeDisarmingStrike()

func MakeDisarmingStrike() game.Action {
	ID := uuid.MustParse("bfeccb7c-4d50-4575-8a95-f3357d6f81ae")

	config := game.ActionConfig{
		Name:        "Disarming Strike",
		Description: "Target loses their held item.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(50),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(50),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				ctx := game.MakeContextForActor(target)
				tx := game.MakeTransaction(mutations.RemoveItem, ctx)
				transactions = append(transactions, tx)
			}

			return transactions
		},
	})
}
