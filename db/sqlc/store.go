package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
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

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) exeTx(ctx context.Context, fn func(q *Queries) error) error {
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

func (store *Store) TransferTx(ctx context.Context, q *Queries, arg TransferTxParams) (TransferTxResult, error) {
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

		return err

	})
	if err != nil {
		return result, err
	}
	return result, err
}
