package game

import (
	"maps"

	"github.com/google/uuid"
)

type actorResolveHandler struct {
	pre              Actor
	actor            Actor
	bypass_modifiers bool
	game             Game
	mutations        []ActorMutation
	transactions     []Transaction[Modifier]
}

func newActorResolveHandler(actor Actor, g Game, bypass_modifiers bool) actorResolveHandler {
	mutations, transactions := GetAllActorMutations(g, bypass_modifiers)
	clone := actor.Clone()
	return actorResolveHandler{
		actor:            clone,
		pre:              clone,
		game:             g,
		bypass_modifiers: bypass_modifiers,
		mutations:        mutations,
		transactions:     transactions,
	}
}

func (ah *actorResolveHandler) applyModifierMutation(mutation ActorMutation) (Actor, bool) {
	if mutation.ModifierGroupID != nil && ah.actor.HasImmunity(*mutation.ModifierGroupID) {
		return ah.actor, false
	}
	context := getContext(ah.actor, ah.transactions, mutation)
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

func toResolved(actor Actor, pre Actor) ResolvedActor {
	return ResolvedActor{
		Actor:     actor,
		BaseStats: maps.Clone(pre.Stats),
	}
}
func (ah *actorResolveHandler) resolveMutations() ResolvedActor {
	for _, mutation := range ah.mutations {
		next, did_apply := ah.applyModifierMutation(mutation)
		if !did_apply {
			continue
		}

		ah.actor = next
	}

	return toResolved(ah.actor, ah.pre)
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

		incomingMultiplier := natureResult.Result
		if incomingMultiplier == 0 {
			resolved.ResolvedNatureResistance[nature] = 0
			continue
		}

		resolved.ResolvedNatureResistance[nature] = 1.0 / incomingMultiplier
		ns := NatureSet(nature)
		resolved.ResolvedNatureDamage[nature] = GetStabModifier(*resolved, &ns)
	}
}

func (ah *actorResolveHandler) resolveActions(resolved *ResolvedActor) {
	if ah.actor.PositionID == nil {
		return
	}

	allDisabled := true
	for i, a := range resolved.Actions {
		// static cooldown offset
		baseCooldown := 0
		if a.Config.Cooldown != nil {
			baseCooldown = *a.Config.Cooldown
		}
		resolved.Actions[i].Config.Cooldown = Ptr(baseCooldown + resolved.CooldownOffset)

		// set dynamic disabled state
		if !a.State.Disabled {
			context := MakeContextForActor(ah.actor)
			context.ActionID = &a.ID

			if resolved.ActionLocked && resolved.LastUsedActionTX != nil {
				if resolved.LastUsedActionTX.Mutation.ID != a.ID && !a.Config.Switch {
					resolved.Actions[i].State.Disabled = true
				}
			}
			if resolved.SwitchLocked && a.Config.Switch {
				resolved.Actions[i].State.Disabled = true
			}

			filterGame := ah.game.WithActor(ah.actor).WithoutActionFilterEval()
			if !ah.game.disableActionFilterEval && !a.Filter(ah.game, filterGame, context) {
				resolved.Actions[i].State.Disabled = true
			}

			player, ok := ah.game.GetPlayerByID(resolved.PlayerID)
			if ok && player.UsedSummon && a.Config.Summon {
				resolved.Actions[i].State.Disabled = true
			}

			_, ok = ah.game.FindQueuedAction(func(t Transaction[Action]) bool {
				if t.Context.SourcePlayerID == nil {
					return false
				}
				is_player := *t.Context.SourcePlayerID == resolved.PlayerID
				return is_player && t.Mutation.Config.Summon
			})

			if ok && a.Config.Summon {
				resolved.Actions[i].State.Disabled = true
			}

			if !a.Config.Switch && !resolved.Actions[i].State.Disabled {
				allDisabled = false
			}
		}
	}

	// if everything is disabled, add struggle
	if allDisabled {
		resolved.Actions = append(resolved.Actions, MakeStruggle())
	}
}
