package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var perishDuration = 4

var pdmID = uuid.New()
var PerishDeathModifier game.Modifier = game.Modifier{
	ID:             pdmID,
	GroupID:        &pdmID,
	Name:           "Perish",
	Show:           true,
	Duration:       game.ModifierDurationInf,
	Delay:          perishDuration,
	ActorMutations: []game.ActorMutation{},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: pdmID,
			On:         game.OnTurnEnd,
			Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
				return true
			},
			ActionMutation: game.ActionMutation{
				Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					var transactions []game.Transaction[game.GameMutation]

					source, ok := g.GetSource(context)
					if !ok {
						return transactions
					}

					mut := game.RatioDamage(1.0)
					tx := game.MakeTransaction(mut, game.MakeContextForActor(source))
					transactions = append(transactions, tx)

					return transactions
				},
			},
		},
	},
}

var pimID = uuid.New()
var PerishInfoModifier game.Modifier = game.Modifier{
	ID:       pimID,
	GroupID:  &pimID,
	Name:     "Perish",
	Show:     true,
	Duration: perishDuration,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&pimID),
	},
	Triggers: []game.Trigger{},
}

var PerishSong = MakePerishSong()

func MakePerishSong() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Perish Song",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Fuinjutsu,
		Description: "Kills target after 5 turns.",
	})

	return game.Action{
		ID:              uuid.MustParse("ead0d88a-1b00-417d-a45e-c77a7bcaeb74"),
		Config:          config,
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
					mut := mutations.AddModifiers(true, PerishDeathModifier, PerishInfoModifier)
					tx := game.MakeTransaction(mut, mut_ctx)
					transactions = append(transactions, tx)
				}

				return transactions
			},
		},
	}
}
