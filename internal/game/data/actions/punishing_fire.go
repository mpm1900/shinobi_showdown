package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PunishingFire = MakePunishingFire()

func MakePunishingFire() game.Action {
	ID := uuid.MustParse("91b571af-0f7b-42ef-8275-fec11e52c372")

	config := game.ActionConfig{
		Name:        "Punishing Fire",
		Description: "Double power if target is statused.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(60),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
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
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				for _, target := range g.GetTargets(context) {
					action_config, _ := game.GetActiveActionConfig(g, config)
					if target.Statused && action_config.Power != nil {
						action_config.Power = game.Ptr(*action_config.Power * 2)
					}
					crit_result := game.MakeCriticalCheck(action_config)
					dmg_config := game.NewDamageConfig(crit_result.Ratio, game.RandomDamageFactor())
					damages := game.NewDamage(action_config, dmg_config)
					ctx := context
					ctx.TargetPositionIDs = []uuid.UUID{*target.PositionID}
					transactions = append(
						transactions,
						game.MakeDamageTransactions(ctx, damages)...,
					)
				}

				return transactions
			},
		},
	}
}
