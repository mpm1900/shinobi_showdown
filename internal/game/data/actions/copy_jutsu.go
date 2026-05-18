package actions

import (
	"fmt"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var CopyJutsu = MakeCopyJutsu()

func MakeCopyJutsu() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Copy Jutsu",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Ninjutsu,
		Description: "Copies target's last used action.",
	})

	return game.Action{
		ID:              uuid.MustParse("7bc916c9-f762-472c-ac29-d5e7996b64e3"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				targets := g.GetTargets(context)
				for _, target := range targets {
					if target.LastUsedActionTX == nil {
						continue
					}

					resolved := target.Resolve(g)
					action, ok := resolved.GetActionByID(resolved.LastUsedActionTX.Mutation.ID)
					if !ok {
						fmt.Println("ERROR: COPY JUTSU ACTION NOT FOUND")
						continue
					}

					action_ctx := game.MakeContextForActor(source)
					action_ctx.ActionID = &action.ID
					action_ctx.TargetActorIDs = []uuid.UUID{target.ID}
					action_ctx.TargetPositionIDs = []uuid.UUID{}

					// Copied attack deltas read active action config from game state.
					// Run against a local copy with no active transaction so the
					// copied action's own config is used (power/stat/accuracy).
					gCopy := g
					gCopy.ActiveTransaction = nil
					transactions = append(transactions, action.Delta(p, gCopy, action_ctx)...)
					// TODO, maybe add the action to use until switch out? (modifier add action)
					//transaction := game.MakeTransaction(mut, context)
					//transactions = append(transactions, transaction)

					fmt.Println(action.Config.Name, len(transactions))
				}

				return transactions
			},
		},
	}
}
