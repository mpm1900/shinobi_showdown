package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var LightningLigerBomb = MakeLightningLigerBomb()

func MakeLightningLigerBomb() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Lightning Liger Bomb",
		Description: "Deals double damage in Electrified Terrain.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(65),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("4b45490f-b1b9-4c76-b15f-0e8d1e5019cd"),
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

				ratio := 1.0
				state, _ := g.GetState(context)
				if state.Terrain == game.GameTerrainElectrified {
					ratio = 2.0
				}

				action_config, _ := game.GetActiveActionConfig(g, config)
				power := game.Round(float64(*action_config.Power) * ratio)
				action_config.Power = game.Ptr(power)
				crit_result := game.MakeCriticalCheck(action_config)
				dmg_config := game.NewDamageConfig(crit_result.Ratio, game.RandomDamageFactor())
				damages := game.NewDamage(action_config, dmg_config)
				transactions = append(
					transactions,
					game.MakeDamageTransactions(context, damages)...,
				)

				return transactions
			},
		},
	}
}
