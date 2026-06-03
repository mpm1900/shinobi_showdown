package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Graft = MakeGraft()

func MakeGraft() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Graft",
		Description: "Heals an ally or damages an enemy.",
		Nature:      game.Ptr(game.NsWood),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(70),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("fdbfd320-071a-46a2-b449-e1455d1a3d14"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := game.NewTransactionBuilder()
				action_config, _ := game.GetActiveActionConfig(g, config)
				damage_config := game.NewDamageConfig(game.RandomDamageFactor())

				for _, target := range g.GetTargets(context) {
					isTeam := context.SourcePlayerID != nil && target.PlayerID == *context.SourcePlayerID
					ctx := context
					ctx.TargetPositionIDs = []uuid.UUID{*target.PositionID}
					if isTeam {
						heals := game.RatioHeal(0.5)
						transactions.PushOne(game.MakeTransaction(heals, ctx))
					} else {
						transactions.Push(game.ResolveDamageCore(action_config, damage_config, g, ctx))
					}
				}

				return transactions.Build()
			},
		},
	}
}
