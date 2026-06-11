package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Flash = MakeBlindingFlash()

func MakeBlindingFlash() game.Action {
	action := makeSingleStatus(
		uuid.MustParse("4cf69985-6785-56a6-b879-e02cb6207960"),
		makeStatusConfig(game.ActionConfig{
			Name:        "Flash",
			Nature:      game.Ptr(game.NsYin),
			Jutsu:       game.Genjutsu,
			Description: "Stuns the target this turn. Usable on the turn after the user switched in. +3 priority.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mod := modifiers.Stunned
			mod.Duration = 0

			mutation := mutations.AddModifiers(false, mod)
			transaction := game.MakeTransaction(mutation, context)
			transactions.PushOne(transaction)

			return transactions.Build()
		},
	)

	action.ActionMutation.Priority = game.ActionPriorityP3
	action.ActionMutation.Filter = game.ComposeGF(
		game.SourceIsAlive,
		game.SourceIsActionOffCooldown,
		game.SourceHasActiveTurns(1),
	)

	return action
}
