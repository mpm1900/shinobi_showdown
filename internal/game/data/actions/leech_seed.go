package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var LeechSeed = MakeLeechSeed()

func MakeLeechSeed() game.Action {
	config := game.ActionConfig{
		Name:        "Plant Roots",
		Description: "Saps health from target every turn.",
		Nature:      game.Ptr(game.NsYang),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(30),
		Jutsu:       game.Senjutsu,
	}

	return game.Action{
		ID:              uuid.MustParse("9ad36f89-03c5-5b52-9f50-66864b06ca03"),
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				source, ok := g.GetSource(context)
				if !ok || source.PositionID == nil {
					return transactions
				}

				targets := g.GetTargets(context)
				for _, target := range targets {
					mut_ctx := game.Context{
						SourcePlayerID:    &source.PlayerID,
						SourceActorID:     &source.ID,
						ParentActorID:     &target.ID,
						TargetPositionIDs: []uuid.UUID{*source.PositionID},
					}
					mutation := mutations.AddModifiers(false, LeechSeedModifier)
					transaction := game.MakeTransaction(mutation, mut_ctx)
					transactions = append(transactions, transaction)
				}

				return transactions
			},
		},
	}
}

var leechSeedModifierID = uuid.MustParse("531ca1a4-40af-5de4-a3c6-99f8468cd368")

var LeechSeedTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: leechSeedModifierID,
	On:         game.OnTurnEnd,
	Check: func(p game.Game, g game.Game, ctx game.Context, t game.Transaction[game.Modifier]) bool {
		return true
	},
	ActionMutation: game.ActionMutation{
		Priority: 0,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			transactions := []game.GameTransaction{}
			parent, ok := g.GetParent(context)
			if !ok {
				return transactions
			}

			resolved_parent := parent.Resolve(g)
			targets := g.GetTargets(context)
			for _, target := range targets {
				ratio := 0.125
				hp_loss := game.Round(float64(resolved_parent.Stats[game.StatHP]) * ratio)
				hp_loss_ctx := context
				hp_loss_ctx.TargetPositionIDs = []uuid.UUID{*resolved_parent.PositionID}
				hp_loss_mut := game.PureDamage(hp_loss, false)
				hp_loss_tx := game.MakeTransaction(hp_loss_mut, hp_loss_ctx)

				heal_mut := game.PureHeal(hp_loss)
				heal_ctx := game.MakeContextForActor(target)
				heal_tx := game.MakeTransaction(heal_mut, heal_ctx)

				transactions = append(transactions, hp_loss_tx, heal_tx)
			}

			return transactions
		},
	},
}

var LeechSeedModifier game.Modifier = game.Modifier{
	ID:       leechSeedModifierID,
	GroupID:  &leechSeedModifierID,
	Name:     "Seeded",
	Icon:     "seeded",
	Show:     true,
	Duration: game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopParent(&leechSeedModifierID),
	},
	Triggers: []game.Trigger{
		LeechSeedTrigger,
	},
}
