package game

import (
	"fmt"

	"github.com/google/uuid"
)

type PlayerPosition struct {
	ID      uuid.UUID  `json:"ID"`
	ActorID *uuid.UUID `json:"actor_ID"`
}

type Player struct {
	ID                uuid.UUID        `json:"ID"`
	User              User             `json:"user"`
	PositionsCapacity int              `json:"positions_capacity"`
	Positions         []PlayerPosition `json:"positions"`
	TeamCapacity      int              `json:"team_capacity"`
	UsedSummon        bool             `json:"used_summon"`
	Ready             bool             `json:"ready"`
}

func NewPlayer(ID uuid.UUID, capacity int, user User) Player {
	positions := make([]PlayerPosition, capacity)
	for i := range capacity {
		positions[i] = PlayerPosition{
			ID:      uuid.New(),
			ActorID: nil,
		}
	}

	return Player{
		ID:                ID,
		User:              user,
		PositionsCapacity: 2,
		Positions:         positions,
		TeamCapacity:      6,
		UsedSummon:        false,
	}
}

func (p *Player) Reset() {
	p.Ready = false
	for i, _ := range p.Positions {
		p.Positions[i].ActorID = nil
	}
}

func (p Player) HasPosition(pid uuid.UUID) bool {
	for _, position := range p.Positions {
		if position.ID == pid {
			return true
		}
	}

	return false
}

func (p Player) GetActorAtPosition(pid uuid.UUID) *uuid.UUID {
	for _, position := range p.Positions {
		if position.ID == pid {
			return position.ActorID
		}
	}
	return nil
}

func (p *Player) SetPosition(pid uuid.UUID, aid *uuid.UUID) {
	for i, pos := range p.Positions {
		if pos.ID == pid {
			p.Positions[i].ActorID = aid
		}
	}
}

func (p *Player) AddPosition(aid *uuid.UUID) error {
	for i, pos := range p.Positions {
		if pos.ActorID == nil {
			p.Positions[i].ActorID = aid
			return nil
		}
	}

	return fmt.Errorf("no room in positions")
}
