package instance

import (
	"shinobi_showdown/internal/game"
	data "shinobi_showdown/internal/game/data"
	"slices"

	"github.com/google/uuid"
)

func findAction(instance *Instance, request Request) (game.Action, bool) {
	if request.PromptID != nil {
		tx, ok := instance.Game.GetPromptTxByID(*request.PromptID)
		if ok {
			return tx.Mutation, true
		}
		return game.Action{}, false
	}

	if request.Context.ActionID == nil {
		return game.Action{}, false
	}

	actor, ok := instance.Game.GetSource(request.Context)
	if !ok {
		return game.Action{}, false
	}

	action, ok := actor.GetActionByID(instance.Game, *request.Context.ActionID)
	if !ok {
		action, ok = data.ACTIONS[*request.Context.ActionID]
	}

	return action, ok
}

func getTargets(instance *Instance, request Request) int {
	action, ok := findAction(instance, request)
	if !ok {
		instance.TargetIDsResponse(request.ClientID, request.Context, nil)
		return none
	}

	targetIDs := make([]uuid.UUID, 0, len(instance.Game.Actors))
	for i := range instance.Game.Actors {
		if action.TargetPredicate(instance.Game, instance.Game.Actors[i], request.Context) {
			targetIDs = append(targetIDs, instance.Game.Actors[i].ID)
		}
	}

	instance.TargetIDsResponse(request.ClientID, request.Context, targetIDs)
	return none
}

func validateContext(instance *Instance, request Request) int {
	action, ok := findAction(instance, request)
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

	if action.State.Disabled {
		return none
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

	case SendChat:
		if request.ChatMessage == nil {
			return none
		}

		request.ChatMessage.ID = uuid.New()
		instance.BroadcastChatMessage(*request.ChatMessage)
		return none

	default:
		return none
	}
}
