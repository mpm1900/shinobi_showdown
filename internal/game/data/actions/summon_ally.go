package actions

import (
	"fmt"
	"maps"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var SummonAlly = MakeSummonAlly()

func MakeSummonAlly() game.Action {
	return game.Action{
		ID: uuid.MustParse("c06b7803-7b52-4e20-a359-e92695920896"),
		Config: game.ActionConfig{
			Name:        "Summon Ally",
			Nature:      game.Ptr(game.NsYin),
			Description: "Summons ally to battle, they gain user's stat up/downs. Switches user out.",
			TargetType:  game.TargetActorID,
		},
		TargetPredicate: game.ComposeAF(game.TeamFilter, game.InactiveFilter, game.AliveFilter),
		ContextValidate: game.TargetLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.TrueGameFilter,
			Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{
					game.MakeTransaction(game.SwitchPositions, context),
				}

				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				resolved := source.Resolve(g)
				stages := maps.Clone(resolved.Stages)
				for _, target := range g.GetTargets(context) {
					mod := modifiers.NewStageDeltaFromMap(
						stages,
						nil,
						game.ComposeAF(game.ActiveFilter, game.TargetFilter),
						game.MutPriorityDefault,
					)
					mod.Name = fmt.Sprintf("%s's stat changes", source.Name)
					mut := mutations.AddModifiers(false, mod)
					tx := game.MakeTransaction(mut, game.MakeContextForActor(target))
					transactions = append(transactions, tx)
				}

				return transactions
			},
		},
	}
}
