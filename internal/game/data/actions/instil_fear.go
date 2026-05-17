package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var InstilFear = MakeInstilFear()

func MakeInstilFear() game.Action {
	ID := uuid.MustParse("0b2947a8-5caf-4588-b4d5-192443bd51e3")

	config := game.ActionConfig{
		Name:        "Instil Fear",
		Description: "Paralyzes the target.",
		Nature:      game.Ptr(game.NsYin),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Accuracy:    game.Ptr(100),
		Cost:        game.Ptr(30),
		Jutsu:       game.Genjutsu,
	}

	return game.Action{
		ID:              ID,
		Config:          config,
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
					transactions = append(
						transactions,
						modifiers.ApplyParalysis(config, g, target, context)...,
					)
				}

				return transactions
			},
		},
	}
}
