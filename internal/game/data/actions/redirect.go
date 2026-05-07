package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Redirect = MakeRedirect()

func MakeRedirect() game.Action {
	config := game.ActionConfig{
		Name:        "Redirect",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Genjutsu,
		Description: "Changes the target of single-target enemy actions.",
	}
	return game.Action{
		ID:              uuid.MustParse("3d0b6e04-f5f0-50db-9eb6-4aede4c11701"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityP2,
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

				mutation := mutations.RedirectSingleTargetEnemyActions(source)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
