package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

func FilterTerrain() game.GameTransaction {
	filter := game.NewGameMutation()
	filter.Delta = func(p, g game.Game, context game.Context) game.Game {
		g.FilterModifiers(func(mod game.Transaction[game.Modifier]) bool {
			return !mod.Mutation.Terrain
		})
		return g
	}
	return game.MakeTransaction(filter, game.NewContext())
}

func SetTerrain(gid uuid.UUID, terrain game.GameTerrain) game.Modifier {
	return game.Modifier{
		ID:       gid,
		GroupID:  &gid,
		Show:     true,
		Terrain:  true,
		Duration: 4,
		GameStateMutations: []game.GameStateMutation{
			game.MakeGameStateMutation(
				&gid,
				game.MutPriorityGameState0,
				game.GS_TrueFilter,
				func(g game.Game, gs game.GameState, context game.Context) game.GameState {
					gs.Terrain = terrain
					return gs
				},
			),
		},
		ActorMutations: []game.ActorMutation{},
		Triggers:       []game.Trigger{},
	}
}

var floodedTerrainID = uuid.MustParse("f1784ee6-ba9a-4672-8eb5-e44b62021fec")

func FloodedTerrain() game.Modifier {
	mod := SetTerrain(floodedTerrainID, game.GameTerrainFlooded)
	mod.Name = "Flooded Terrain"
	mod.Icon = "flooded"
	mod.Description = "no innate effect."
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&floodedTerrainID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.GameHasTerrain(game.GameTerrainFlooded)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Terrain != game.GameTerrainFlooded {
					return actor
				}

				// TODO
				return actor
			},
		),
	}
	return mod
}

var electrifiedTerrainID = uuid.MustParse("b8b348b6-fd35-4eca-a0fa-790994d6f205")

func ElectrifiedTerrain() game.Modifier {
	mod := SetTerrain(electrifiedTerrainID, game.GameTerrainElectrified)
	mod.Name = "Electrified Terrain"
	mod.Icon = "electrified"
	mod.Description = "no innate effect."
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&electrifiedTerrainID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.GameHasTerrain(game.GameTerrainElectrified)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Terrain != game.GameTerrainElectrified {
					return actor
				}

				// TODO
				return actor
			},
		),
	}
	return mod
}

var flamableTerrainID = uuid.MustParse("3cb62ea5-6657-464c-84b8-684721b6bfdb")

func FlamableTerrain() game.Modifier {
	mod := SetTerrain(flamableTerrainID, game.GameTerrainFlamable)
	mod.Name = "Flammable Terrain"
	mod.Icon = "flamable"
	mod.Description = ""
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&flamableTerrainID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.GameHasTerrain(game.GameTerrainFlamable)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Terrain != game.GameTerrainFlamable {
					return actor
				}

				// TODO
				return actor
			},
		),
	}
	return mod
}

var rockyTerrainID = uuid.MustParse("24aa5833-2579-4926-ae43-7bc1e891a210")

func RockyTerrain() game.Modifier {
	mod := SetTerrain(rockyTerrainID, game.GameTerrainRocky)
	mod.Name = "Rocky Terrain"
	mod.Icon = "rocky_terrain"
	mod.Description = "Non-earth, grounded shinobi have decreased accuracy."
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&rockyTerrainID,
			game.MutPriorityPreStagedStats,
			game.ComposeAF(game.ActiveFilter, game.GameHasTerrain(game.GameTerrainRocky)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				_, is_earth := actor.Natures[game.NsEarth]
				if state.Terrain != game.GameTerrainRocky || is_earth {
					return actor
				}

				actor.Stages[game.StatAccuracy] = actor.Stages[game.StatAccuracy] - 1
				return actor
			},
		),
	}
	return mod
}

var chakraTerrainID = uuid.MustParse("37041a49-3f69-481c-a807-357085fa91ef")

func ChakraTerrain() game.Modifier {
	mod := SetTerrain(chakraTerrainID, game.GameTerrainChakra)
	mod.Name = "Chakra Terrain"
	mod.Icon = "chakra_terrain"
	mod.Description = "On turn end: all shinobi heal 1/16th HP."
	mod.ActorMutations = []game.ActorMutation{
		game.MakeActorMutation(
			&chakraTerrainID,
			game.MutPriorityPreStagedStats,
			game.ComposeAF(game.ActiveFilter, game.GameHasTerrain(game.GameTerrainChakra)),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Terrain != game.GameTerrainChakra {
					return actor
				}

				return actor
			},
		),
	}
	mod.Triggers = []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: chakraTerrainID,
			On:         game.OnTurnEnd,
			Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
				return g.HasTerrain(game.GameTerrainChakra, context)
			},
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.HasWeather(game.GameWeatherSandstorm),
				Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					mut := game.RatioHeal(0.0625)
					mut_ctx := context
					mut_ctx.TargetActorIDs = []uuid.UUID{}
					for _, target := range g.GetActiveActors() {
						mut_ctx.TargetActorIDs = append(mut_ctx.TargetActorIDs, target.ID)
					}
					return []game.Transaction[game.GameMutation]{
						game.MakeTransaction(mut, mut_ctx),
					}
				},
			},
		},
	}
	return mod
}

func ApplyTerrain(g game.Game, context game.Context, terrain game.GameTerrain, mod game.Modifier) []game.GameTransaction {
	transactions := []game.GameTransaction{}

	state, _ := g.GetState(context)
	if state.Terrain == terrain {
		return transactions
	}

	filter := FilterTerrain()
	transactions = append(transactions, filter)

	mut := mutations.AddModifiers(false, mod)
	transaction := game.MakeTransaction(mut, game.NewContext())
	transactions = append(transactions, transaction)

	return transactions
}
