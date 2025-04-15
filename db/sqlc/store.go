package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) exeTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) //创建一个数据库事务变量，得到指针
	if err != nil {
		return err
	}
	q := New(tx) //通过sql.Tx指针创建Queries结构体指针
	err = fn(q)  //执行事务函数
	if err != nil {
		RBerr := tx.Rollback()
		if RBerr != nil {
			return fmt.Errorf(" fn(q) error:%v tx.Rollback() error: %v", err, RBerr)
		}
		return err
	}
	return tx.Commit()
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.exeTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	Amount:    -arg.Amount,
		// 	AccountID: arg.FromAccountID,
		// })
		// if err != nil {
		// 	return err
		// }
		// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		// 	Amount:    arg.Amount,
		// 	AccountID: arg.ToAccountID,
		// })
		// if err != nil {
		// 	return err
		// }

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err

	})
	if err != nil {
		return result, err
	}
	return result, err
}

func addMoney(ctx context.Context, q *Queries,
	account1ID int64, amount1 int64,
	account2ID int64, amount2 int64) (updateAccount1 Account, updateAccount2 Account, err error) {
	updateAccount1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount:    amount1,
		AccountID: account1ID,
	})
	if err != nil {
		return
	}
	updateAccount2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount:    amount2,
		AccountID: account2ID,
	})
	if err != nil {
		return
	}

	return
}
