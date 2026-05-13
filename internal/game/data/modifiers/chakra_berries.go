package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var FireChakraBerry = MakeChakraBerry(
	uuid.MustParse("113e38c7-9d50-4870-9515-09803ee35b37"),
	game.NsFire,
	"Fire Chakra Berry",
	"Halves the attack of super-effective fire damage. Consumes after use.",
	"fire_berry",
)

var WindChakraBerry = MakeChakraBerry(
	uuid.MustParse("e94dc891-5682-49fe-a4e7-4aaf487ffdd6"),
	game.NsWind,
	"Wind Chakra Berry",
	"Halves the attack of super-effective wind damage. Consumes after use.",
	"wind_berry",
)

var LightningChakraBerry = MakeChakraBerry(
	uuid.MustParse("9befaf46-f5d7-49d7-b028-5d34feadca69"),
	game.NsWind,
	"Lightning Chakra Berry",
	"Halves the attack of super-effective lightning damage. Consumes after use.",
	"lightning_berry",
)

var EarthChakraBerry = MakeChakraBerry(
	uuid.MustParse("43127a1e-4be9-490e-b3af-ed213e64d478"),
	game.NsEarth,
	"Earth Chakra Berry",
	"Halves the attack of super-effective earth damage. Consumes after use.",
	"earth_berry",
)

var WaterChakraBerry = MakeChakraBerry(
	uuid.MustParse("8c24e7a3-02bd-45bd-9e01-a37597bf9099"),
	game.NsWater,
	"Water Chakra Berry",
	"Halves the attack of super-effective water damage. Consumes after use.",
	"water_berry",
)

func MakeChakraBerry(id uuid.UUID, nature game.NatureSet, name string, description string, icon string) game.Modifier {
	return game.Modifier{
		ID:          id,
		Name:        name,
		Description: description,
		Icon:        icon,
		Show:        true,
		GroupID:     &id,
		Duration:    0,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				&id,
				game.MutPriorityPostStagedStats,
				game.ComposeAF(game.ActiveFilter, game.SourceFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					atx := g.ActiveTransaction
					if atx == nil {
						return a
					}

					action := atx.Mutation
					if action.Config.Nature == nil || *action.Config.Nature != nature {
						return a
					}
					atx_source, ok := g.GetSource(atx.Context)
					if !ok {
						return a
					}

					targets := g.GetTargets(atx.Context)
					for _, target := range targets {
						if target.ID == a.ID {
							natures := game.NATURES[*action.Config.Nature]
							nature_result := game.ResolveNatures(natures, atx_source.NatureDamage, a.NatureResistance, a.Natures)
							if nature_result.Average < game.NATURE_WEAKNESS_FULL {
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
				ModifierID: id,
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
						return *cfg.Nature == nature
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
}
