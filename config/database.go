package config

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type Transaction struct {
	Builder *squirrel.StatementBuilderType
	Tx      *sql.Tx
}

func (t *Transaction) Commit() error {
	return t.Tx.Commit()
}

func (t *Transaction) Rollback() {
	_ = t.Tx.Rollback()
}

var DatabaseConnection *sql.DB

func OpenConnection() (err error) {
	connStr := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=disable", Environment.Database.Username, Environment.Database.Password, Environment.Database.Host, Environment.Database.Name)

	DatabaseConnection, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DatabaseConnection.Ping()
	if err != nil {
		return err
	}

	return
}

func NewTransaction(readOnly bool) (transaction *Transaction, err error) {
	tx, err := DatabaseConnection.BeginTx(context.Background(), &sql.TxOptions{ReadOnly: readOnly})
	if err != nil {
		return
	}

	builder := NewBuilder(tx)

	transaction = &Transaction{
		Builder: &builder,
		Tx:      tx,
	}

	return
}

func NewBuilder(tx *sql.Tx) squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx)
}
