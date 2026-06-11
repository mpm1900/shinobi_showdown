package actions

import (
	"math/rand/v2"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

func makeAttackConfig(base game.ActionConfig) game.ActionConfig {
	if base.Cost == nil {
		base.Cost = game.Ptr(0)
	}
	if base.Cooldown == nil {
		base.Cooldown = game.Ptr(0)
	}
	if base.CritStage == nil {
		base.CritStage = game.Ptr(0)
	}
	base.CritMod = 1.5
	if base.TargetCount == nil {
		base.TargetCount = game.Ptr(1)
	}
	base.TargetType = game.TargetPositionID
	return base
}

func makeSpreadAttackConfig(base game.ActionConfig) game.ActionConfig {
	if base.Cost == nil {
		base.Cost = game.Ptr(0)
	}
	if base.Cooldown == nil {
		base.Cooldown = game.Ptr(0)
	}
	if base.CritStage == nil {
		base.CritStage = game.Ptr(0)
	}
	base.CritMod = 1.5
	if base.TargetCount == nil {
		base.TargetCount = game.Ptr(0)
	}
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
	if base.Cooldown == nil {
		base.Cooldown = game.Ptr(0)
	}
	if base.TargetCount == nil {
		base.TargetCount = game.Ptr(1)
	}
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
	BeforeAttack    func(game.Game, game.Context, game.ActionConfig) []game.GameTransaction
	OnSuccess       func(game.Game, game.Context, game.Context, game.ActionConfig) []game.GameTransaction
	OnFailure       func(game.Game, game.Context, game.Context, game.ActionConfig) []game.GameTransaction
	AfterAttack     func(game.Game, game.Context, game.ActionConfig) []game.GameTransaction
}

func (ac AttackConfig) actionConfig(g game.Game, context game.Context) game.ActionConfig {
	action_config, _ := game.GetActiveActionConfig(g, ac.Config)
	if ac.MapConfig != nil {
		action_config = ac.MapConfig(g, context, action_config)
	}

	return action_config
}

func (ac AttackConfig) damageConfig() game.DamageConfig {
	dmg_config := game.NewDamageConfig(game.RandomDamageFactor())
	if ac.OnSuccess != nil {
		dmg_config.OnSuccess = func(g game.Game, ctx, tctx game.Context) []game.GameTransaction {
			return ac.OnSuccess(g, ctx, tctx, ac.actionConfig(g, ctx))
		}
	}
	if ac.OnFailure != nil {
		dmg_config.OnFailure = func(g game.Game, ctx, tctx game.Context) []game.GameTransaction {
			return ac.OnFailure(g, ctx, tctx, ac.actionConfig(g, ctx))
		}
	}
	return dmg_config
}

func ResolveAttack(config AttackConfig, g game.Game, context game.Context) []game.GameTransaction {
	return game.ResolveDamageCore(config.actionConfig(g, context), config.damageConfig(), g, context)
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
				transactions := game.NewTransactionBuilder()

				if config.BeforeAttack != nil {
					transactions.Push(config.BeforeAttack(g, context, config.actionConfig(g, context)))
				}
				transactions.Push(ResolveAttack(config, g, context))
				if config.AfterAttack != nil {
					transactions.Push(config.AfterAttack(g, context, config.actionConfig(g, context)))
				}

				return transactions.Build()
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

func makeSelfStatus(
	id uuid.UUID,
	config game.ActionConfig,
	delta func(game.Game, game.Game, game.Context) []game.GameTransaction,
) game.Action {
	return game.Action{
		ID:              id,
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta:    delta,
		},
	}
}

func applySummon(context game.Context, def game.ActorDef, actions []game.Action) []game.GameTransaction {
	transactions := game.NewTransactionBuilder()

	mut := game.GameMutation{
		Delta: func(mp, g game.Game, m_ctx game.Context) game.Game {
			g.UpdateActor(*m_ctx.SourceActorID, func(a game.Actor) game.Actor {
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
			g.UpdatePlayer(*m_ctx.SourcePlayerID, func(p game.Player) game.Player {
				p.UsedSummon = true
				return p
			})
			return g
		},
	}

	transactions.PushOne(game.MakeTransaction(mut, context))
	return transactions.Build()
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

func makeRepeats(config game.DamageConfig, min int, max int, g game.Game, context game.Context) game.DamageConfig {
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
