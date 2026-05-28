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
	Check: Match__True,
	ActionMutation: ActionMutation{
		Delta: func(parent Game, input Game, context Context) []Transaction[GameMutation] {
			transactions := NewTransactionBuilder()

			mut := GameMutation{
				Delta: func(p Game, g Game, c Context) Game {
					t := g.Turn.Count
					for i := range g.Actors {
						if t > 0 {
							g.Actors[i].DecrementCooldowns()
							g.Actors[i].RecoverStamina(g, STAMINA_RECOVERY)
							if g.Actors[i].Poisoned {
								g.Actors[i].PoisonedCounter++
							}
						}
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

			transactions.PushOne(MakeTransaction(mut, context))

			return transactions.Build()
		},
	},
}
