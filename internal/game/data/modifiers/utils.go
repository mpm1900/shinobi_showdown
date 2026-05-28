package modifiers

import (
	"fmt"
	"math/rand/v2"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"
	"slices"

	"github.com/google/uuid"
)

func RemoveModifierSource(id uuid.UUID, actor game.Actor) []game.GameTransaction {
	transactions := game.NewTransactionBuilder()

	ctx := game.MakeContextForActor(actor)
	ctx.ModifierID = &id
	mod := mutations.RemoveModifierWhere(func(t game.Transaction[game.Modifier]) bool {
		if !((t.Mutation.ID == id) || (t.Mutation.GroupID != nil && *t.Mutation.GroupID == id)) {
			return false
		}

		if slices.Contains(t.Context.TargetActorIDs, actor.ID) {
			return true
		}
		if actor.PositionID != nil && slices.Contains(t.Context.TargetPositionIDs, *actor.PositionID) {
			return true
		}

		return false
	})
	mod_tx := game.MakeTransaction(mod, ctx)
	transactions.PushOne(mod_tx)

	return transactions.Build()
}
func applyModifier(checkWarded bool, config game.ActionConfig, context game.Context, actor game.Actor, modifier game.Modifier) []game.GameTransaction {
	transactions := game.NewTransactionBuilder()

	if mutations.CheckJutsuImmunity(config, actor) {
		log_ctx := game.MakeContextForActor(actor)
		log := game.MakeGameLog(fmt.Sprintf("$source$ was immune to %s.", config.Jutsu), log_ctx, 1)
		tx := game.AddLogs(log)
		transactions.PushOne(game.MakeTransaction(tx, log_ctx))

		return transactions.Build()
	}

	ctx := game.MakeContextForActor(actor)
	ctx.ModifierID = modifier.GroupID
	mod := mutations.AddModifiers(checkWarded, modifier)
	mod_tx := game.MakeTransaction(mod, ctx)
	transactions.PushOne(mod_tx)

	return transactions.Build()
}
func ApplyModifier(config game.ActionConfig, context game.Context, actor game.Actor, modifier game.Modifier) []game.GameTransaction {
	return applyModifier(true, config, context, actor, modifier)
}
func ApplyModifierBypass(config game.ActionConfig, context game.Context, actor game.Actor, modifier game.Modifier) []game.GameTransaction {
	return applyModifier(false, config, context, actor, modifier)
}
func getModifierChance(g game.Game, actor game.Actor, chance int) (bool, int) {
	resolved := actor.Resolve(g)
	b := resolved.ModifierChanceOffset
	m := resolved.ModifierChanceMult
	y := game.Round(m*float64(chance)) + b
	roll := rand.IntN(100)
	fmt.Printf("modifier: roll=%d, check=%d\n", roll, y)
	return y >= roll, roll
}
func getStatusChance(g game.Game, actor game.Actor, chance int) (bool, int) {
	resolved := actor.Resolve(g)
	b := resolved.StatusChanceOffset
	m := resolved.StatusChanceMult
	y := game.Round(m*float64(chance)) + b
	roll := rand.IntN(100)
	fmt.Printf("status: roll=%d, check=%d\n", roll, y)
	return y >= roll, roll
}
func ChanceModifier(config game.ActionConfig, g game.Game, context game.Context, actor game.Actor, modifier game.Modifier, chance int) []game.GameTransaction {
	source, ok := g.GetSource(context)
	if !ok {
		return []game.GameTransaction{}
	}
	success, _ := getModifierChance(g, source, chance)
	if !success {
		return []game.GameTransaction{}
	}

	return ApplyModifier(config, context, actor, modifier)
}

func applyStatus(checkWarded bool, config game.ActionConfig, g game.Game, actor game.Actor, modifier game.Modifier, mutation game.GameMutation, actionCtx game.Context) []game.GameTransaction {
	transactions := game.NewTransactionBuilder()

	resolved := actor.Resolve(g)
	if resolved.Safeguarded {
		log_ctx := game.MakeContextForActor(actor)
		log := game.MakeGameLog("$source$ was safeguarded.", log_ctx, 1)
		tx := game.AddLogs(log)
		transactions.PushOne(game.MakeTransaction(tx, log_ctx))

		return transactions.Build()
	}
	if mutations.CheckJutsuImmunity(config, actor) {
		log_ctx := game.MakeContextForActor(actor)
		log := game.MakeGameLog(fmt.Sprintf("$source$ was immune to %s.", config.Jutsu), log_ctx, 1)
		tx := game.AddLogs(log)
		transactions.PushOne(game.MakeTransaction(tx, log_ctx))

		return transactions.Build()
	}
	if mutations.CheckImmunity(modifier.ID, actor) {
		log_ctx := game.MakeContextForActor(actor)
		log := game.MakeGameLog("$source$ was immune.", log_ctx, 1)
		tx := game.AddLogs(log)
		transactions.PushOne(game.MakeTransaction(tx, log_ctx))

		return transactions.Build()
	}
	if resolved.Statused && config.Name != "" {
		log_ctx := game.MakeContextForActor(actor)
		log := game.MakeGameLog(fmt.Sprintf("%s failed.", config.Name), log_ctx, 1)
		tx := game.AddLogs(log)
		transactions.PushOne(game.MakeTransaction(tx, log_ctx))

		return transactions.Build()
	}

	ctx := game.NewContext()
	ctx.ActionID = actionCtx.ActionID
	ctx.SourcePlayerID = actionCtx.SourcePlayerID
	ctx.SourceActorID = actionCtx.SourceActorID
	ctx.TargetActorIDs = []uuid.UUID{actor.ID}
	ctx.ParentActorID = nil // do not remove on switch

	mod := mutations.AddStatus(checkWarded, modifier)
	mod_tx := game.MakeTransaction(mod, ctx)

	mut_tx := game.MakeTransaction(mutation, ctx)
	transactions.PushOne(mod_tx)
	transactions.PushOne(mut_tx)

	return transactions.Build()
}
func ApplyStatus(config game.ActionConfig, g game.Game, actor game.Actor, modifier game.Modifier, mutation game.GameMutation, actionCtx game.Context) []game.GameTransaction {
	return applyStatus(true, config, g, actor, modifier, mutation, actionCtx)
}
func ApplyStatusBypass(config game.ActionConfig, g game.Game, actor game.Actor, modifier game.Modifier, mutation game.GameMutation, actionCtx game.Context) []game.GameTransaction {
	return applyStatus(false, config, g, actor, modifier, mutation, actionCtx)
}
func ClearStatus(g game.Game, actor game.Actor) []game.GameTransaction {
	transactions := game.NewTransactionBuilder()

	rmod_mut := mutations.RemoveModifierWhere(func(t game.Transaction[game.Modifier]) bool {
		targets := g.GetTargets(t.Context)
		remove := false
		for _, t := range targets {
			if t.ID == actor.ID {
				remove = true
				break
			}
		}

		return remove && t.Mutation.Status
	})
	rmod_tx := game.MakeTransaction(rmod_mut, game.MakeContextForActor(actor))
	transactions.PushOne(rmod_tx)

	rmut_mut := game.GameMutation{
		Delta: func(p, g game.Game, context game.Context) game.Game {
			g.UpdateActor(*context.SourceActorID, func(a game.Actor) game.Actor {
				a.Statused = false
				a.Burned = false
				a.Sleeping = false
				a.Poisoned = false
				a.Paralyzed = false
				return a
			})
			return g
		},
	}
	rmut_tx := game.MakeTransaction(rmut_mut, game.MakeContextForActor(actor))
	transactions.PushOne(rmut_tx)

	return transactions.Build()
}

func ApplyBurn(config game.ActionConfig, g game.Game, actor game.Actor, actionCtx game.Context) []game.GameTransaction {
	_, ok := actor.Natures[game.NsFire]
	if ok {
		// fire nature immune to burn
		return []game.GameTransaction{}
	}

	return ApplyStatus(config, g, actor, Burned, mutations.Burn, actionCtx)
}
func ChanceBurn(config game.ActionConfig, g game.Game, context game.Context, actor game.Actor, chance int) []game.GameTransaction {
	source, ok := g.GetSource(context)
	if !ok {
		return []game.GameTransaction{}
	}
	success, _ := getStatusChance(g, source, chance)
	if !success {
		return []game.GameTransaction{}
	}

	return ApplyBurn(config, g, actor, context)
}

func ApplyParalysis(config game.ActionConfig, g game.Game, actor game.Actor, actionCtx game.Context) []game.GameTransaction {
	_, ok := actor.Natures[game.NsLightning]
	if ok {
		// lightning nature immune to paralysis
		return []game.GameTransaction{}
	}
	return ApplyStatus(config, g, actor, Paralysis, mutations.Paralyze, actionCtx)
}
func ChanceParalysis(config game.ActionConfig, g game.Game, context game.Context, actor game.Actor, chance int) []game.GameTransaction {
	source, ok := g.GetSource(context)
	if !ok {
		return []game.GameTransaction{}
	}
	success, _ := getStatusChance(g, source, chance)
	if !success {
		return []game.GameTransaction{}
	}

	return ApplyParalysis(config, g, actor, context)
}

func ApplySleep(config game.ActionConfig, g game.Game, actor game.Actor, actionCtx game.Context) []game.GameTransaction {
	return ApplyStatus(config, g, actor, Sleeping, mutations.Sleep, actionCtx)
}
func ChanceSleep(config game.ActionConfig, g game.Game, context game.Context, actor game.Actor, chance int) []game.GameTransaction {
	source, ok := g.GetSource(context)
	if !ok {
		return []game.GameTransaction{}
	}
	success, _ := getStatusChance(g, source, chance)
	if !success {
		return []game.GameTransaction{}
	}

	return ApplySleep(config, g, actor, context)
}

func ApplyPoison(config game.ActionConfig, g game.Game, actor game.Actor, actionCtx game.Context) []game.GameTransaction {
	return ApplyStatus(config, g, actor, Poisoned, mutations.Poison, actionCtx)
}
func ChancePoison(config game.ActionConfig, g game.Game, context game.Context, actor game.Actor, chance int) []game.GameTransaction {
	source, ok := g.GetSource(context)
	if !ok {
		return []game.GameTransaction{}
	}
	success, _ := getStatusChance(g, source, chance)
	if !success {
		return []game.GameTransaction{}
	}

	return ApplyPoison(config, g, actor, context)
}

func ApplyImmunity(modifier_id uuid.UUID, immunity_id uuid.UUID) game.ActorMutation {
	return game.MakeActorMutation(
		&modifier_id,
		game.MutPriorityDefault,
		game.SourceFilter,
		func(g game.Game, a game.Actor, ctx game.Context) game.Actor {
			if a.Immunities == nil {
				a.Immunities = map[uuid.UUID]struct{}{}
			}
			a.Immunities[immunity_id] = struct{}{}
			return a
		})
}
func applyWeather(g game.Game, context game.Context, weather game.GameWeather, mod game.Modifier, duration int) []game.GameTransaction {
	transactions := game.NewTransactionBuilder()

	state, _ := g.GetState(context)
	if state.Weather == weather {
		return transactions.Build()
	}

	transactions.PushOne(FilterWeather())

	mod.Duration = duration
	mut := mutations.AddModifiers(false, mod)
	transactions.PushOne(game.MakeTransaction(mut, game.NewContext()))

	return transactions.Build()
}

func ApplyRain(g game.Game, context game.Context) []game.GameTransaction {
	return applyWeather(g, context, game.GameWeatherRain, RainWeather(), 4)
}
func ApplySandstorm(g game.Game, context game.Context) []game.GameTransaction {
	return applyWeather(g, context, game.GameWeatherSandstorm, SandstormWeather(), 4)
}
func ApplySunlight(g game.Game, context game.Context) []game.GameTransaction {
	return applyWeather(g, context, game.GameWeatherSunlight, SunlightWeather(), 4)
}
