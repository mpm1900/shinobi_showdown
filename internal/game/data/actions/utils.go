package actions

import (
	"math/rand/v2"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

func makeAttackConfig(base game.ActionConfig) game.ActionConfig {
	base.Cost = game.Ptr(0)
	base.Cooldown = game.Ptr(0)
	base.CritStage = game.Ptr(0)
	base.CritMod = 1.5
	base.TargetCount = game.Ptr(1)
	base.TargetType = game.TargetPositionID
	return base
}

func makeSpreadAttackConfig(base game.ActionConfig) game.ActionConfig {
	base.Cost = game.Ptr(0)
	base.Cooldown = game.Ptr(0)
	base.CritStage = game.Ptr(0)
	base.CritMod = 1.5
	base.TargetCount = game.Ptr(0)
	base.TargetType = game.TargetPositionID
	return base
}

func makeNoTargetStatusConfig(base game.ActionConfig) game.ActionConfig {
	base.Cooldown = game.Ptr(0)
	base.TargetCount = game.Ptr(0)
	base.TargetType = game.TargetActorID
	return base
}

func makeStatusConfig(base game.ActionConfig) game.ActionConfig {
	base.Cooldown = game.Ptr(0)
	base.TargetCount = game.Ptr(1)
	base.TargetType = game.TargetPositionID
	return base
}

type AttackConfig struct {
	ID              uuid.UUID
	Config          game.ActionConfig
	MapContext      func(game.Game, game.Context) game.Context
	MapConfig       func(game.Game, game.Context, game.ActionConfig) game.ActionConfig
	TargetPredicate func(game.Game, game.Actor, game.Context) bool
	Priority        *int
	BeforeAttack    func(game.Game, game.Context) []game.GameTransaction
	OnSuccess       func(game.Game, game.Context, game.Context) []game.GameTransaction
	OnFailure       func(game.Game, game.Context, game.Context) []game.GameTransaction
	AfterAttack     func(game.Game, game.Context) []game.GameTransaction
}

func makeAttack(config AttackConfig) game.Action {
	action := game.Action{
		ID:              config.ID,
		Config:          config.Config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.Config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Config.Cost),
		MapContext:      config.MapContext,
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				if config.BeforeAttack != nil {
					transactions = append(transactions, config.BeforeAttack(g, context)...)
				}

				action_config, _ := game.GetActiveActionConfig(g, config.Config)
				if config.MapConfig != nil {
					action_config = config.MapConfig(g, context, action_config)
				}

				dmg_config := game.NewDamageConfig(game.RandomDamageFactor())
				if config.OnSuccess != nil {
					dmg_config.OnSuccess = config.OnSuccess
				}
				if config.OnFailure != nil {
					dmg_config.OnFailure = config.OnFailure
				}
				damages := game.DamageCoreMutation(action_config, dmg_config)
				transactions = append(
					transactions,
					game.MakeDamageTransactions(context, damages)...,
				)

				if config.AfterAttack != nil {
					transactions = append(transactions, config.AfterAttack(g, context)...)
				}

				return transactions
			},
		},
	}

	if config.TargetPredicate != nil {
		action.TargetPredicate = config.TargetPredicate
	}

	if config.Priority != nil {
		action.Priority = *config.Priority
	}

	return action
}

func applySummon(context game.Context, def game.ActorDef, actions []game.Action) []game.GameTransaction {
	transactions := []game.GameTransaction{}

	mut := game.GameMutation{
		Delta: func(mp, mg game.Game, mc game.Context) game.Game {
			mg.UpdateActor(*mc.SourceActorID, func(a game.Actor) game.Actor {
				summon := game.MakeActor(
					def,
					a.PlayerID,
					a.Experience,
					nil,
					nil,
					append(actions, game.CancelSummon),
					game.FocusNone,
					map[game.ActorStat]int{},
				)
				a.SetSummonFromActor(&summon, false)
				return a
			})
			mg.UpdatePlayer(*mc.SourcePlayerID, func(p game.Player) game.Player {
				p.UsedSummon = true
				return p
			})
			return mg
		},
	}

	transactions = append(
		transactions,
		game.MakeTransaction(mut, context),
	)

	return transactions
}

func checkPlayerHasModifier(g game.Game, context game.Context, modifierID uuid.UUID) bool {
	for _, tx := range g.GetModifiers() {
		if tx.Context.SourcePlayerID == nil {
			continue
		}

		if *tx.Context.SourcePlayerID == *context.SourcePlayerID && tx.Mutation.ID == modifierID {
			return true
		}

	}

	return false
}

func MakeRepeats(config game.DamageCoreConfig, min int, max int, g game.Game, context game.Context) game.DamageCoreConfig {
	source, ok := g.GetSource(context)
	if !ok {
		return config
	}

	resolved := source.Resolve(g)
	min_offset := resolved.RepeatsMinOffset
	max_offset := resolved.RepeatsMaxOffset
	offset := min + min_offset
	rand_max := max + max_offset - offset
	if rand_max <= 0 {
		config.Repeat = (offset) > 1
		config.RepeatMax = (offset)
		return config
	}

	config.Repeat = true
	config.RepeatMax = rand.IntN(rand_max) + (offset)
	return config
}
