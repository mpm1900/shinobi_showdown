package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var SharinganGlare = MakeSharinganGlare()

func MakeSharinganGlare() game.Action {
	config := makeSpreadAttackConfig(game.ActionConfig{
		Name:        "Sharingan: Glare",
		Description: "Hits all enemy shinobi. Lowers targets' speed.",
		Nature:      game.Ptr(game.NsYin),
		Accuracy:    game.Ptr(95),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(100),
		Jutsu:       game.Genjutsu,
	})

	return makeAttack(AttackConfig{
		ID:              uuid.MustParse("ac1ce2ba-7bec-4b26-b72e-9a657d498697"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		OnSuccess: func(g game.Game, context, tcontext game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			tcontext.SourceActorID = context.SourceActorID
			mod := modifiers.SpeedDownTarget
			mutation := mutations.AddModifiers(false, mod)
			transaction := game.MakeTransaction(mutation, tcontext)
			transactions = append(transactions, transaction)

			return transactions
		},
	})
}
