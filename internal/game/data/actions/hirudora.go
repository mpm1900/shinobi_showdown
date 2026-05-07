package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Hirudora = MakeHirudora()

func MakeHirudora() game.Action {
	ID := uuid.MustParse("36c88bee-1eb0-4c55-bdc6-b704221ea846")
	config := game.ActionConfig{
		Name:        "7th Gate: Hirudora",
		Description: "Lowers user's Attack by 2 stages.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(130),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(50),
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			mutation := mutations.AddModifiers(false, modifiers.AttackDown2Source)
			transaction := game.MakeTransaction(mutation, game.MakeContextForActor(source))
			transactions = append(transactions, transaction)

			return transactions
		},
	})
}
