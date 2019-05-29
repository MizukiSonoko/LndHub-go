package entity

import "time"

type PaymentHashState int

const (
	PAYMENT_HASH_STATE_UNSPECIFIED PaymentHashState = iota
	PAYMENT_HASH_STATE_PAIED
)

type User struct {
	id          string
	balance     uint
	paymentHash map[string]PaymentHashState
	tx          Transaction
	invoice     string
	address     string
}

func (u User) Id() string {
	return u.id
}

func (u *User) Getaddress() string {
	return u.address
}

func (u User) Balance() uint {
	return u.balance
}

func (u *User) UpdateBalance(n uint) {
	u.balance = n
}

func (u *User) AttachTransaction(tx Transaction) {
	u.tx = tx
}

func (u *User) AttachUserInvoice(invoice string) {
	u.invoice = invoice
}

func (u *User) UnlockFounds(invoice string) {

}

func (u *User) Invoice() string {
	return u.invoice
}

func (u *User) Txs() []Transaction {
	return u.tx
}

func (u *User) GetPaymentHashState(hash string) PaymentHashState {
	st, ok := u.paymentHash[hash]
	if !ok {
		return PAYMENT_HASH_STATE_UNSPECIFIED
	}
	return st
}

func (u *User) UpdatePaymentHashState(hash string, st PaymentHashState) {
	u.paymentHash[hash] = st
}

func NewUser(id string, balance uint) *User {
	return &User{
		id: id, balance: balance}
}

type Transaction struct {
	timestamp time.Time
	txType    string
	value     uint
	fee       uint
	memo      string
}

func NewTx(timestamp time.Time, txType string, value uint, fee uint, memo string) *Transaction {
	return &Transaction{
		timestamp: timestamp,
		txType:    txType,
		value:     value,
		fee:       fee,
		memo:      memo,
	}
}
