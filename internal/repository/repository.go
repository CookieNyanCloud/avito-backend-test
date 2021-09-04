package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)




type FinanceRepo struct {
	db *sqlx.DB
}

func NewFinanceRepo(db *sqlx.DB) *FinanceRepo {
	return &FinanceRepo{db: db}
}

type Finance interface {
	Transaction(id int, sum float64) error
	Remittance(idFrom int, idTo int, sum float64) error
	Balance(id int) (float64,error)

	NewTransaction(idFrom int, operation string, sum float64, idTo int) error
	//GetTransactionsList(id int,)
}

func (r *FinanceRepo) NewFinanceRepo(db *sqlx.DB) *FinanceRepo {
	return &FinanceRepo{db: db}
}

const (
	financeTable = "userbalance"
	transactionTable = "transactions"
)

const (
	Minus   = "недостаточно средств"
)

const (
	transaction = "transaction"
	remittance = "remittance"
)


func (r *FinanceRepo) NewTransaction(idFrom int, operation string,sum float64, idTo int) error {
	if idTo >0 {
		query := fmt.Sprintf("INSERT INTO %s (user_id, operation,sum, user_to) values ($1, $2, $3, $4)",
			transactionTable)
		_, err := r.db.Exec(query,idFrom,operation,sum, idTo)
		if err != nil {
			return err
		}
	}else {
		query := fmt.Sprintf("INSERT INTO %s (user_id, operation, sum, user_to) values ($1, $2, $3, NULL)",
			transactionTable)
		_, err := r.db.Exec(query,idFrom,operation,sum)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *FinanceRepo) Balance (id int) (float64, error)  {
	var currentBalance float64
	query := fmt.Sprintf(`SELECT balance FROM %s WHERE id=$1`, financeTable)
	err := r.db.Get(&currentBalance, query, id)
	if err != nil {
		return -1, err
	}
	return currentBalance, nil
}

func (r *FinanceRepo) Transaction(id int, sum float64) error {

	currentBalance, err:= r.Balance(id)
	if err != nil {
		return err
	}

	newBalance := currentBalance + sum
	if newBalance >= 0 {
		query := fmt.Sprintf("UPDATE %s SET balance = $1  WHERE id = $2", financeTable)
		_, err = r.db.Exec(query, newBalance, id)
		if err != nil {
			return err
		}
		err = r.NewTransaction(id, transaction, sum,-1)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(Minus)
}

func (r *FinanceRepo) Remittance(idFrom int, idTo int, sum float64) error {
	currentBalanceFrom, err := r.Balance(idFrom)
	if err != nil {
		return err
	}

	currentBalanceTo, err:= r.Balance(idTo)
	if err != nil {
		return err
	}

	newBalanceFrom:= currentBalanceFrom - sum
	newBalanceTo:= currentBalanceTo + sum
	if newBalanceFrom >= 0 {
		query := fmt.Sprintf("UPDATE %s SET balance = $1  WHERE id = $2",
			financeTable)
		_, err = r.db.Exec(query, newBalanceFrom, idFrom)
		if err != nil {
			return err
		}

		query = fmt.Sprintf("UPDATE %s SET balance = $1  WHERE id = $2",
			financeTable)
		_, err = r.db.Exec(query, newBalanceTo, idTo)
		if err != nil {
			return err
		}

		err = r.NewTransaction(idFrom, transaction, sum, idTo)
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New(Minus)
}

