package actions

import (
	"fmt"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"
	"slices"

	"github.com/google/uuid"
)

var focusingID = uuid.MustParse("6bc5e8a3-ddae-4fe2-8ddd-bc0af50b99a3")
var focusPunchID = uuid.MustParse("89f545e9-af99-444a-ac55-eb65b9483632")
var focusPunchAttackID = uuid.MustParse("a52d55c4-11b3-43ce-a680-09c32f453092")

var Focusing = game.Modifier{
	ID:       focusingID,
	GroupID:  &focusingID,
	Icon:     "focusing",
	Name:     "Focusing",
	Show:     true,
	Duration: 0,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&focusingID),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: focusingID,
			On:         game.OnDamageReceive,
			Check:      game.Match__TargetActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					transactions := []game.GameTransaction{}

					mod_mut := mutations.RemoveModifierWhere(func(tx game.Transaction[game.Modifier]) bool {
						if tx.Context.SourceActorID == nil || tx.Mutation.GroupID == nil {
							return false
						}

						if *tx.Mutation.GroupID != focusingID {
							return false
						}

						if slices.Contains(context.TargetActorIDs, *tx.Context.SourceActorID) {
							return true
						}

						return false
					})
					mod_tx := game.MakeTransaction(mod_mut, context)
					transactions = append(transactions, mod_tx)

					return transactions
				},
			},
		},
	},
}

var focusPunchConfig = game.ActionConfig{
	Name:        "Focus Punch",
	Nature:      game.Ptr(game.NsTai),
	Jutsu:       game.Taijutsu,
	Description: "User tightens their focus. If not damaged this turn, they unleash a powerful attack.",
	TargetCount: game.Ptr(1),
	Accuracy:    game.Ptr(100),
	Power:       game.Ptr(150),
	Stat:        game.Ptr(game.StatAttack),
	Cooldown:    game.Ptr(0),
	Cost:        game.Ptr(0),
	CritChance:  game.Ptr(getCriticalStage(0)),
	CritMod:     1.5,
}

var FocusPunch = MakeFocusPunch()

func MakeFocusPunch() game.Action {
	return game.Action{
		ID:              focusPunchID,
		Config:          focusPunchConfig,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*focusPunchConfig.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityP3,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}

				mod_mut := mutations.AddModifiers(false, Focusing)
				mod_tx := game.MakeTransaction(mod_mut, context)
				transactions = append(transactions, mod_tx)

				action_mut := mutations.PushExtraAction(MakeFocusPunchAttack(), context)
				action_tx := game.MakeTransaction(action_mut, context)
				transactions = append(transactions, action_tx)

				return transactions
			},
		},
	}
}

func MakeFocusPunchAttack() game.Action {
	action := makeAttack(AttackConfig{
		ID:       focusPunchAttackID,
		Config:   focusPunchConfig,
		Priority: game.Ptr(game.ActionPrioritySlow3),
	})

	delta := action.Delta
	action.Delta = func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
		source, ok := g.GetSource(context)
		if !ok {
			return []game.GameTransaction{}
		}

		resolved := source.Resolve(g)
		_, ok = resolved.AppliedModifiers[focusingID]
		if !ok {
			log := game.AddLogs(
				game.MakeGameLog(fmt.Sprintf("%s failed.", focusPunchConfig.Name), context, 1),
			)
			return []game.GameTransaction{
				game.MakeTransaction(log, context),
			}
		}

		return delta(p, g, context)
	}

	return action
}
