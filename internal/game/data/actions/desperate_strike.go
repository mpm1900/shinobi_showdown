package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var DesperateStrike = MakeDesperateStrike()

func MakeDesperateStrike() game.Action {
	config := game.ActionConfig{
		Name:        "Desperate Strike",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Taijutsu,
		Description: "Target's HP becomes equal to the user's.",
		TargetCount: game.Ptr(1),
	}

	return game.Action{
		ID:              uuid.MustParse("55665151-dfec-40d3-ad45-96ef53d716e9"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mut := game.SetHpSourceToTargets()
				transactions = append(transactions, game.MakeTransaction(mut, context))

				return transactions
			},
		},
	}
}
