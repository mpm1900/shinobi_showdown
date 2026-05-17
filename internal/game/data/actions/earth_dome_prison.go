package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var EarthDomePrison = MakeEarthDomePrison()

func MakeEarthDomePrison() game.Action {
	ID := uuid.MustParse("c0066f1e-7b7e-40ca-a06b-ade3ad06d19b")

	config := game.ActionConfig{
		Name:        "Earth Dome Prison",
		Description: "Lowers target's Speed.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(95),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(50),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
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
				ctx.SourceActorID = context.SourceActorID
				mod := modifiers.SpeedDownTarget
				mutation := mutations.AddModifiers(false, mod)
				transaction := game.MakeTransaction(mutation, ctx)
				transactions = append(transactions, transaction)
			}

			return transactions
		},
	})
}
