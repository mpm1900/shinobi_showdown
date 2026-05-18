package game

import (
	"fmt"

	"github.com/google/uuid"
)

func queueSwitchPrompt(g *Game, player Player, missing_pos []uuid.UUID) bool {
	action := SwitchIn(len(missing_pos))

	context := NewContext()
	context.SourcePlayerID = &player.ID
	context.TargetPositionIDs = missing_pos
	possible_targets := g.GetActors(func(a Actor) bool {
		return action.TargetPredicate(*g, a, context)
	})

	if len(possible_targets) == 0 {
		fmt.Printf("Invalid state, but no possible targets, likely game-over. \n")
		return true
	}

	switch_count := min(len(missing_pos), len(possible_targets))
	action = SwitchIn(switch_count)
	context.ActionID = &action.ID
	transaction := MakeTransaction(action, context)
	transaction.Ready = false
	if !g.HasPlayerPrompt(player.ID) {
		g.AddPrompt(transaction)
	}

	return false
}

func validatePlayer(g *Game, player Player, valid *bool) {
	missing_pos := make([]uuid.UUID, 0)
	for _, pos := range player.Positions {
		if pos.ActorID == nil {
			missing_pos = append(missing_pos, pos.ID)
			continue
		}

		actor, ok := g.GetActorByID(*pos.ActorID)
		if !ok {
			missing_pos = append(missing_pos, pos.ID)
			continue
		}

		if actor.Summon != nil {
			if !actor.Summon.Alive {
				g.UpdateActor(actor.ID, func(a Actor) Actor {
					a.SetSummon(nil)
					return a
				})
			}
		}

		if !actor.Alive {
			missing_pos = append(missing_pos, pos.ID)
			context := NewContext().WithTargetIDs([]uuid.UUID{actor.ID})
			transaction := MakeTransaction(RemovePositions, context)
			g.JumpTransaction(transaction)
		}
	}

	if len(missing_pos) > 0 {
		*valid = queueSwitchPrompt(g, player, missing_pos)
	}
}

func (g *Game) Validate() bool {
	valid := true

	for _, player := range g.Players {
		validatePlayer(g, player, &valid)
	}

	return valid
}
