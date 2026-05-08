package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var sleepyModifierID = uuid.MustParse("c45e694d-9211-479d-95f5-c63d6511d1d3")
var SleepyModifier = game.Modifier{
	ID:          sleepyModifierID,
	GroupID:     &sleepyModifierID,
	Icon:        "sleepy",
	Name:        "Sleepy",
	Description: "On turn end: fall asleep",
	Show:        true,
	Duration:    2,
	Delay:       1,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&sleepyModifierID),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: sleepyModifierID,
			On:         game.OnTurnEnd,
			Check: func(p, g game.Game, ctx game.Context, t game.Transaction[game.Modifier]) bool {
				return true
			},
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					parent, ok := g.GetParent(context)
					if !ok {
						return transactions
					}

					transactions = append(transactions, modifiers.RemoveModifierSource(sleepyModifierID, parent)...)
					transactions = append(transactions, modifiers.ApplySleep(game.ActionConfig{}, g, parent)...)

					return transactions
				},
			},
		},
	},
}

var Yawn = MakeYawn()

func MakeYawn() game.Action {
	config := game.ActionConfig{
		Name:        "Yawn",
		Nature:      game.Ptr(game.NsYang),
		Jutsu:       game.Senjutsu,
		Description: "Applies sleepy to the target. (Target will fall asleep at the end of the next turn.)",
		TargetCount: game.Ptr(1),
	}

	return game.Action{
		ID:              uuid.MustParse("2ac1ffa2-d197-48fc-a21e-2cca9afe0e19"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}

				for _, target := range g.GetTargets(context) {
					mut_ctx := game.MakeContextForActor(target)
					mut := mutations.AddModifiers(true, SleepyModifier)
					tx := game.MakeTransaction(mut, mut_ctx)
					transactions = append(transactions, tx)
				}

				return transactions
			},
		},
	}
}
