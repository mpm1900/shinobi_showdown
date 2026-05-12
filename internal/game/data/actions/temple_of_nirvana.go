package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var TempleOfNirvana = MakeTempleOfNirvana()

func MakeTempleOfNirvana() game.Action {
	ID := uuid.MustParse("d59535f2-9cb5-4268-854e-4d9a1d6b7c70")

	config := game.ActionConfig{
		Name:        "Temple Of Nirvana",
		Description: "Puts target to sleep",
		Nature:      game.Ptr(game.NsYin),
		TargetCount: game.Ptr(1),
		Accuracy:    game.Ptr(85),
		Cost:        game.Ptr(30),
		Jutsu:       game.Genjutsu,
	}

	return game.Action{
		ID:              ID,
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				targets := g.GetTargets(context)
				for _, target := range targets {
					transactions = append(transactions, modifiers.ApplySleep(config, g, target, context)...)
				}

				return transactions
			},
		},
	}
}
