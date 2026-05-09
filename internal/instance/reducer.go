package instance

import (
	"shinobi_showdown/internal/game"
	data "shinobi_showdown/internal/game/data"
	"slices"

	"github.com/google/uuid"
)

func getTargets(instance *Instance, request Request) int {
	if request.Context.ActionID == nil && request.PromptID == nil {
		instance.TargetIDsResponse(request.ClientID, request.Context, nil)
		return none
	}

	var action = game.Action{}
	if request.PromptID != nil {
		tx, ok := instance.Game.GetPromptTxByID(*request.PromptID)
		if !ok {
			instance.TargetIDsResponse(request.ClientID, request.Context, nil)
			return none
		}
		action = tx.Mutation
	} else {
		actor, ok := instance.Game.GetSource(request.Context)
		if !ok {
			instance.ValidateContextResponse(request.ClientID, request.Context, false)
			return none
		}

		a, ok := actor.GetActionByID(instance.Game, *request.Context.ActionID)
		if !ok {
			// if it's not on the source, do a root look up
			// this is probably due to something like recharing or a
			// non-stored, special action
			a, ok = data.ACTIONS[*request.Context.ActionID]
			if !ok {
				instance.TargetIDsResponse(request.ClientID, request.Context, nil)
				return none
			}
		}

		action = a
	}

	targets := instance.Game.GetActors(func(a game.Actor) bool {
		return action.TargetPredicate(instance.Game, a, request.Context)
	})
	targetIDs := make([]uuid.UUID, 0, len(targets))
	for _, a := range targets {
		targetIDs = append(targetIDs, a.ID)
	}

	instance.TargetIDsResponse(request.ClientID, request.Context, targetIDs)
	return none
}
func validateContext(instance *Instance, request Request) int {
	if request.Context.ActionID == nil && request.PromptID == nil {
		instance.ValidateContextResponse(request.ClientID, request.Context, false)
		return none
	}

	if request.PromptID != nil {
		action, ok := instance.Game.GetPromptTxByID(*request.PromptID)
		if !ok {
			instance.ValidateContextResponse(request.ClientID, request.Context, false)
			return none
		}

		valid := action.Mutation.ContextValidate(request.Context)
		instance.ValidateContextResponse(request.ClientID, request.Context, valid)
		return none
	}

	actor, ok := instance.Game.GetSource(request.Context)
	if !ok {
		instance.ValidateContextResponse(request.ClientID, request.Context, false)
		return none
	}

	action, ok := actor.GetActionByID(instance.Game, *request.Context.ActionID)
	if !ok {
		instance.ValidateContextResponse(request.ClientID, request.Context, false)
		return none
	}

	valid := action.ContextValidate(request.Context)
	instance.ValidateContextResponse(request.ClientID, request.Context, valid)
	return none
}
func setTeam(instance *Instance, request Request) int {
	if request.TeamConfig == nil {
		return none
	}
	player, ok := instance.Game.GetPlayerByID(request.ClientID)
	if !ok {
		return none
	}

	if len(request.TeamConfig.Actors) > player.TeamCapacity {
		return none
	}

	teamActors := make([]game.Actor, len(request.TeamConfig.Actors))
	for i, actorConfig := range request.TeamConfig.Actors {
		def, ok := data.ACTORS[actorConfig.ActorID]
		if !ok {
			return none
		}

		hydrated := HydrateActorConfig(actorConfig.Config, def.Abilities)

		actor := game.MakeActor(
			def,
			request.ClientID,
			game.LV_100_XP,
			hydrated.Ability,
			hydrated.Item,
			hydrated.Actions,
			hydrated.Focus,
			hydrated.AuxStats,
		)
		teamActors[i] = actor
	}
	instance.Game.SetPlayerActors(request.ClientID, teamActors)
	return state
}
func readyTeam(instance *Instance, request Request) int {
	instance.Game.EnableActors(request.ClientID, request.Context.TargetActorIDs)
	return state
}
func cancelTeam(instance *Instance, request Request) int {
	instance.Game.ResetPlayer(request.ClientID)
	return state
}
func startBattle(instance *Instance) int {
	if instance.Game.Status == game.GameStatusRunning {
		return none
	}

	actors := make([]game.Actor, 0)
	for _, a := range instance.Game.Actors {
		if a.Enabled {
			actors = append(actors, a)
		}
	}

	instance.Game.Actors = actors
	instance.RunGameActions()
	return state
}
func pushAction(instance *Instance, request Request) int {
	if request.Context.ActionID == nil || request.Context.SourceActorID == nil {
		return none
	}

	if request.Context.ParentActorID == nil {
		request.Context.ParentActorID = request.Context.SourceActorID
	}

	actor, ok := instance.Game.GetSource(request.Context)
	if !ok {
		return none
	}

	action, ok := actor.GetActionByID(instance.Game, *request.Context.ActionID)
	if !ok {
		return none
	}

	if action.Meta.Switch && actor.SwitchLocked {
		return none
	}

	if !action.Meta.Switch && actor.ActionLocked && actor.LastUsedActionTX != nil {
		if *request.Context.ActionID != actor.LastUsedActionTX.Mutation.ID {
			return none
		}
	}

	transaction := game.MakeTransaction(action, request.Context)
	if instance.Game.PushAction(transaction) {
		instance.RunGameActions()
	}

	return state
}
func cancelAction(instance *Instance, request Request) int {
	if request.Context.ActionID == nil {
		return none
	}

	instance.Game.Actions = slices.DeleteFunc(instance.Game.Actions, func(tx game.Transaction[game.Action]) bool {
		if tx.Context.SourcePlayerID != nil && *tx.Context.SourcePlayerID == request.ClientID {
			return tx.ID == *request.Context.ActionID
		}

		return false
	})

	return state
}
func resolvePrompt(instance *Instance, request Request) int {
	if request.PromptID == nil {
		return none
	}

	instance.Game.ReadyPrompt(*request.PromptID, request.Context)
	if instance.Game.AllPromptsReady() {
		instance.RunGameActions()
	}

	return state

}

func Reducer(instance *Instance, request Request) int {
	switch request.Type {
	case GetTargets:
		return getTargets(instance, request)
	case ValidateContext:
		return validateContext(instance, request)
	case SetTeam:
		return setTeam(instance, request)
	case ReadyTeam:
		return readyTeam(instance, request)
	case CancelTeam:
		return cancelTeam(instance, request)
	case StartBattle:
		return startBattle(instance)
	case Reset:
		instance.Game.Reset()
		return state
	case PushAction:
		return pushAction(instance, request)
	case RemoveAction: //CancelAction
		return cancelAction(instance, request)
	case ResolvePrompt:
		return resolvePrompt(instance, request)
	case RunGameActions:
		if instance.Game.Status == game.GameStatusRunning {
			return none
		}

		instance.RunGameActions()
		return state

	default:
		return none
	}
}
