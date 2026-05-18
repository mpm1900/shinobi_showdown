package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Revival = MakeRevival()

func MakeRevival() game.Action {
	config := game.ActionConfig{
		Name:        "Revival",
		Nature:      game.Ptr(game.NsYin),
		Cooldown:    game.Ptr(100),
		Jutsu:       game.Genjutsu,
		Description: "Revives a fallen ally. Usable once per battle.",
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetActorID,
	}

	return game.Action{
		ID:              uuid.MustParse("fa748b4e-2832-476c-b792-f1d525641280"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.TeamFilter, game.InactiveFilter, game.NotAliveFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				for _, target := range g.GetTargets(context) {
					mut_ctx := game.MakeContextForActor(target)
					mut := mutations.Revive
					tx := game.MakeTransaction(mut, mut_ctx)
					transactions = append(transactions, tx)
				}

				return transactions
			},
		},
	}
}
