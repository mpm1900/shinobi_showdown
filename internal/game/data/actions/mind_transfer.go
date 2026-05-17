package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var MindTransfer = MakeMindTransfer()

func MakeMindTransfer() game.Action {
	config := game.ActionConfig{
		Name:        "Mind Transfer",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Ninjutsu,
		Description: "Switches abilities with the target temporarily.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("f7a33bde-db98-45e1-8d4d-028afe124aeb"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.TeamFilter, game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.TargetLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}
				targets := g.GetTargets(context)
				for _, target := range targets {
					mut := game.GameMutation{
						Delta: func(p, g game.Game, context game.Context) game.Game {
							g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
								a.AuxAbility = target.Ability
								return a
							})
							g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
								a.AuxAbility = source.Ability
								return a
							})
							return g
						},
					}
					transaction := game.MakeTransaction(mut, context)
					transactions = append(transactions, transaction)
				}

				return transactions
			},
		},
	}
}
