package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Flash = MakeBlindingFlash()

func MakeBlindingFlash() game.Action {
	config := game.ActionConfig{
		Name:        "Flash",
		Nature:      game.Ptr(game.NsYin),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Genjutsu,
		Description: "Stuns the target this turn. Usable on the turn after the user switched in. +3 priority.",
	}
	return game.Action{
		ID:              uuid.MustParse("4cf69985-6785-56a6-b879-e02cb6207960"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityP3,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
				game.SourceHasActiveTurns(1),
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}
				mod := modifiers.Stunned
				mod.Duration = 0

				mutation := mutations.AddModifiers(false, mod)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
