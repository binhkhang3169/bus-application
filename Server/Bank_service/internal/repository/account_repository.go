package repository

import (
	db "bank/internal/db" // Import package db từ sqlc
	"bank/utils"
	"context"
	"database/sql" // Quan trọng: import database/sql
)

// AccountRepository định nghĩa interface cho các thao tác CSDL liên quan đến Account.
// Sử dụng Querier interface từ sqlc giúp dễ dàng mock và test.
type AccountRepository interface {
	CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error)
	GetAccount(ctx context.Context, id int64) (db.Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (db.Account, error)
	ListAccounts(ctx context.Context, arg db.ListAccountsParams) ([]db.Account, error)
	UpdateAccountBalance(ctx context.Context, arg db.UpdateAccountBalanceParams) (db.Account, error)
	UpdateAccountStatus(ctx context.Context, arg db.UpdateAccountStatusParams) (db.Account, error)
	AddAccountBalance(ctx context.Context, arg db.AddAccountBalanceParams) (db.Account, error)
	// ExecTx dùng để thực thi một hàm callback trong một transaction CSDL
	ExecTx(ctx context.Context, fn func(*db.Queries) error) error
	ListTransactionHistoryByAccountID(ctx context.Context, arg db.ListTransactionHistoryByAccountIDParams) ([]db.TransactionHistory, error)
}

type Store interface {
	GetQuerier() *db.Queries
	GetDB() *sql.DB
}

type SQLStore struct {
	db      *sql.DB
	queries *db.Queries
}

func NewStore(dbConn *sql.DB) *SQLStore {
	return &SQLStore{
		db:      dbConn,
		queries: db.New(dbConn),
	}
}

func (s *SQLStore) GetQuerier() *db.Queries {
	return s.queries
}

func (s *SQLStore) GetDB() *sql.DB {
	return s.db
}

// SQLAccountRepository triển khai AccountRepository sử dụng sqlc.Queries.
type SQLAccountRepository struct {
	*db.Queries         // Nhúng Queries từ sqlc để có tất cả các phương thức query
	db          *sql.DB // Cần đối tượng *sql.DB để quản lý transaction
}

// NewAccountRepository tạo một instance mới của SQLAccountRepository.
func NewAccountRepository(store Store) AccountRepository { // store là interface từ sqlc (db.go)
	return &SQLAccountRepository{
		Queries: store.GetQuerier(), // Lấy Querier từ Store
		db:      store.GetDB(),      // Lấy *sql.DB từ Store
	}
}

// ExecTx thực thi một hàm callback trong một transaction CSDL.
// Hàm này rất quan trọng cho các nghiệp vụ cần nhiều thao tác CSDL (ví dụ: chuyển tiền).
func (r *SQLAccountRepository) ExecTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.New(tx) // Tạo một Querier mới với transaction
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return utils.NewTransactionError(err, rbErr) // Lỗi gốc và lỗi rollback
		}
		return err
	}

	return tx.Commit()
}

func (r *SQLAccountRepository) ListTransactionHistoryByAccountID(ctx context.Context, arg db.ListTransactionHistoryByAccountIDParams) ([]db.TransactionHistory, error) {
	return r.Queries.ListTransactionHistoryByAccountID(ctx, arg)
}

// Các phương thức của SQLAccountRepository sẽ gọi trực tiếp các phương thức từ db.Queries
// Ví dụ:
// func (r *SQLAccountRepository) CreateAccount(ctx context.Context, arg db.CreateAccountParams) (db.Account, error) {
// 	return r.Queries.CreateAccount(ctx, arg)
// }
// ... và tương tự cho các phương thức khác.
// sqlc đã tự sinh các phương thức này trên *db.Queries, nên chúng ta không cần viết lại nếu nhúng Queries.
// Tuy nhiên, nếu bạn muốn thêm logic logging hoặc xử lý lỗi cụ thể ở tầng repository, bạn có thể override chúng.

// Ví dụ về việc override nếu cần thêm logic:
// func (r *SQLAccountRepository) GetAccount(ctx context.Context, id int64) (db.Account, error) {
// 	account, err := r.Queries.GetAccount(ctx, id)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return db.Account{}, model.ErrAccountNotFound
// 		}
// 		// log error
// 		return db.Account{}, err
// 	}
// 	return account, nil
// }
