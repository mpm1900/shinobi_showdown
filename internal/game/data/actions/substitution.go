package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Substitution = MakeSubstitution()

func MakeSubstitution() game.Action {
	action := makeNoneAction(
		uuid.MustParse("d3765608-4b30-5c4c-b5a9-f4132f0bbb7c"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Substitution",
			Nature:      game.Ptr(game.NsYin),
			Jutsu:       game.Ninjutsu,
			Description: "Protects the user from actions this turn. +4 priority, 1 turn cooldown.",
			Cooldown:    game.Ptr(1),
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.Protected)
			transaction := game.MakeTransaction(mutation, context)
			transactions.Push(transaction)

			return transactions.Build()
		},
	)

	action.Priority = game.ActionPriorityProtect
	return action
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

var IaiBlock = MakeIaiBlock()

func MakeIaiBlock() game.Action {
	action := MakeSubstitution()
	action.ID = uuid.MustParse("4a1ebde5-adb7-457c-a8a3-57fa0d13e70f")
	action.Config.Name = "Iai: Block"
	action.Config.Nature = game.Ptr(game.NsTai)
	action.Config.Jutsu = game.Bukijutsu
	return action
}
