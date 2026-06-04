package game

import (
	"slices"
	"sort"

	"github.com/google/uuid"
)

const ModifierDurationInf = -1

type ModifierMetadata struct {
	Hazard bool
	Status bool
}

type Modifier struct {
	ID                uuid.UUID  `json:"ID"`
	GroupID           *uuid.UUID `json:"group_ID"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	ParentDescription string     `json:"parent_description"`
	Delay             int        `json:"-"`
	Duration          int        `json:"duration"`
	Icon              string     `json:"icon"`
	Show              bool       `json:"show"`

	// flags
	Status  bool `json:"status"`
	Weather bool `json:"weather"`
	Terrain bool `json:"terrain"`

	ActorMutations     []ActorMutation     `json:"-"`
	GameStateMutations []GameStateMutation `json:"-"`
	Triggers           []Trigger           `json:"-"`
}

func (m *Modifier) DecrementTimers() {
	if m.Duration > 0 {
		m.Duration -= 1
	}
	if m.Delay > 0 {
		m.Delay -= 1
	}
}

func ResolveTrigger(game Game, transaction Transaction[Trigger]) []Transaction[GameMutation] {
	transactions := []Transaction[GameMutation]{}
	if transaction.Mutation.Delta == nil {
		return transactions
	}

	targets := game.GetTargets(transaction.Context)
	for _, target := range targets {
		resolved := target.Resolve(game)
		if resolved.HasImmunity(transaction.Mutation.ModifierID) {
			log_ctx := MakeContextForActor(target)
			log_mut := AddLogs(NewLogContext("$source$ was immune.", log_ctx))
			log_tx := MakeTransaction(log_mut, log_ctx)
			transactions = append(transactions, log_tx)
		}
	}

	if len(transactions) > 0 {
		return transactions
	}

	return transaction.Mutation.Delta(game, game, transaction.Context)
}

func CheckModifierForActor(tx Transaction[Modifier], game Game, actor Actor) bool {
	game.Actors = []Actor{actor}
	for _, mut := range tx.Mutation.ActorMutations {
		if mut.Filter(game, actor, tx.Context) {
			return true
		}
	}

	return false
}

func GetAllActorMutations(g Game, bypassModifiers bool) ([]ActorMutation, []Transaction[Modifier]) {
	if bypassModifiers {
		mutations := slices.Clone(specialMutations)
		sort.SliceStable(mutations, func(i, j int) bool {
			return mutations[i].Priority < mutations[j].Priority
		})
		return mutations, nil
	}

	transactions := g.GetActiveModifierTransactions()
	mutations := make([]ActorMutation, 0, (len(transactions)*2)+len(specialMutations))

	for i := range transactions {
		tx := &transactions[i]
		for j := range tx.Mutation.ActorMutations {
			mut := tx.Mutation.ActorMutations[j]
			mut.TransactionID = &tx.ID
			mutations = append(mutations, mut)
		}
	}

	mutations = append(mutations, specialMutations...)

	sort.SliceStable(mutations, func(i, j int) bool {
		return mutations[i].Priority < mutations[j].Priority
	})

	return mutations, transactions
}

func GetAllGameStateMutations(g Game) ([]GameStateMutation, []Transaction[Modifier]) {
	transactions := g.GetActiveModifierTransactions()
	mutations := make([]GameStateMutation, 0, len(transactions))

	for i := range transactions {
		tx := &transactions[i]
		for j := range tx.Mutation.GameStateMutations {
			mut := tx.Mutation.GameStateMutations[j]
			mut.TransactionID = &tx.ID
			mutations = append(mutations, mut)
		}
	}

	sort.SliceStable(mutations, func(i, j int) bool {
		return mutations[i].Priority < mutations[j].Priority
	})

	return mutations, transactions
}
