package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var PuppetAssault = MakePuppetAssault()

func MakePuppetAssault() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Puppet Assault",
		Description: "Lowers user's Defense and Chakra Defense.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Bukijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("9c79987a-6cc3-44eb-a3fa-f1691c989490"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			mutation := mutations.AddModifiers(false, modifiers.DefenseDownSource, modifiers.ChakraDefenseDownSource)
			transaction := game.MakeTransaction(mutation, game.MakeContextForActor(source))
			transactions = append(transactions, transaction)

			return transactions
		},
	})
}
