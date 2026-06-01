package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

func FilterWeather() game.GameTransaction {
	filter := game.NewGameMutation()
	filter.Delta = func(p, g game.Game, context game.Context) game.Game {
		g.FilterModifiers(func(mod game.Transaction[game.Modifier]) bool {
			return !mod.Mutation.Weather
		})
		return g
	}
	return game.MakeTransaction(filter, game.NewContext())
}

func ClearWeather() []game.GameTransaction {
	var transactions = []game.GameTransaction{}
	transactions = append(transactions, FilterWeather())

	return transactions
}

func SetWeather(gid uuid.UUID, weather game.GameWeather) game.Modifier {
	return game.Modifier{
		ID:       gid,
		GroupID:  &gid,
		Show:     true,
		Weather:  true,
		Duration: game.ModifierDurationInf,
		GameStateMutations: []game.GameStateMutation{
			game.MakeGameStateMutation(
				&gid,
				game.MutPriorityGameState0,
				game.GS_TrueFilter,
				func(g game.Game, gs game.GameState, context game.Context) game.GameState {
					gs.Weather = weather
					return gs
				},
			),
		},
		ActorMutations: []game.ActorMutation{},
		Triggers:       []game.Trigger{},
	}
}

var rainWeatherID = uuid.MustParse("b28933cd-ad6a-5b83-acff-4dd084dad6e5")

func RainWeather() game.Modifier {
	mod := SetWeather(rainWeatherID, game.GameWeatherRain)
	mod.Name = "Heavy Rain"
	mod.Icon = "raining"
	mod.Description = "Fire damage reduced by 1.5x, Water damage increased by 1.5x."
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&rainWeatherID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.GameHasWeather(game.GameWeatherRain)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Weather != game.GameWeatherRain {
					return actor
				}

				actor.NatureDamage[game.NatureFire] /= 1.5
				actor.NatureDamage[game.NatureWater] *= 1.5
				return actor
			},
		),
	}
	return mod
}

var sandstormWeatherID = uuid.MustParse("6dfb1ee5-a0dc-4ee0-8a52-2c062e8d8b3e")

func SandstormWeather() game.Modifier {
	mod := SetWeather(sandstormWeatherID, game.GameWeatherSandstorm)
	mod.Name = "Sandstorm"
	mod.Description = "On turn end: lose 1/6th HP. Earth shinobi are immune."
	mod.Icon = "sandstorm"
	mod.ActorMutations = []game.ActorMutation{
		game.NewNoopActive(&sandstormWeatherID),
	}
	mod.Triggers = append(mod.Triggers, game.Trigger{
		ID:         uuid.New(),
		ModifierID: sandstormWeatherID,
		On:         game.OnTurnEnd,
		Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
			return g.HasWeather(game.GameWeatherSandstorm, context)
		},
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.HasWeather(game.GameWeatherSandstorm),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction{
				mut := game.RatioDamage(0.0625)
				mut_ctx := context
				mut_ctx.TargetActorIDs = []uuid.UUID{}
				for _, target := range g.GetActiveActors() {
					_, ok := target.Natures[game.NsEarth]
					if ok {
						continue
					}
					mut_ctx.TargetActorIDs = append(mut_ctx.TargetActorIDs, target.ID)
				}
				return []game.GameTransaction{
					game.MakeTransaction(mut, mut_ctx),
				}
			},
		},
	})
	return mod
}

var sunlightWeatherID = uuid.MustParse("1d911f5f-1cfd-40b2-b44f-e365fc500da6")

func SunlightWeather() game.Modifier {
	mod := SetWeather(sunlightWeatherID, game.GameWeatherSunlight)
	mod.Name = "Heavy Rain"
	mod.Icon = "raining"
	mod.Description = "Water damage reduced by 1.5x, Fire damage increased by 1.5x."
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&sunlightWeatherID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.GameHasWeather(game.GameWeatherSunlight)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Weather != game.GameWeatherSunlight {
					return actor
				}

				actor.NatureDamage[game.NatureWater] /= 1.5
				actor.NatureDamage[game.NatureFire] *= 1.5
				return actor
			},
		),
	}
	return mod
}
