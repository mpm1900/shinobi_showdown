package mutations

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

func RedirectSingleTargetEnemyActions(source game.Actor) game.GameMutation {
	return game.GameMutation{
		Delta: func(p game.Game, g game.Game, context game.Context) game.Game {
			for i, a := range g.Actions {
				if a.Mutation.Config.IgnoreRedirect {
					continue
				}

				if a.Context.SourcePlayerID == nil || a.Context.SourceActorID == nil {
					continue
				}

				targets := g.GetTargets(a.Context)
				if len(targets) != 1 {
					continue
				}

				if *a.Context.SourcePlayerID != source.PlayerID {
					c_source, ok := g.GetSource(a.Context)
					if !ok {
						continue
					}
					rc_source := c_source.Resolve(g)
					if rc_source.IgnoreRedirect {
						continue
					}

					if a.Mutation.TargetType == game.TargetActorID {
						g.Actions[i].Context.TargetActorIDs = []uuid.UUID{source.ID}
					}
					if a.Mutation.TargetType == game.TargetPositionID && source.IsActive() {
						g.Actions[i].Context.TargetPositionIDs = []uuid.UUID{*source.PositionID}
					}
				}
			}

			return g
		},
	}
}
