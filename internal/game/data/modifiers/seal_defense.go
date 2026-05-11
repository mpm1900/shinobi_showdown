package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var sodefID = uuid.MustParse("94d8c719-b6df-4b73-9bb9-f3843cc2b9f7")
var sealOfDefenseTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: sodefID,
	On:         game.OnDamageReceive,
	Check:      game.ComposeTF(game.Match__SourceActor_SourceActor),
	ActionMutation: game.ActionMutation{
		Priority: game.MutPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			mut_ctx := game.MakeContextForActor(source)
			consume_tx := game.MakeTransaction(mutations.ConsumeItem, mut_ctx)
			transactions = append(transactions, consume_tx)
			return transactions
		},
	},
}

var SealOfDefense = game.Modifier{
	ID:          sodefID,
	GroupID:     &sodefID,
	Icon:        "seal_up",
	Name:        "Seal of Defense",
	Description: "Holder takes half damage. On damaged: break this seal.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&sodefID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, a game.Actor, ctx game.Context) game.Actor {
				a.DamageReduction[game.ChakraAttack] *= 0.5
				a.DamageReduction[game.Attack] *= 0.5
				return a
			},
		),
	},
	Triggers: []game.Trigger{
		sealOfDefenseTrigger,
	},
}
