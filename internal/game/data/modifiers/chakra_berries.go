package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"
	"slices"

	"github.com/google/uuid"
)

var fireChakraBerryID = uuid.MustParse("113e38c7-9d50-4870-9515-09803ee35b37")
var FireChakraBerry = game.Modifier{
	ID:          fireChakraBerryID,
	Name:        "Fire Chakra Berry",
	Description: "Halves the attack of super-effective fire damage. Consumes after use.",
	Icon:        "fire_berry",
	Show:        true,
	GroupID:     &fireChakraBerryID,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&fireChakraBerryID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				atx := g.ActiveTransaction
				if atx == nil {
					return a
				}

				action := atx.Mutation
				atx_source, ok := g.GetSource(atx.Context)
				if !ok || action.Config.Nature == nil {
					return a
				}

				targets := g.GetTargets(atx.Context)
				for _, target := range targets {
					if target.ID == a.ID {
						natures := game.NATURES[*action.Config.Nature]
						nature_result := game.ResolveNatures(natures, atx_source.NatureDamage, a.NatureResistance, a.Natures)
						if !slices.Contains(natures, game.NatureFire) || nature_result.Average < game.NATURE_WEAKNESS_FULL {
							return a
						}

						a.DamageReduction[game.ChakraAttack] *= 0.5
						a.DamageReduction[game.Attack] *= 0.5
					}
				}

				return a
			},
		),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: fireChakraBerryID,
			On:         game.OnWeakness,
			Check: game.ComposeTF(
				game.Match__TargetActor_SourceActor,
				func(parent game.Game, g game.Game, context game.Context, modifier_tx game.Transaction[game.Modifier]) bool {
					if g.ActiveTransaction == nil {
						return false
					}
					cfg := g.ActiveTransaction.Mutation.Config
					if cfg.Nature == nil {
						return false
					}
					return *cfg.Nature == game.NsFire
				},
			),
			ActionMutation: game.ActionMutation{
				Priority: game.MutPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
					var transactions []game.GameTransaction
					for _, target := range g.GetTargets(context) {
						mut_ctx := game.MakeContextForActor(target)
						transactions = append(transactions, game.MakeTransaction(mutations.ConsumeItem, mut_ctx))
					}
					return transactions
				},
			},
		},
	},
}
