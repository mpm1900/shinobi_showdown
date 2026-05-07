package mutations

import (
	"fmt"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

func CheckGameJutsuImmunity(g *game.Game, source game.Actor) bool {
	config, ok := game.GetActiveActionConfig(*g, game.ActionConfig{})
	if ok {
		if CheckJutsuImmunity(config, source) {
			log_ctx := game.MakeContextForActor(source)
			g.PushLog(game.MakeGameLog(fmt.Sprintf("$source$ was immune to %s.", config.Jutsu), log_ctx, 1))
			return true
		}
	}

	return false
}
func CheckJutsuImmunity(config game.ActionConfig, source game.Actor) bool {
	return source.HasJutsuImmunity(config.Jutsu)
}
func CheckImmunity(id uuid.UUID, source game.Actor) bool {
	return source.HasImmunity(id)
}

func AddStatus(checkWarded bool, modifiers ...game.Modifier) game.GameMutation {
	mut := AddModifiers(checkWarded, modifiers...)
	baseDelta := mut.Delta
	mut.Delta = func(p, g game.Game, context game.Context) game.Game {
		source, ok := g.GetSource(context)
		if !ok {
			return g
		}

		if CheckGameJutsuImmunity(&g, source) {
			return g
		}

		resolved := source.Resolve(g)
		if resolved.Statused {
			return g
		}

		return baseDelta(p, g, context)
	}
	return mut
}

func AddModifiers(checkWarded bool, modifiers ...game.Modifier) game.GameMutation {
	return game.GameMutation{
		Delta: func(p game.Game, g game.Game, context game.Context) game.Game {
			for _, modifier := range modifiers {
				mod_tx := game.MakeTransaction(modifier, context)
				mod_tx.Context.ModifierID = &mod_tx.ID

				// logs
				if context.SourceActorID == nil && len(modifier.GameStateMutations) > 0 {
					g.PushLog(game.MakeGameLog(fmt.Sprintf("The battlefield gained %s.", mod_tx.Mutation.Name), game.NewContext(), 1))
				}
				hasCandidate := false
				hasApplicableTarget := false
				for _, actor := range g.GetActionableActors() {
					if !game.CheckModifierForActor(mod_tx, g, actor) {
						continue
					}

					hasCandidate = true
					resolved := actor.Resolve(g)

					/**
					 * Filtering out via immune, safeguarded, and warded check
					 */
					if CheckGameJutsuImmunity(&g, resolved.Actor) {
						log_ctx := game.MakeContextForActor(resolved.Actor)
						g.PushLog(game.MakeGameLog("$source$ was immune.", log_ctx, 1))
						continue
					}
					if (modifier.GroupID != nil && resolved.HasImmunity(*modifier.GroupID)) || resolved.HasImmunity(modifier.ID) {
						mod_tx.Context.FilterOutTarget(actor)

						log_ctx := game.MakeContextForActor(resolved.Actor)
						g.PushLog(game.MakeGameLog(fmt.Sprintf("$source$ was immune to %s.", modifier.Name), log_ctx, 1))
						continue
					}
					if context.ModifierID != nil {
						parent_mod, ok := g.GetModifierByID(*context.ModifierID)
						if ok && resolved.HasImmunity(*context.ModifierID) {
							mod_tx.Context.FilterOutTarget(actor)

							log_ctx := game.MakeContextForActor(resolved.Actor)
							g.PushLog(game.MakeGameLog(fmt.Sprintf("$source$ was immune to %s.", parent_mod.Name), log_ctx, 1))
							continue
						}
					}

					if resolved.Safeguarded && context.SourcePlayerID != nil && resolved.PlayerID != *context.SourcePlayerID {
						mod_tx.Context.FilterOutTarget(actor)

						context.SourceActorID = &actor.ID
						g.PushLog(game.MakeGameLog("$source$ was safeguarded.", context.WithSource(actor.ID), 1))
						continue
					}
					if checkWarded && resolved.Warded {
						mod_tx.Context.FilterOutTarget(actor)

						context.SourceActorID = &actor.ID
						g.PushLog(game.MakeGameLog("$source$ was warded.", context.WithSource(actor.ID), 1))
						continue
					}

					hasApplicableTarget = true
					if hasCandidate {
						g.PushLog(game.MakeGameLog(fmt.Sprintf("$source$ gained %s.", modifier.Name), context.WithSource(actor.ID), 1))
					}
				}

				if hasCandidate && !hasApplicableTarget {
					continue
				}

				g.Modifiers = append(g.Modifiers, mod_tx)
				g.On(game.OnModifierAdd, &mod_tx.Context)
			}

			return g
		},
	}
}

func RemoveModifierTxByID(ID uuid.UUID) game.GameMutation {
	return game.GameMutation{
		Delta: func(p game.Game, g game.Game, context game.Context) game.Game {
			modifiers := []game.Transaction[game.Modifier]{}
			for _, tx := range g.Modifiers {
				if tx.ID != ID {
					modifiers = append(modifiers, tx)
				}
			}

			g.Modifiers = modifiers
			return g
		},
	}
}

func RemoveModifierWhere(where func(game.Transaction[game.Modifier]) bool) game.GameMutation {
	return game.GameMutation{
		Delta: func(p, g game.Game, context game.Context) game.Game {
			modifiers := []game.Transaction[game.Modifier]{}
			for _, tx := range g.Modifiers {
				if !where(tx) {
					modifiers = append(modifiers, tx)
				}
			}

			g.Modifiers = modifiers
			return g
		},
	}
}

var ConsumeItem game.GameMutation = game.GameMutation{
	Delta: func(p game.Game, g game.Game, context game.Context) game.Game {
		if context.SourceActorID == nil {
			return g
		}

		g.UpdateActor(*context.SourceActorID, func(a game.Actor) game.Actor {
			a.Item = nil
			return a
		})

		g.On(game.OnItemConsume, &context)

		return g
	},
}
