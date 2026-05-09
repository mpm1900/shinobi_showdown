package game

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type GameTransaction = Transaction[GameMutation]

type GameStatus string

type GameActiveTransaction struct {
	Transaction[Action]
	Message string `json:"message"`
}

const (
	GameStatusRunning GameStatus = "running"
	GameStatusIdle    GameStatus = "idle"
	GameStatusWaiting GameStatus = "waiting"
)

/**
 * Game is the main state container state for a game instance
 */
type Game struct {
	Status            GameStatus             `json:"status"`
	Turn              Turn                   `json:"turn"`
	ActiveTransaction *GameActiveTransaction `json:"active_transaction"`
	Players           []Player               `json:"players"`

	state     GameState
	Actors    []Actor                 `json:"actors"`
	Modifiers []Transaction[Modifier] `json:"modifiers"`

	Tick time.Duration `json:"-"`

	/*
	 * [Transactions] are the quable/storable change in state
	 * 	- mostly to animate and modify state changes while running
	 *  - so one aciton can do multiple things, in order
	 */
	Transactions Queue[GameTransaction] `json:"transactions"`
	/*
	 * [Triggers] are basically Actions that do not need input to resolve
	 *  - triggers are to delay, repeat, or statically/conditionally respond to events
	 */
	Triggers Queue[Transaction[Trigger]] `json:"triggers"`
	/*
	 * [Actions] are a sorted list of actions to be resolved
	 *  - actions resolve down to reansactions
	 *
	 * [Prompts] are special actions that pause Action resolution loops.
	 *  - can stop actions from resolving while waiting for user input
	 *  - switch ins will be a common Prompt
	 */
	Actions Queue[Transaction[Action]] `json:"actions"`
	Prompts Queue[Transaction[Action]]

	/**
	 * [QueuedActions] are a map of actor IDs to action IDs
	 *  - the actions need not be on the user, it does a global lookup.
	 */
	QueuedActions map[uuid.UUID]Transaction[uuid.UUID] `json:"queued_actions"`

	Log []GameLog `json:"log"`

	/**
	 * [ActionRegistry]
	 */
	ActionRegistry map[uuid.UUID]Action

	disableActionFilterEval bool `json:"-"`
}

func MakeGameActiveTransaction(tx Transaction[Action]) *GameActiveTransaction {
	return &GameActiveTransaction{
		Transaction: tx,
		Message:     "",
	}
}

func NewGame(actionRegistry map[uuid.UUID]Action) Game {
	return Game{
		Status: GameStatusIdle,
		Turn: Turn{
			Count: 0,
			Phase: TurnInit,
		},
		state:             NewGameState(),
		ActiveTransaction: nil,
		Players:           make([]Player, 0),
		Actors:            make([]Actor, 0),
		Modifiers:         make([]Transaction[Modifier], 0),
		Tick:              time.Second / 2,
		Transactions:      MakeQueue[GameTransaction](),
		Actions:           MakeQueue[Transaction[Action]](),
		Prompts:           MakeQueue[Transaction[Action]](),
		Triggers:          MakeQueue[Transaction[Trigger]](),
		QueuedActions:     make(map[uuid.UUID]Transaction[uuid.UUID]),
		Log:               []GameLog{},

		ActionRegistry: actionRegistry,
	}
}

func (g *Game) Reset() {
	g.state = NewGameState()
	g.ActiveTransaction = nil
	g.Turn.Count = 0
	g.Turn.Phase = TurnInit
	g.Log = []GameLog{}
	g.Modifiers = []Transaction[Modifier]{}
	g.Actions = MakeQueue[Transaction[Action]]()
	g.QueuedActions = make(map[uuid.UUID]Transaction[uuid.UUID])

	for i, _ := range g.Players {
		g.Players[i].Reset()
	}

	for i, _ := range g.Actors {
		g.Actors[i].Reset()
	}
}
func (g *Game) ResetPlayer(playerID uuid.UUID) {
	for i, _ := range g.Players {
		if g.Players[i].ID != playerID {
			continue
		}

		g.Players[i].Reset()
	}

	for i, _ := range g.Actors {
		if g.Actors[i].PlayerID != playerID {
			continue
		}

		g.Actors[i].Reset()
	}
}

/**
 * Getters
 */

// Player
func (g Game) GetPlayer(predicate func(Player) bool) (Player, bool) {
	for _, a := range g.Players {
		if predicate(a) {
			return a, true
		}
	}

	return Player{}, false
}
func (g Game) GetPlayerByID(ID uuid.UUID) (Player, bool) {
	return g.GetPlayer(func(p Player) bool {
		return p.ID == ID
	})
}

// Actor
func (g Game) GetActor(predicate func(Actor) bool) (Actor, bool) {
	for _, a := range g.Actors {
		if predicate(a) {
			return a, true
		}
	}

	return Actor{}, false
}
func (g Game) GetActorByID(ID uuid.UUID) (Actor, bool) {
	return g.GetActor(func(a Actor) bool {
		return a.ID == ID
	})
}
func (g Game) GetActors(predicate func(Actor) bool) []Actor {
	actors := make([]Actor, 0)
	for _, a := range g.Actors {
		if predicate(a) {
			actors = append(actors, a)
		}
	}
	return actors
}
func (g Game) GetActorsByPlayer(playerID uuid.UUID) []Actor {
	return g.GetActors(func(a Actor) bool {
		return a.PlayerID == playerID
	})
}
func (g Game) GetActiveActors() []Actor {
	return g.GetActors(func(a Actor) bool {
		active := a.IsActive()
		return active && a.Alive
	})
}
func (g Game) GetActionableActors() []Actor {
	return g.GetActors(func(a Actor) bool {
		active := a.IsActive()
		return active && a.Alive && !a.Stunned
	})
}
func (g Game) GetActorsFilters(context Context, filters ...ActorFilter) []Actor {
	filter := ComposeAF(filters...)
	return g.GetActors(func(a Actor) bool {
		return filter(g, a, context)
	})
}

func (g Game) GetResolvedActors() map[uuid.UUID]ResolvedActor {
	actors := make(map[uuid.UUID]ResolvedActor, len(g.Actors))
	for _, a := range g.Actors {
		resolvedActor := a.Resolve(g)
		actors[a.ID] = resolvedActor
	}

	return actors
}

func (g Game) GetTriggers(on TriggerOn, context *Context) []Transaction[Trigger] {
	ctx := NewContext()
	if context != nil {
		ctx = *context
	}
	triggers := []Transaction[Trigger]{
		MakeTransaction(END_OF_TURN_TRIGGER, ctx),
	}

	transactions := g.GetModifiers()
	modifiers := make([]Transaction[Modifier], 0, len(transactions))
	modifiers = append(modifiers, transactions...)
	modifiers = append(modifiers, GetActorModifiers(g)...)

	/**
	 * This is so that leaving actos are temporarily considered active
	 * for trigger resolution
	 */
	if on == OnActorLeave && context != nil && context.SourceActorID != nil {
		source, ok := g.GetActorByID(*context.SourceActorID)
		if ok {
			source_ctx := newActorContext(source)
			for _, mod := range source.GetModifiers() {
				modifiers = append(modifiers, MakeTransaction(mod, source_ctx))
			}
		}
	}

	for _, tx := range modifiers {
		if context == nil {
			ctx = tx.Context
		}
		for _, trig := range tx.Mutation.Triggers {
			if trig.On != on {
				continue
			}

			if trig.Check != nil && !trig.Check(g, g, ctx, tx) {
				fmt.Println(tx.Mutation.Name, "FAILED")
				continue
			}
			triggers = append(triggers, MakeTransaction(trig, ctx))
		}
	}

	return triggers
}
func (g Game) GetModifiers() []Transaction[Modifier] {
	var transactions []Transaction[Modifier] = []Transaction[Modifier]{}
	for _, tx := range g.Modifiers {
		if tx.Mutation.Delay <= 0 {
			transactions = append(transactions, tx)
		}
	}

	return transactions
}
func (g Game) GetModifierTxByID(txID uuid.UUID) (Transaction[Modifier], bool) {
	modifiers := make([]Transaction[Modifier], 0, len(g.Modifiers))
	modifiers = append(modifiers, g.Modifiers...)
	modifiers = append(modifiers, GetActorModifiers(g)...)

	for _, m := range modifiers {
		if m.ID == txID {
			return m, m.Mutation.Delay <= 0
		}
	}
	return Transaction[Modifier]{}, false
}
func (g Game) GetModifierByTxID(txID uuid.UUID) (Modifier, bool) {
	modifiers := g.GetModifiers()
	modifiers = append(modifiers, GetActorModifiers(g)...)

	for _, m := range modifiers {
		if m.ID == txID {
			return m.Mutation, m.Mutation.Delay <= 0
		}
	}
	return Modifier{}, false
}
func (g Game) GetModifierByID(ID uuid.UUID) (Modifier, bool) {
	modifiers := g.GetModifiers()
	modifiers = append(modifiers, GetActorModifiers(g)...)

	for _, m := range modifiers {
		if m.Mutation.ID == ID || (m.Mutation.GroupID != nil && *m.Mutation.GroupID == ID) {
			return m.Mutation, m.Mutation.Delay <= 0
		}
	}
	return Modifier{}, false
}
func (g Game) GetPromptTxByID(ID uuid.UUID) (Transaction[Action], bool) {
	for _, tx := range g.Prompts {
		if tx.ID == ID {
			return tx, true
		}
	}

	return Transaction[Action]{}, false
}
func (g Game) GetQueuedAction(context Context) (Action, bool) {
	for _, action := range g.Actions {
		if action.Mutation.ID == *context.ActionID && action.Context.SourceActorID == context.SourceActorID {
			return action.Mutation, true
		}
	}

	return Action{}, false
}
func (g Game) FindQueuedAction(predicate func(Transaction[Action]) bool) (Transaction[Action], bool) {
	for _, action := range g.Actions {
		if predicate(action) {
			return action, true
		}
	}

	return Transaction[Action]{}, false
}
func (g Game) GetActiveAction() (Action, bool) {
	if g.ActiveTransaction != nil {
		return g.ActiveTransaction.Mutation, true
	}

	return Action{}, false
}
func (g Game) GetState(context Context) (GameState, []uuid.UUID) {
	applied := make([]uuid.UUID, 0)
	mutations, transactions := GetAllGameStateMutations(g)
	state := g.state
	for _, mut := range mutations {
		mutContext := ResolveModifierTransactionContext(context, transactions, mut.TransactionID)
		tx := MakeTransaction(mut.Mutation, mutContext)
		s, ok := ResolveTransaction(g, state, tx, state)
		if ok {
			if mut.TransactionID != nil {
				applied = append(applied, *mut.TransactionID)
			}
			state = s
		}
	}
	return state, applied
}
func (g Game) HasWeather(weather GameWeather, context Context) bool {
	state, _ := g.GetState(context)
	return state.Weather == weather
}
func (g Game) HasTerrain(terrain GameTerrain, context Context) bool {
	state, _ := g.GetState(context)
	return state.Terrain == terrain
}
func (g Game) WithActor(actor Actor) Game {
	next := g
	next.Actors = append([]Actor{}, g.Actors...)
	next.UpdateActor(actor.ID, func(a Actor) Actor {
		return actor
	})
	return next
}

func (g Game) WithoutActionFilterEval() Game {
	next := g
	next.disableActionFilterEval = true
	return next
}

func (g *Game) FilterModifiers(predicate func(modifier Transaction[Modifier]) bool) {
	filtered := g.Modifiers[:0]
	for _, m := range g.Modifiers {
		if predicate(m) {
			filtered = append(filtered, m)
		}
	}

	g.Modifiers = filtered
}

func (g *Game) FilterParentModifiers(actorID uuid.UUID) {
	g.FilterModifiers(func(modifier Transaction[Modifier]) bool {
		if modifier.Context.ParentActorID != nil {
			return *modifier.Context.ParentActorID != actorID
		}

		return true
	})
}
func (g *Game) FilterParentActions(actorID uuid.UUID) {
	actions := g.Actions[:0]

	for _, tx := range g.Actions {
		parent := tx.Context.ParentActorID
		if parent == nil || *parent != actorID {
			actions = append(actions, tx)
		}
	}

	g.Actions = actions
}

func (g *Game) AddPlayer(player Player) {
	g.Players = append(g.Players, player)
}

func (g *Game) RemovePlayer(playerID uuid.UUID) {
	g.Players = slices.DeleteFunc(g.Players, func(p Player) bool {
		return p.ID == playerID
	})

	for i, _ := range g.Actors {
		g.Actors[i].PositionID = nil
	}
}

func (g *Game) UpdatePlayer(playerID uuid.UUID, updater func(Player) Player) {
	index := slices.IndexFunc(g.Players, func(p Player) bool {
		return p.ID == playerID
	})

	if index == -1 {
		return
	}

	g.Players[index] = updater(g.Players[index])
}

/**
 * [SetPosition]
 * This is a very important function
 */
func (g *Game) SetPosition(actor Actor, positionID *uuid.UUID) {
	prev := actor.PositionID
	if (prev == nil && positionID == nil) || (prev != nil && positionID != nil && *prev == *positionID) {
		return
	}

	player, ok := g.GetPlayerByID(actor.PlayerID)
	if !ok {
		return
	}

	if positionID != nil {
		curr := player.GetActorAtPosition(*positionID)
		if curr != nil {
			displaced, ok := g.GetActorByID(*curr)
			if ok {
				g.SetPosition(displaced, nil)
			}
		}

		g.UpdatePlayer(actor.PlayerID, func(p Player) Player {
			p.SetPosition(*positionID, &actor.ID)
			return p
		})
	}

	g.UpdateActor(actor.ID, func(a Actor) Actor {
		a.SetPosition(positionID)
		return a
	})

	if prev != nil {
		g.UpdatePlayer(actor.PlayerID, func(p Player) Player {
			p.SetPosition(*prev, nil)
			return p
		})
	}

	if positionID != nil {
		g.ActorSwitchInSideEffects(actor)
	}

	if positionID == nil {
		g.ActorSwitchOutSideEffects(actor)
	}
}

func (g *Game) ActorSwitchInSideEffects(actor Actor) {
	context := NewContext()
	context.SourceActorID = &actor.ID
	g.PushLog(NewLogContext("$source$ joined the battle.", context))
	t_context := Context{
		ParentActorID:  &actor.ID,
		SourceActorID:  &actor.ID,
		SourcePlayerID: &actor.PlayerID,
	}
	g.On(OnActorEnter, &t_context)
}

func (g *Game) ActorSwitchOutSideEffects(actor Actor) {
	g.FilterParentModifiers(actor.ID)
	g.FilterParentActions(actor.ID)

	context := NewContext()
	context.SourceActorID = &actor.ID
	g.PushLog(NewLogContext("$source$ left the battle.", context))
	t_context := Context{
		ParentActorID:  &actor.ID,
		SourceActorID:  &actor.ID,
		SourcePlayerID: &actor.PlayerID,
	}
	g.On(OnActorLeave, &t_context)
}

func (g *Game) SetActorPlayerIndex(actor Actor, index *int) {
	if index == nil {
		g.SetPosition(actor, nil)
	}

	player, ok := g.GetPlayerByID(actor.PlayerID)
	if !ok || *index >= len(player.Positions) {
		return
	}

	positionID := player.Positions[*index]
	g.SetPosition(actor, &positionID.ID)
}

func (g *Game) AddModifier(modifier Transaction[Modifier]) {
	g.Modifiers = append(g.Modifiers, modifier)
}

func (g *Game) AddActor(actor Actor) {
	g.Actors = append(g.Actors, actor)
}

func (g *Game) RemoveActor(actorID uuid.UUID) {
	g.Actors = slices.DeleteFunc(g.Actors, func(a Actor) bool {
		return a.ID == actorID
	})
}
func (g *Game) SetPlayerActors(playerID uuid.UUID, actors []Actor) {
	result := make([]Actor, 0)
	for _, a := range g.Actors {
		if a.PlayerID != playerID {
			actors = append(actors, a)
		}
	}

	result = append(result, actors...)
	g.Actors = result
}
func (g *Game) EnableActors(playerID uuid.UUID, ids []uuid.UUID) {
	for i, a := range g.Actors {
		if a.PlayerID != playerID {
			continue
		}

		if slices.Contains(ids, a.ID) {
			g.Actors[i].Enabled = true
		} else {
			g.Actors[i].Enabled = false
		}
	}

	g.UpdatePlayer(playerID, func(p Player) Player {
		p.Ready = len(ids) == BATTLE_TEAM_SIZE
		return p
	})
}

func (g *Game) UpdateActor(actorID uuid.UUID, updater func(Actor) Actor) {
	index := slices.IndexFunc(g.Actors, func(a Actor) bool {
		return a.ID == actorID
	})

	if index == -1 {
		return
	}

	g.Actors[index] = updater(g.Actors[index])
}

func (g *Game) SetActionCooldown(actorID uuid.UUID, actionID uuid.UUID, cooldown int) {
	g.UpdateActor(actorID, func(a Actor) Actor {
		a.SetActionCooldown(actionID, cooldown)
		return a
	})
}

func (g *Game) SortActions() {
	slices.SortFunc(g.Actions, func(a, b Transaction[Action]) int {
		if a.Mutation.Priority != b.Mutation.Priority {
			return b.Mutation.Priority - a.Mutation.Priority
		}

		if a.Mutation.Config.SubPriority != b.Mutation.Config.SubPriority {
			return b.Mutation.Config.SubPriority - a.Mutation.Config.SubPriority
		}

		a_source, ok := g.GetActorByID(*a.Context.SourceActorID)
		if !ok {
			return 1
		}
		b_source, ok := g.GetActorByID(*b.Context.SourceActorID)
		if !ok {
			return -1
		}

		a_res := a_source.Resolve(*g)
		b_res := b_source.Resolve(*g)
		return b_res.Stats[StatSpeed] - a_res.Stats[StatSpeed]
	})
}
func (g *Game) SortTriggers() {
	slices.SortFunc(g.Triggers, func(a, b Transaction[Trigger]) int {
		if a.Mutation.Priority != b.Mutation.Priority {
			return b.Mutation.Priority - a.Mutation.Priority
		}

		return 0
	})
}

func (g *Game) PushAction(transaction Transaction[Action]) bool {
	for _, t := range g.Actions {
		if t.Context.ParentActorID != nil &&
			transaction.Context.ParentActorID != nil &&
			*t.Context.ParentActorID == *transaction.Context.ParentActorID {
			return false
		}
	}

	g.Actions.Enqueue(transaction)
	g.SortActions()

	if len(g.Actions) == len(g.GetActionableActors()) {
		return true
	}

	return false
}

func (g Game) HasPlayerPrompt(playerID uuid.UUID) bool {
	for _, prompt := range g.Prompts {
		if *prompt.Context.SourcePlayerID == playerID {
			return true
		}
	}

	return false
}

func (g Game) AllPromptsReady() bool {
	if len(g.Prompts) == 0 {
		return false
	}

	for _, prompt := range g.Prompts {
		if !prompt.Ready {
			return false
		}
	}

	return true
}

func (g *Game) AddPrompt(transaction Transaction[Action]) {
	g.Prompts = append(g.Prompts, transaction)
}

func (g *Game) RemovePrompt(ID uuid.UUID) {
	g.Prompts = slices.DeleteFunc(g.Prompts, func(t Transaction[Action]) bool {
		return t.ID == ID
	})
}

func (g *Game) ReadyPrompt(ID uuid.UUID, context Context) {
	for i, prompt := range g.Prompts {
		if prompt.ID == ID {
			g.Prompts[i].Context = context
			g.Prompts[i].Ready = true
			return
		}
	}
}

func (g *Game) RunPrompt(transaction Transaction[Action]) {
	transactions := ResolveAction(g, transaction)
	g.Transactions = append(g.Transactions, transactions...)
}

func (g *Game) RunAction(transaction Transaction[Action]) {
	transactions := ResolveAction(g, transaction)
	g.Transactions = append(g.Transactions, transactions...)
}

func (g *Game) RunTrigger(transaction Transaction[Trigger]) {
	if transaction.Mutation.ModifierID != uuid.Nil {
		modifier, ok := g.GetModifierByID(transaction.Mutation.ModifierID)
		if !ok {
			return
		}

		text := fmt.Sprintf("%s: %s", strings.ToUpper(string(transaction.Mutation.On)), modifier.Name)
		g.PushLog(NewLog(text))
	}

	transactions := ResolveTrigger(*g, transaction)
	g.Transactions = append(transactions, g.Transactions...)
}
func (g *Game) On(on TriggerOn, context *Context) int {
	triggers := make([]Transaction[Trigger], 0)
	for _, trigger := range g.GetTriggers(on, context) {
		if trigger.Mutation.On == on {
			triggers = append(triggers, trigger)
		}
	}

	count := len(triggers)
	g.Triggers = append(g.Triggers, triggers...)
	return count
}

func (g *Game) JumpTransaction(transaction Transaction[GameMutation]) {
	next := Queue[GameTransaction]{transaction}
	g.Transactions = append(next, g.Transactions...)
}

func (g *Game) PushLog(log GameLog) {
	g.Log = append(g.Log, log)
}

func (g *Game) NextTurn() {
	g.Turn.Count++
	g.Turn.Phase = TurnMain
	for _, tx := range g.QueuedActions {
		action, ok := g.ActionRegistry[tx.Mutation]
		if !ok {
			continue
		}

		g.PushAction(MakeTransaction(action, tx.Context))
	}
}

type GameJSON struct {
	Status             GameStatus             `json:"status"`
	Turn               Turn                   `json:"turn"`
	ActiveTransaction  *GameActiveTransaction `json:"active_transaction"`
	Players            []Player               `json:"players"`
	Actors             []ResolvedActor        `json:"actors"`
	State              GameState              `json:"state"`
	AppliedGameStateTx []uuid.UUID            `json:"applied_game_state_tx"`

	Modifiers    []Transaction[Modifier]     `json:"modifiers"`
	Transactions []Transaction[GameMutation] `json:"-"`
	Actions      []Transaction[Action]       `json:"actions"`
	Prompt       *Transaction[Action]        `json:"prompt"`
	Triggers     []Transaction[Trigger]      `json:"-"`
	Log          []GameLog                   `json:"log"`

	QueuedActions map[uuid.UUID]Transaction[uuid.UUID] `json:"queued_actions"`
}

func getLastN[T any](s []T, n int) []T {
	if n >= len(s) {
		return s
	}
	return s[len(s)-n:]
}
func (g Game) ToJSON(playerID *uuid.UUID) GameJSON {
	resolvedMap := g.GetResolvedActors()
	resolved := make([]ResolvedActor, 0, len(g.Actors))
	for _, a := range g.Actors {
		resolved = append(resolved, resolvedMap[a.ID])
	}

	var prompt *Transaction[Action]
	if playerID != nil {
		for _, p := range g.Prompts {
			if p.Context.SourcePlayerID != nil && *p.Context.SourcePlayerID == *playerID && !p.Ready {
				prompt = &p
				break
			}
		}
	}

	status := g.Status
	if status == GameStatusIdle && len(g.Prompts) > 0 && prompt == nil {
		status = GameStatusWaiting
	}

	state, applied := g.GetState(NewContext())
	return GameJSON{
		Status:             status,
		Turn:               g.Turn,
		ActiveTransaction:  g.ActiveTransaction,
		Players:            g.Players,
		Actors:             resolved,
		State:              state,
		AppliedGameStateTx: applied,
		Modifiers:          g.Modifiers,
		Transactions:       g.Transactions,
		Actions:            g.Actions,
		Prompt:             prompt,
		Triggers:           g.Triggers,
		Log:                getLastN(g.Log, CLIENT_LOG_SIZE),
		QueuedActions:      g.QueuedActions,
	}
}

func (g Game) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.ToJSON(nil))
}
