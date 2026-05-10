package instance

import (
	"shinobi_showdown/internal/chat"
	"shinobi_showdown/internal/db"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data"

	"github.com/google/uuid"
)

type RequestType = string

const (
	SetTeam        RequestType = "set-team"
	Reset          RequestType = "reset"
	ReadyTeam      RequestType = "ready-team"
	CancelTeam     RequestType = "cancel-team"
	StartBattle    RequestType = "start-battle"
	PushAction     RequestType = "push-action"
	RemoveAction   RequestType = "remove-action"
	RunGameActions RequestType = "run-game-actions" // TEMP
	ResolvePrompt  RequestType = "resolve-prompt"

	GetTargets      RequestType = "get-targets"
	ValidateContext RequestType = "validate-context"

	SendChat RequestType = "send-chat"
)

type Request struct {
	Type        RequestType     `json:"type"`
	ClientID    uuid.UUID       `json:"client_ID"`
	PromptID    *uuid.UUID      `json:"prompt_ID"`
	Context     game.Context    `json:"context"`
	ActorConfig *db.ActorConfig `json:"actor_config"`
	TeamConfig  *db.TeamConfig  `json:"team_config"`
	ChatMessage *chat.Message   `json:"chat_message"`
}

type HydratedActorConfig struct {
	Ability  *game.Modifier
	Actions  []game.Action
	Focus    game.ActorFocus
	Item     *game.Modifier
	AuxStats map[game.ActorStat]int `json:"aux_stats"`
}

func HydrateActorConfig(config db.ActorConfig, abilities []game.Modifier) HydratedActorConfig {
	var ability *game.Modifier = nil
	if config.AbilityID != nil {
		for _, a := range abilities {
			if a.ID == *config.AbilityID {
				ability = &a
			}
		}
	} else {
		if len(abilities) > 0 {
			ability = &abilities[0]
		}
	}

	actions := make([]game.Action, 0, len(config.ActionIDs))
	for _, id := range config.ActionIDs {
		action, ok := data.ACTIONS[id]
		if !ok {
			continue
		}
		actions = append(actions, action)
	}

	var item *game.Modifier = nil
	if config.ItemID != nil {
		i, ok := data.ITEMS[*config.ItemID]
		if ok {
			item = &i
		}
	}

	focus := game.FocusNone
	if config.Focus != nil {
		focus = *config.Focus
	}

	return HydratedActorConfig{
		Ability:  ability,
		Actions:  actions,
		AuxStats: config.AuxStats,
		Focus:    focus,
		Item:     item,
	}
}
