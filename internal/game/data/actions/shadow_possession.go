package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ShadowPossession = MakeShadowPossession()

func MakeShadowPossession() game.Action {
	ID := uuid.MustParse("6005ead3-37ca-4793-a166-c41eb7c2f3bd")

	config := game.ActionConfig{
		Name:        "Shadow Possession",
		Description: "Increases an ally's speed or paralyzes an enemy.",
		Nature:      game.Ptr(game.NsYin),
		Accuracy:    game.Ptr(100),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(30),
		Cooldown:    game.Ptr(1),
		Jutsu:       game.Ninjutsu,
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
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				for _, target := range g.GetTargets(context) {
					isTeam := context.SourcePlayerID != nil && target.PlayerID == *context.SourcePlayerID
					ctx := context
					ctx.TargetPositionIDs = []uuid.UUID{*target.PositionID}
					if isTeam {
						mut := mutations.AddModifiers(false, modifiers.SpeedUpTarget)
						transactions = append(
							transactions,
							game.MakeTransaction(mut, ctx),
						)
					} else {
						transactions = append(
							transactions,
							modifiers.ApplyParalysis(config, g, target, context)...,
						)
					}
				}

				return transactions
			},
		},
	}
}
