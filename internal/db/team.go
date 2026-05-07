package db

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

type ActorConfig struct {
	AbilityID *uuid.UUID             `json:"ability_ID"`
	ActionIDs []uuid.UUID            `json:"action_IDs"`
	Focus     *game.ActorFocus       `json:"focus"`
	ItemID    *uuid.UUID             `json:"item_ID"`
	AuxStats  map[game.ActorStat]int `json:"aux_stats"`
}

type TeamActor struct {
	ActorID uuid.UUID   `json:"actor_ID"`
	Config  ActorConfig `json:"config"`
}

type TeamConfig struct {
	Actors []TeamActor `json:"actors"`
	Name   string      `json:"name"`
}
