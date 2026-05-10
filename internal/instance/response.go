package instance

import (
	"shinobi_showdown/internal/chat"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

type ResponseType string

const (
	ResponseTypeGame            = "game"
	ResponseTypeClients         = "clients"
	ResponseTypeJoinSuccess     = "join-success"
	ResponseTypeValidateContext = "validate-context"
	ResponseTypeTargetIDs       = "target-IDs"

	ResponseTypeNewChat = "new-chat"
)

type Response struct {
	Type        ResponseType   `json:"type"`
	State       *game.GameJSON `json:"state"`
	Clients     []*Client      `json:"clients"`
	Valid       *bool          `json:"valid"`
	Context     *game.Context  `json:"context"`
	ChatMessage *chat.Message  `json:"chat_message"`
}

func NewGameMessage(client *Client, state *game.Game) Response {
	var json game.GameJSON
	if client == nil {
		json = state.ToJSON(nil)
	} else {
		json = state.ToJSON(&client.ID)
	}

	return Response{
		Type:        ResponseTypeGame,
		State:       &json,
		Clients:     nil,
		ChatMessage: nil,
	}
}

func NewClientsMessage(clients []*Client) Response {
	return Response{
		Type:        ResponseTypeClients,
		State:       nil,
		Clients:     clients,
		ChatMessage: nil,
	}
}

func NewChatMessage(message chat.Message) Response {
	return Response{
		Type:        ResponseTypeNewChat,
		State:       nil,
		Clients:     nil,
		ChatMessage: &message,
	}
}

func PostRegisterMessage(client *Client, state *game.Game) Response {
	var json game.GameJSON
	if client == nil {
		json = state.ToJSON(nil)
	} else {
		json = state.ToJSON(&client.ID)
	}

	return Response{
		Type:    ResponseTypeJoinSuccess,
		State:   &json,
		Clients: []*Client{client},
	}
}

func TargetIDsResponse(client *Client, context game.Context, targetIDs []uuid.UUID) Response {
	context.TargetActorIDs = targetIDs
	return Response{
		Type:    ResponseTypeTargetIDs,
		Context: &context,
	}
}

func ValidateContextMessage(client *Client, context game.Context, valid bool) Response {
	return Response{
		Type:    ResponseTypeValidateContext,
		Context: &context,
		Valid:   &valid,
	}
}
