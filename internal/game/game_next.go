package game

import (
	"time"
)

func (g *Game) NextPhase() {
	g.ActiveTransaction = nil
	switch g.Turn.Phase {
	case TurnStart:
		g.Turn.Phase = TurnMain
	case TurnInit, TurnMain:
		g.Turn.Phase = TurnEnd
	case TurnEnd:
		g.Turn.Phase = TurnCleanup
	case TurnCleanup:
		// Keep cleanup stable so callers can run end-of-turn bookkeeping once
		// without immediately wrapping back to main in the same loop tick.
	}
}

func (g *Game) NextTransaction() bool {
	transaction, err := g.Transactions.Dequeue()
	if err != nil {
		return false
	}

	n, ok := ResolveTransaction(*g, *g, transaction, *g)
	if !ok {
		return false
	}

	*g = n
	return true
}

func (g *Game) PreAction() bool {
	g.SortActions()
	transaction, err := g.Actions.Peek()
	if err != nil {
		return false
	}

	source, ok := g.GetSource(transaction.Context)
	if !ok || !source.Alive {
		return false
	}

	if !g.Turn.PreAction {
		g.Turn.PreAction = true
		g.ActiveTransaction = MakeGameActiveTransaction(transaction)
		count := g.On(OnActionStart, &transaction.Context)
		return count > 0
	}

	return false
}

func (g *Game) PostAction(context Context) {
	g.Turn.PreAction = false
	g.On(OnActionEnd, &context)
}

func (g *Game) NextAction() bool {
	g.SortActions()
	transaction, err := g.Actions.Dequeue()
	if err != nil {
		g.ActiveTransaction = nil
		return false
	}

	source, ok := g.GetSource(transaction.Context)
	if !ok || !source.Alive {
		g.ActiveTransaction = nil
		return true
	}

	resolved := source.Resolve(*g)
	action, ok := resolved.GetActionByID(transaction.Mutation.ID)
	if ok {
		queuedConfig := transaction.Mutation.Config
		action.Config = queuedConfig
		transaction.Mutation = action
	}

	if transaction.Mutation.MapContext != nil {
		c := transaction.Mutation.MapContext(*g, transaction.Context)
		transaction.Context = c
	}

	g.ActiveTransaction = MakeGameActiveTransaction(transaction)
	g.RunAction(transaction)
	g.PostAction(transaction.Context)
	return true
}
func (g *Game) NextPrompt() bool {
	transaction, err := g.Prompts.Dequeue()
	g.ActiveTransaction = MakeGameActiveTransaction(transaction)
	if err != nil {
		return false
	}

	g.RunPrompt(transaction)
	return true
}

func (g *Game) NextTrigger() bool {
	g.SortTriggers()
	transaction, err := g.Triggers.Dequeue()
	a_tx := Transaction[Action]{}
	a_tx.Context = transaction.Context
	g.ActiveTransaction = MakeGameActiveTransaction(a_tx)
	if err != nil {
		return false
	}

	g.RunTrigger(transaction)
	return true
}

func (g Game) GetBaseTick() time.Duration {
	if g.Turn.PreAction {
		return 0
	}
	return time.Second / 2
}

func (g *Game) Next() bool {
	g.Tick = g.GetBaseTick()
	if g.NextTransaction() {
		return true
	}

	g.Tick = g.GetBaseTick()
	if g.AllPromptsReady() {
		if g.NextPrompt() {
			return true
		}
	}

	g.Tick = g.GetBaseTick()
	if g.NextTrigger() {
		return true
	}

	g.Tick = g.GetBaseTick()
	if g.AllPromptsReady() {
		if g.NextPrompt() {
			return true
		}
	}

	if !g.Validate() {
		return false
	}

	if g.PreAction() {
		return true
	}
	g.Tick = time.Second * 2
	if g.NextAction() {
		return true
	}

	g.Tick = g.GetBaseTick()
	if !g.Validate() {
		return false
	}
	g.NextPhase()

	return false
}
