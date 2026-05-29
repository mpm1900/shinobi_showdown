package mutations

import (
	"math/rand/v2"
	"shinobi_showdown/internal/game"
)

var Sleep = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		for _, target := range g.GetTargets(context) {
			resolved := target.Resolve(g)
			if resolved.Statused {
				continue
			}

			g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
				if a.Statused {
					return a
				}

				a.SleepCounter = rand.IntN(2) + 1
				a.Sleeping = true
				a.Statused = true

				return a
			})
			g.PushLog(game.NewLogContext("$source$ went to sleep.", game.MakeContextForActor(target)))
		}

		return g
	},
}

var Burn = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		for _, target := range g.GetTargets(context) {
			resolved := target.Resolve(g)
			if resolved.Statused {
				continue
			}

			g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
				if a.Statused {
					return a
				}

				a.Burned = true
				a.Statused = true

				return a
			})
			g.PushLog(game.NewLogContext("$source$ became burned.", game.MakeContextForActor(target)))
		}

		return g
	},
}

var Paralyze = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		targets := g.GetTargets(context)
		for _, target := range targets {
			resolved := target.Resolve(g)
			if resolved.Statused {
				continue
			}

			g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
				if a.Statused {
					return a
				}

				a.Paralyzed = true
				a.Statused = true

				return a
			})
			g.PushLog(game.NewLogContext("$source$ became paralyzed.", game.MakeContextForActor(target)))
		}

		return g
	},
}

var Poison = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		targets := g.GetTargets(context)
		for _, target := range targets {
			resolved := target.Resolve(g)
			if resolved.Statused {
				continue
			}

			g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
				if a.Statused {
					return a
				}

				a.Poisoned = true
				a.PoisonedCounter = 1
				a.Statused = true

				return a
			})
			g.PushLog(game.NewLogContext("$source$ became poisoned.", game.MakeContextForActor(target)))
		}

		return g
	},
}

var Revive = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		source, ok := g.GetSource(context)
		if !ok {
			return g
		}

		resolved := source.Resolve(g)

		g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
			a.Alive = true
			a.Damage = resolved.Stats[game.StatHP] / 2
			return a
		})

		return g
	},
}

var RemoveItem = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		source, ok := g.GetSource(context)
		if !ok {
			return g
		}

		g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
			a.Item = nil
			return a
		})

		return g
	},
}
var ExchangeItems = game.GameMutation{
	Delta: func(p, g game.Game, context game.Context) game.Game {
		source, ok := g.GetSource(context)
		if !ok {
			return g
		}

		targets := g.GetTargets(context)
		if len(targets) == 0 {
			return g
		}

		var item = game.Modifier{}
		for _, target := range targets {
			g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
				item = *a.Item
				a.Item = source.Item
				return a
			})
		}

		g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
			a.Item = &item
			return a
		})

		return g
	},
}

func Transform(def game.ActorDef) game.GameMutation {
	return game.GameMutation{
		Delta: func(p, g game.Game, context game.Context) game.Game {
			source, ok := g.GetSource(context)
			if !ok {
				return g
			}

			g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
				a.Transform(def)
				return a
			})

			tctx := game.MakeContextForActor(source)
			g.On(game.OnActorTransform, &tctx)

			return g
		},
	}
}

func KillSource() game.GameMutation {
	return game.GameMutation{
		Delta: func(p, g game.Game, context game.Context) game.Game {
			source, ok := g.GetSource(context)
			if !ok {
				return g
			}

			resolved := source.Resolve(g)
			g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
				a.Damage = resolved.Stats[game.StatHP]
				a.Alive = false
				return a
			})

			return g
		},
	}
}
