package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var ShadowClone = MakeShadowClone()

func MakeShadowClone() game.Action {
	config := game.ActionConfig{
		Name:        "Shadow Clone",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Ninjutsu,
		Description: "Summons a shadow clone to take damage. Pay 1/4th of Max HP.",
		TargetType:  game.TargetActorID,
	}

	return game.Action{
		ID:              uuid.MustParse("e2a1768a-fb9a-5891-a703-a20cf8bcbd6e"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.ComposeGF(game.SourceIsAlive, game.SourceHasHpRatio(0.25)),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}
				s, ok := g.GetSource(context)
				if !ok {
					return transactions
				}
				source := s.Resolve(g)
				damage := game.RatioDamage(0.25)
				damage_context := game.MakeContextForActor(s)

				mut := game.GameMutation{
					Delta: func(mp, mg game.Game, mc game.Context) game.Game {
						mg.UpdateActor(*mc.SourceActorID, func(a game.Actor) game.Actor {
							hp := game.Round(float64(source.Stats[game.StatHP]) * 0.25)
							summon := game.MakeActor(
								game.ActorDef{
									ActorID:   uuid.New(),
									SpriteURL: "/sprites/sub_64.png",
									Name:      source.Name,
									Stats: map[game.ActorStat]int{
										game.StatHP: hp,
									},
								},
								a.PlayerID,
								a.Experience,
								nil,
								nil,
								[]game.Action{},
								game.FocusNone,
								map[game.ActorStat]int{},
							)
							a.SetSummonFromActor(&summon, true)
							return a
						})
						return mg
					},
				}

				transactions = append(
					transactions,
					game.MakeTransaction(damage, damage_context),
					game.MakeTransaction(mut, context),
				)

				return transactions
			},
		},
	}
}
