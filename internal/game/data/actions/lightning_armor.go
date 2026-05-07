package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var LightningArmor = MakeLightningArmor()

func MakeLightningArmor() game.Action {
	config := game.ActionConfig{
		Name:        "Lightning Armor",
		Nature:      game.Ptr(game.NsLightning),
		Jutsu:       game.Ninjutsu,
		Description: "Raises the user's Attack and Defense stats. Sets electrified terrain.",
	}
	return game.Action{
		ID:              uuid.MustParse("32808365-9a64-4102-9ed7-39fe7b795f7d"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}
				mutation := mutations.AddModifiers(false, modifiers.AttackUpSource, modifiers.DefenseUpSource)
				mod_tx := game.MakeTransaction(mutation, context)
				transactions = append(transactions, mod_tx)

				transactions = append(
					transactions,
					modifiers.ApplyTerrain(g, context, game.GameTerrainElectrified, modifiers.ElectrifiedTerrain())...,
				)

				return transactions
			},
		},
	}
}
