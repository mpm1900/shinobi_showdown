package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Devistation = MakeDevistation()

func MakeDevistation() game.Action {
	ID := uuid.MustParse("3f1e20b2-7a85-4d12-80fd-c2a0991d2ee8")
	config := game.ActionConfig{
		Name:        "Devistation",
		Description: "Deals more damage proportional to remaining HP.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(150),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(1),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
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
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				resolved := source.Resolve(g)
				ratio := resolved.GetHealthRatio()

				action_config, _ := game.GetActiveActionConfig(g, config)
				action_config.Power = game.Ptr(game.Round(float64(*action_config.Power) * ratio))
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
