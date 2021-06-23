package account

import (
	"errors"
	"fmt"
)

//Account is bank account
type Account struct {
	owner string
	balance int
}

var errNoMoney = errors.New("Can't withdraw. You don't have enough money.")

//NewAccount creates Account
func NewAccount(owner string) *Account{
	account := Account{owner: owner, balance: 0}
	return &account
}

// Balance of your account
func (a Account) Balance() int {
	return a.balance
}

// Owner of your account
func (a Account) Owner() string {
	return a.owner
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance+=amount
}

//Withdraw x amount on your account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -=amount
	return nil
}

//ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

func (a Account) String() string{
	return fmt.Sprint(a.Owner(),"'s account.\nHas : ",a.Balance())
}