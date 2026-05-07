package modifiers

import (
	"fmt"
	"math/rand/v2"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var copyAbilityID = uuid.MustParse("b2f94cae-3d34-417d-9b36-160d3d37bd5d")

var CopyAbility = game.Modifier{
	ID:          copyAbilityID,
	GroupID:     &copyAbilityID,
	Name:        "Copy Ability",
	Icon:        "copy_ability",
	Description: "On enter: copies random allie's ability.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&copyAbilityID),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: copyAbilityID,
			On:         game.OnActorEnter,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityP3,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					source, ok := g.GetSource(context)
					if !ok {
						return transactions
					}

					targets := g.GetActorsFilters(context, game.ComposeAF(
						game.ActiveFilter,
						game.AliveFilter,
						game.OtherTeamFilter,
					))

					if len(targets) == 0 {
						return transactions
					}

					index := rand.IntN(len(targets))
					target := targets[index]

					fmt.Println(target.Name)

					mut_ctx := game.MakeContextForActor(source)
					mut := game.GameMutation{
						Delta: func(p, g game.Game, context game.Context) game.Game {
							g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
								a.AuxAbility = target.Ability
								return a
							})
							return g
						},
					}
					transaction := game.MakeTransaction(mut, mut_ctx)
					transactions = append(transactions, transaction)

					return transactions
				},
			},
		},
	},
}
