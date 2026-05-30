package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ShinigamiCurse = MakeShinigamiCurse()

func MakeShinigamiCurse() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Shinigami Curse",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Fuinjutsu,
		Description: "Lowers the target's Attack and Chakra Attack by 2 stages. User dies.",
		Cost:        game.Ptr(0),
	})

	return game.Action{
		ID:              uuid.MustParse("f3be1a59-8c31-43ca-bf96-4e5d109e81e8"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}
				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				targets := g.GetTargets(context)
				for _, target := range targets {
					mut_ctx := game.MakeContextForActor(target)
					mutation := mutations.AddModifiers(false, modifiers.AttackDown2Target, modifiers.ChakraAttackDown2Target)
					transaction := game.MakeTransaction(mutation, mut_ctx)
					transactions = append(transactions, transaction)
				}

				self_dmg := mutations.KillSource()
				self_dmg_ctx := game.MakeContextForActor(source)
				transactions = append(transactions, game.MakeTransaction(self_dmg, self_dmg_ctx))

				return transactions
			},
		},
	}
}
