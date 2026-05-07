package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var sharkSkinID = uuid.MustParse("42b677d6-cabe-4fe0-af0f-5956299e7e52")

var SkarkSkinTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: sharkSkinID,
	On:         game.OnDamageReceive,
	Check:      game.Match__TargetActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: 0,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			transactions := []game.GameTransaction{}

			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			mut_ctx := game.MakeContextForActor(source)
			mutation := game.RatioDamage(0.125)
			transaction := game.MakeTransaction(mutation, mut_ctx)
			transactions = append(transactions, transaction)

			return transactions
		},
	},
}

var SharkSkin game.Modifier = game.Modifier{
	ID:          sharkSkinID,
	GroupID:     &sharkSkinID,
	Icon:        "shark_skin",
	Name:        "Shark Skin",
	Description: "On physical damage: attacker loses 1/8th HP.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&sharkSkinID),
	},
	Triggers: []game.Trigger{
		SkarkSkinTrigger,
	},
}
