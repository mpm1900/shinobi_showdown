package game

type TransactionBuilder struct {
	transactions []GameTransaction
}

func NewTransactionBuilder() *TransactionBuilder {
	return &TransactionBuilder{
		transactions: make([]GameTransaction, 0),
	}
}

func (t *TransactionBuilder) Push(tx []GameTransaction) {
	t.transactions = append(t.transactions, tx...)
}

func (t *TransactionBuilder) PushOne(tx GameTransaction) {
	t.transactions = append(t.transactions, tx)
}

func (t *TransactionBuilder) Build() []GameTransaction {
	return t.transactions
}
