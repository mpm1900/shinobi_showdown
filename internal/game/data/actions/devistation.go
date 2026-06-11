package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Devistation = MakeDevistation()

func MakeDevistation() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Devistation",
		Description: "Deals more damage proportional to remaining HP.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(150),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("3f1e20b2-7a85-4d12-80fd-c2a0991d2ee8"),
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

				source, ok := g.GetSource(context)
				if !ok {
					return transactions.Build()
				}

				resolved := source.Resolve(g)
				ratio := resolved.GetHealthRatio()

				action_config, _ := game.GetActiveActionConfig(g)
				action_config.Power = game.Ptr(game.Round(float64(*action_config.Power) * ratio))
				damage_config := game.NewDamageConfig(game.RandomDamageFactor())
				transactions.Concat(game.ResolveDamageCore(action_config, damage_config, g, context))

				return transactions.Build()
			},
		},
	}
}
