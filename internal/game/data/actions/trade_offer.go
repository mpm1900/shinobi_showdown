package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var TradeOffer = MakeTradeOffer()

func MakeTradeOffer() game.Action {
	config := game.ActionConfig{
		Name:        "Trade Offer",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Taijutsu,
		Description: "Exchanges held items with the target.",
	}

	return game.Action{
		ID:              uuid.MustParse("932069be-0954-486f-b4c7-8ee59cf28ee2"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}

				tx := game.MakeTransaction(mutations.ExchangeItems, context)
				transactions = append(transactions, tx)

				return transactions
			},
		},
	}
}
