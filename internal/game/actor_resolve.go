package game

import (
	"maps"

	"github.com/google/uuid"
)

func MapBaseStat(stat, level int, focus float64, ev int) int {
	base := float64((stat * 2) + BASE_IV)
	ratio := float64((base+(float64(ev*2)))*float64(level)) / 100
	return Round((ratio + 5) * focus)
}

func MapResourceStat(stat, level int, focus float64, ev int) int {
	return MapBaseStat(stat, level, focus, ev) + level + 5
}

func (actor *Actor) MapBase(stat ActorStat) {
	actor.Stats[stat] = MapBaseStat(actor.Stats[stat], actor.Level, actor.GetFocusModifier(stat), actor.AuxStats[stat])
}

func (actor *Actor) MapResource(stat ActorStat) {
	actor.Stats[stat] = MapResourceStat(actor.Stats[stat], actor.Level, 1.0, actor.AuxStats[stat])
}

func (actor *Actor) MapBaseStats() {
	actor.MapResource(StatHP)
	actor.MapResource(StatStamina)

	actor.MapBase(StatAttack)
	actor.MapBase(StatDefense)
	actor.MapBase(StatChakraAttack)
	actor.MapBase(StatChakraDefense)
	actor.MapBase(StatSpeed)
}

func MapStagedStat(stat, stage, mod int) int {
	m := 1.0
	if stage > 0 {
		stage = min(stage, 6)
		m = float64(stage+mod) / float64(mod)
	} else if stage < 0 {
		stage = max(stage, -6)
		m = float64(mod) / float64(-stage+mod)
	}

	return Round(float64(stat) * m)
}

func (actor *Actor) MapStaged(stat ActorStat, mod int) {
	actor.Stats[stat] = MapStagedStat(actor.Stats[stat], actor.Stages[stat], mod)
}

func (actor *Actor) MapStagedStats() {
	actor.MapStaged(StatAttack, 2)
	actor.MapStaged(StatDefense, 2)
	actor.MapStaged(StatChakraAttack, 2)
	actor.MapStaged(StatChakraDefense, 2)
	actor.MapStaged(StatSpeed, 2)

	actor.MapStaged(StatEvasion, 3)
	actor.MapStaged(StatAccuracy, 3)
}

func newActorContext(actor *Actor) Context {
	return Context{
		SourcePlayerID:    &actor.PlayerID,
		SourceActorID:     &actor.ID,
		ParentActorID:     &actor.ID,
		TargetActorIDs:    []uuid.UUID{},
		TargetPositionIDs: []uuid.UUID{},
	}
}

func GetActorModifiers(game Game) []Transaction[Modifier] {
	activeActors := game.GetActiveActors()
	modifiers := make([]Transaction[Modifier], 0, len(activeActors)*3)

	for i := range activeActors {
		actor := &activeActors[i]
		context := newActorContext(actor)
		ability := actor.GetAbility()
		if ability != nil {
			modifiers = append(modifiers, MakeTransaction(*ability, context))
		}
		if actor.Item != nil {
			modifiers = append(modifiers, MakeTransaction(*actor.Item, context))
		}

		for j := range actor.ActorDef.DefaultModifiers {
			modifiers = append(modifiers, MakeTransaction(actor.ActorDef.DefaultModifiers[j], context))
		}
	}

	return modifiers
}

var specialMutations = []ActorMutation{
	MakeActorMutation(
		nil,
		MutPriorityMapBaseStats,
		AllFilter,
		func(g Game, input Actor, context Context) Actor {
			input.MapBaseStats()
			return input
		},
	),
	MakeActorMutation(
		nil,
		MutPriorityMapStagedStats,
		AllFilter,
		func(g Game, input Actor, context Context) Actor {
			input.MapStagedStats()
			return input
		},
	),
}

func getContext(actor *Actor, transactions []Transaction[Modifier], mutation ActorMutation) Context {
	context := newActorContext(actor)
	return ResolveModifierTransactionContext(context, transactions, mutation.TransactionID)
}

func resolveActor(actor Actor, g Game, mutations []ActorMutation, transactions []Transaction[Modifier], actions func(*ResolvedActor)) ResolvedActor {
	handler := actorResolveHandler{
		actor:        actor,
		baseStats:    maps.Clone(actor.Stats),
		game:         g,
		mutations:    mutations,
		transactions: transactions,
	}
	return handler.resolve(actions)
}

func (a Actor) getActor() Actor {
	actor := a
	if actor.Summon != nil && !actor.Summon.Proxy {
		form := actor.Summon.Actor
		actor.ActorDef = form.ActorDef
		actor.Actions = form.Actions
	}

	return actor
}

func (a Actor) Resolve(g Game) ResolvedActor {
	mutations, transactions := GetAllActorMutations(g, false)
	return a.ResolveWithMutations(g, mutations, transactions)
}

func (a Actor) ResolveWithMutations(g Game, mutations []ActorMutation, transactions []Transaction[Modifier]) ResolvedActor {
	actor := a.getActor()
	resolved := resolveActor(actor.Clone(), g, mutations, transactions, func(ra *ResolvedActor) {
		resolveActions(g, ra)
	})

	unmodified := actor.Clone()
	unmodified.MapBaseStats()
	unmodified.MapStagedStats()

	resolved.UnmodifiedStats = unmodified.Stats
	resolved.Ability = actor.GetAbility()

	return resolved
}

func (a Actor) ResolveShallow(g Game) ResolvedActor {
	actor := a.getActor()
	mutations, transactions := GetAllActorMutations(g, false)
	resolved := resolveActor(actor.Clone(), g, mutations, transactions, nil)

	return resolved
}

type actorResolveHandler struct {
	actor        Actor
	baseStats    map[ActorStat]int
	game         Game
	mutations    []ActorMutation
	transactions []Transaction[Modifier]
}

func (ah *actorResolveHandler) applyModifierMutation(mutation ActorMutation) (Actor, bool) {
	if mutation.ModifierGroupID != nil && ah.actor.HasImmunity(*mutation.ModifierGroupID) {
		return ah.actor, false
	}
	context := getContext(&ah.actor, ah.transactions, mutation)
	g := ah.game.WithActor(ah.actor)

	tx := MakeTransaction(mutation.Mutation, context)
	next, ok := ResolveTransaction(g, ah.actor, tx, ah.actor)
	if !ok {
		return ah.actor, false
	}

	if mutation.ModifierGroupID == nil {
		return next, true
	}

	if next.AppliedModifiers == nil {
		next.AppliedModifiers = make(map[uuid.UUID]int)
	}

	count, ok := next.AppliedModifiers[*mutation.ModifierGroupID]
	if !ok {
		next.AppliedModifiers[*mutation.ModifierGroupID] = 1
		return next, true
	}

	next.AppliedModifiers[*mutation.ModifierGroupID] = count + 1
	return next, true
}

func toResolved(actor Actor, baseStats map[ActorStat]int) ResolvedActor {
	return ResolvedActor{
		Actor:     actor,
		BaseStats: baseStats,
	}
}

func (ah *actorResolveHandler) resolve(actions func(*ResolvedActor)) ResolvedActor {
	resolved := ah.resolveMutations()
	ah.resolveNatures(&resolved)
	if actions != nil {
		actions(&resolved)
	}
	return resolved
}

func (ah *actorResolveHandler) resolveMutations() ResolvedActor {
	for _, mutation := range ah.mutations {
		next, did_apply := ah.applyModifierMutation(mutation)
		if !did_apply {
			continue
		}

		ah.actor = next
	}

	return toResolved(ah.actor, ah.baseStats)
}

func (ah *actorResolveHandler) resolveNatures(resolved *ResolvedActor) {
	resolved.ResolvedNatureResistance = make(map[Nature]float64)
	resolved.ResolvedNatureDamage = make(map[Nature]float64)
	for nature := range resolved.NatureResistance {
		natureResult := ResolveNatures(
			[]Nature{nature},
			NewNatureSetValues(),
			resolved.NatureResistance,
			resolved.Natures,
		)

		if natureResult.Result == 0 {
			resolved.ResolvedNatureResistance[nature] = 0
			continue
		}

		resolved.ResolvedNatureResistance[nature] = 1.0 / natureResult.Result
		ns := NatureSet(nature)
		resolved.ResolvedNatureDamage[nature] = GetStabModifier(*resolved, &ns)
	}
}

func resolveActions(game Game, resolved *ResolvedActor) {
	if resolved.PositionID == nil {
		return
	}

	player, _ := game.GetPlayerByID(resolved.PlayerID)
	_, hasQueuedSummon := game.FindQueuedAction(func(t Transaction[Action]) bool {
		if t.Context.SourcePlayerID == nil {
			return false
		}
		return *t.Context.SourcePlayerID == resolved.PlayerID && t.Mutation.Config.Summon
	})

	filterGame := game.WithActor(resolved.Actor).WithoutActionFilterEval()
	allDisabled := true
	context := newActorContext(&resolved.Actor)

	for i := range resolved.Actions {
		action := &resolved.Actions[i]
		if action.State.Disabled {
			continue
		}

		// static cooldown offset
		baseCooldown := 0
		if action.Config.Cooldown != nil {
			baseCooldown = *action.Config.Cooldown
		}
		action.Config.Cooldown = Ptr(baseCooldown + resolved.CooldownOffset)

		// set dynamic disabled state
		context.ActionID = &action.ID

		if resolved.ActionLocked && resolved.LastUsedActionTX != nil {
			if resolved.LastUsedActionTX.Mutation.ID != action.ID && !action.Config.Switch {
				action.State.Disabled = true
			}
		}
		if resolved.SwitchLocked && action.Config.Switch {
			action.State.Disabled = true
		}

		if !game.disableActionFilterEval && !action.Filter(game, filterGame, context) {
			action.State.Disabled = true
		}

		if player.UsedSummon && action.Config.Summon {
			action.State.Disabled = true
		}

		if hasQueuedSummon && action.Config.Summon {
			action.State.Disabled = true
		}

		if !action.Config.Switch && !action.State.Disabled {
			allDisabled = false
		}
	}

	// if everything is disabled, add struggle
	if allDisabled {
		resolved.Actions = append(resolved.Actions, MakeStruggle())
	}
}
