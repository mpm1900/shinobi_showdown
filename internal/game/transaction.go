package game

import (
	"github.com/google/uuid"
)

type Transaction[M any] struct {
	ID       uuid.UUID `json:"ID"`
	Ready    bool      `json:"ready"`
	Context  Context   `json:"context"`
	Priority int       `json:"priority"`
	Mutation M         `json:"mutation"`
}

func MakeTransaction[M any](
	mutation M,
	context Context,
) Transaction[M] {
	return Transaction[M]{
		ID:       uuid.New(),
		Context:  context,
		Mutation: mutation,
		Ready:    true,
	}
}

func CheckTransaction[P any, I any, O any](
	parent P,
	input I,
	transaction Transaction[Mutation[P, I, O]],
) bool {
	if transaction.Mutation.Filter == nil {
		return true
	}
	return transaction.Mutation.Filter(parent, input, transaction.Context)
}

func ResolveTransaction[P any, I any, O any](
	parent P,
	input I,
	transaction Transaction[Mutation[P, I, O]],
	fallback O,
) (O, bool) {
	if !CheckTransaction(parent, input, transaction) {
		return fallback, false
	}

	return transaction.Mutation.Delta(parent, input, transaction.Context), true
}

func ResolveTransactionFn[P any, I any, O any](
	parent P,
	input I,
	transaction Transaction[Mutation[P, I, O]],
	fallback func(I) O,
) (O, bool) {
	return ResolveTransaction(parent, input, transaction, fallback(input))
}
