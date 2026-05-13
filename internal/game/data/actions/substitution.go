package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Substitution = MakeSubstitution()

func MakeSubstitution() game.Action {
	config := game.ActionConfig{
		Name:        "Body Replacement",
		Nature:      game.Ptr(game.NsYin),
		Cooldown:    game.Ptr(1),
		Jutsu:       game.Ninjutsu,
		Description: "Protects the user from actions this turn. +4 priority, 1 turn cooldown.",
	}
	return game.Action{
		ID:              uuid.MustParse("d3765608-4b30-5c4c-b5a9-f4132f0bbb7c"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityProtect,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.Protected)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}

// proxies
var Kamui = MakeKamui()

func MakeKamui() game.Action {
	action := MakeSubstitution()
	action.ID = uuid.MustParse("04669c33-cb56-480a-9c94-187b2acab8e1")
	action.Config.Name = "Kamui"
	return action
}

var NegateJutsu = MakeNegateJutsu()

func MakeNegateJutsu() game.Action {
	action := MakeSubstitution()
	action.ID = uuid.MustParse("923b1582-fd96-4f97-9957-1cdd9169600f")
	action.Config.Name = "Negate Jutsu"
	return action
}
