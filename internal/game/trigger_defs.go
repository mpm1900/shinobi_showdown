package game

import (
	"github.com/google/uuid"
)

func RunTriggerTx(on TriggerOn, context Context) GameTransaction {
	return GameTransaction{
		Context: context,
		Mutation: GameMutation{
			Delta: func(p, g Game, context Context) Game {
				g.On(on, &context)
				return g
			},
		},
	}
}

var END_OF_TURN_TRIGGER Trigger = Trigger{
	ID:    uuid.MustParse("f63aefeb-02cf-4dbd-93f9-8f1908f99d4f"),
	On:    OnTurnEnd,
	Check: func(p, g Game, context Context, tx Transaction[Modifier]) bool { return true },
	ActionMutation: ActionMutation{
		Delta: func(parent Game, input Game, context Context) []Transaction[GameMutation] {
			var transactions []Transaction[GameMutation]
			mut := GameMutation{
				Delta: func(p Game, g Game, c Context) Game {
					t := g.Turn.Count
					for i := range g.Actors {
						if t > 0 {
							g.Actors[i].DecrementCooldowns()
							g.Actors[i].RecoverStamina(g, 0.08)
							if g.Actors[i].Poisoned {
								g.Actors[i].PoisonedCounter++
							}
						}
						g.Actors[i].IncrementTurns()
					}

					if t > 0 {
						g.FilterModifiers(func(mod Transaction[Modifier]) bool {
							return mod.Mutation.Duration != 0
						})

						for i, _ := range g.Modifiers {
							g.Modifiers[i].Mutation.DecrementTimers()
						}
					}

					return g
				},
			}

			tx := MakeTransaction(mut, context)
			transactions = append(transactions, tx)

			return transactions
		},
	},
}
