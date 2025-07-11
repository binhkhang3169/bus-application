package service

import (
	"context"
	"database/sql"
	"errors" // Import package errors
	"fmt"
	"log"
	"net/http"
	"time"

	"bank/internal/db"
	"bank/internal/models"
	"bank/internal/repository"
	"bank/pkg/kafkaclient"
	"bank/utils"
)

// AccountService định nghĩa interface cho các nghiệp vụ liên quan đến tài khoản.
type AccountService interface {
	CreateAccount(ctx context.Context, req models.CreateAccountRequest) (db.Account, error)
	GetAccount(ctx context.Context, id int64) (db.Account, error)
	DepositToAccount(ctx context.Context, accountID int64, req models.DepositRequest) (db.Account, error)
	MakePayment(ctx context.Context, accountID int64, req models.PaymentRequest) (db.Account, error)
	CloseAccount(ctx context.Context, accountID int64) (db.Account, error)
	ListAccounts(ctx context.Context, req models.ListAccountsRequest) ([]db.Account, error)
	GetTransactionHistory(ctx context.Context, accountID int64, req models.ListTransactionHistoryRequest) ([]db.TransactionHistory, error)
}

type accountService struct {
	repo      repository.AccountRepository
	publisher *kafkaclient.Publisher // << ADDED

}

// NewAccountService tạo một instance mới của AccountService.
func NewAccountService(repo repository.AccountRepository, publisher *kafkaclient.Publisher) AccountService {
	return &accountService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req models.CreateAccountRequest) (db.Account, error) {
	var account db.Account
	err := s.repo.ExecTx(ctx, func(q *db.Queries) error {
		createAccountArg := db.CreateAccountParams{
			OwnerName: req.OwnerName,
			Balance:   req.Balance,
			Currency:  req.Currency,
			Status:    string(utils.AccountStatusActive),
		}
		if createAccountArg.Balance < 0 {
			createAccountArg.Balance = 0
		}

		var errCreate error
		account, errCreate = q.CreateAccount(ctx, createAccountArg)
		if errCreate != nil {
			return utils.NewInternalServerError("không thể tạo tài khoản", errCreate)
		}

		_, errLog := q.CreateTransactionHistory(ctx, db.CreateTransactionHistoryParams{
			AccountID:       account.ID,
			TransactionType: string(models.TransactionTypeCreateAccount),
			Amount:          sql.NullInt64{Int64: account.Balance, Valid: true},
			Currency:        sql.NullString{String: account.Currency, Valid: true},
			Description:     fmt.Sprintf("Account created for %s. Initial balance: %d %s", account.OwnerName, account.Balance, account.Currency),
		})
		if errLog != nil {
			return utils.NewInternalServerError("không thể ghi lịch sử giao dịch khi tạo tài khoản", errLog)
		}
		return nil
	})

	if err != nil {
		return db.Account{}, err
	}
	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	account, err := s.repo.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, models.ErrAccountNotFound
		}
		return db.Account{}, utils.NewInternalServerError("không thể lấy thông tin tài khoản", err)
	}
	return account, nil
}

func (s *accountService) ListAccounts(ctx context.Context, req models.ListAccountsRequest) ([]db.Account, error) {
	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := s.repo.ListAccounts(ctx, arg)
	if err != nil {
		return nil, utils.NewInternalServerError("không thể liệt kê tài khoản", err)
	}
	return accounts, nil
}

func (s *accountService) DepositToAccount(ctx context.Context, accountID int64, req models.DepositRequest) (db.Account, error) {
	var updatedAccount db.Account
	err := s.repo.ExecTx(ctx, func(q *db.Queries) error {
		acc, err := q.GetAccountForUpdate(ctx, accountID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.ErrAccountNotFound
			}
			return err
		}
		if acc.Status != string(utils.AccountStatusActive) {
			return models.ErrInvalidAccountStatus
		}
		if acc.Currency != req.Currency {
			return models.ErrCurrencyMismatch
		}
		if req.Amount <= 0 {
			return utils.NewAppError("số tiền nạp phải dương", http.StatusBadRequest) // Or a specific error type
		}

		updatedAccount, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     accountID,
			Amount: req.Amount,
		})
		if err != nil {
			return err
		}

		_, errLog := q.CreateTransactionHistory(ctx, db.CreateTransactionHistoryParams{
			AccountID:       updatedAccount.ID,
			TransactionType: string(models.TransactionTypeDeposit),
			Amount:          sql.NullInt64{Int64: req.Amount, Valid: true},
			Currency:        sql.NullString{String: req.Currency, Valid: true},
			Description:     fmt.Sprintf("Deposited %d %s. New balance: %d %s", req.Amount, req.Currency, updatedAccount.Balance, updatedAccount.Currency),
		})
		if errLog != nil {
			return utils.NewInternalServerError("không thể ghi lịch sử giao dịch khi nạp tiền", errLog)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) || errors.Is(err, models.ErrInvalidAccountStatus) || errors.Is(err, models.ErrCurrencyMismatch) {
			return db.Account{}, utils.NewAppError(err.Error(), utils.DetermineStatusCode(err))
		}
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			return db.Account{}, appErr
		}
		return db.Account{}, utils.NewInternalServerError("lỗi khi nạp tiền", err)
	}

	s.publishTransactionNotification(
		context.Background(),
		updatedAccount,
		"DEPOSIT_SUCCESS",
		"Nạp tiền thành công",
		fmt.Sprintf("Bạn đã nạp thành công %d %s vào tài khoản. Số dư mới: %d %s.", req.Amount, req.Currency, updatedAccount.Balance, updatedAccount.Currency),
	)

	return updatedAccount, nil
}

func (s *accountService) MakePayment(ctx context.Context, accountID int64, req models.PaymentRequest) (db.Account, error) {
	var updatedAccount db.Account
	err := s.repo.ExecTx(ctx, func(q *db.Queries) error {
		acc, err := q.GetAccountForUpdate(ctx, accountID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return models.ErrAccountNotFound
			}
			return err
		}
		if acc.Status != string(utils.AccountStatusActive) {
			return models.ErrInvalidAccountStatus
		}
		if acc.Currency != req.Currency {
			return models.ErrCurrencyMismatch
		}
		if req.Amount <= 0 {
			return utils.NewAppError("số tiền thanh toán phải dương", http.StatusBadRequest) // Or a specific error type
		}
		if acc.Balance < req.Amount {
			return models.ErrInsufficientFunds
		}

		updatedAccount, err = q.AddAccountBalance(ctx, db.AddAccountBalanceParams{
			ID:     accountID,
			Amount: -req.Amount,
		})
		if err != nil {
			return err
		}

		_, errLog := q.CreateTransactionHistory(ctx, db.CreateTransactionHistoryParams{
			AccountID:       updatedAccount.ID,
			TransactionType: string(models.TransactionTypePayment),
			Amount:          sql.NullInt64{Int64: req.Amount, Valid: true}, // Log positive amount
			Currency:        sql.NullString{String: req.Currency, Valid: true},
			Description:     fmt.Sprintf("Payment of %d %s made. New balance: %d %s", req.Amount, req.Currency, updatedAccount.Balance, updatedAccount.Currency),
		})
		if errLog != nil {
			return utils.NewInternalServerError("không thể ghi lịch sử giao dịch khi thanh toán", errLog)
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) || errors.Is(err, models.ErrInvalidAccountStatus) || errors.Is(err, models.ErrInsufficientFunds) || errors.Is(err, models.ErrCurrencyMismatch) {
			return db.Account{}, utils.NewAppError(err.Error(), utils.DetermineStatusCode(err))
		}
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			return db.Account{}, appErr
		}
		return db.Account{}, utils.NewInternalServerError("lỗi khi thanh toán", err)
	}

	s.publishTransactionNotification(
		context.Background(),
		updatedAccount,
		"PAYMENT_SUCCESS",
		"Thanh toán thành công",
		fmt.Sprintf("Thực hiện thanh toán %d %s thành công. Số dư còn lại: %d %s.", req.Amount, req.Currency, updatedAccount.Balance, updatedAccount.Currency),
	)

	return updatedAccount, nil
}

func (s *accountService) CloseAccount(ctx context.Context, accountID int64) (db.Account, error) {
	var updatedAccount db.Account

	accInitial, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Account{}, models.ErrAccountNotFound
		}
		return db.Account{}, utils.NewInternalServerError("không thể lấy thông tin tài khoản để đóng (pre-check)", err)
	}

	if accInitial.Status == string(utils.AccountStatusClosed) {
		return accInitial, nil
	}
	if accInitial.Status != string(utils.AccountStatusActive) {
		return db.Account{}, models.ErrInvalidAccountStatus
	}
	// Optional: Check if accInitial.Balance > 0 and prevent closing or trigger other workflows

	errTx := s.repo.ExecTx(ctx, func(q *db.Queries) error {
		var errUpdate error
		updatedAccount, errUpdate = q.UpdateAccountStatus(ctx, db.UpdateAccountStatusParams{
			ID:     accountID,
			Status: string(utils.AccountStatusClosed),
		})
		if errUpdate != nil {
			return utils.NewInternalServerError("không thể cập nhật trạng thái tài khoản", errUpdate)
		}

		_, errLog := q.CreateTransactionHistory(ctx, db.CreateTransactionHistoryParams{
			AccountID:       updatedAccount.ID,
			TransactionType: string(models.TransactionTypeCloseAccount),
			Amount:          sql.NullInt64{Valid: false},  // No amount for closing
			Currency:        sql.NullString{Valid: false}, // No currency for closing
			Description:     fmt.Sprintf("Account %d closed. Final balance: %d %s", updatedAccount.ID, updatedAccount.Balance, updatedAccount.Currency),
		})
		if errLog != nil {
			return utils.NewInternalServerError("không thể ghi lịch sử giao dịch khi đóng tài khoản", errLog)
		}
		return nil
	})

	if errTx != nil {
		var appErr *utils.AppError
		if errors.As(errTx, &appErr) {
			return db.Account{}, errTx
		}
		return db.Account{}, utils.NewInternalServerError("lỗi khi đóng tài khoản (trong tx)", errTx)
	}
	return updatedAccount, nil
}

func (s *accountService) GetTransactionHistory(ctx context.Context, accountID int64, req models.ListTransactionHistoryRequest) ([]db.TransactionHistory, error) {
	account, err := s.repo.GetAccount(ctx, accountID) // Verify account exists
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrAccountNotFound
		}
		return nil, utils.NewInternalServerError("không thể xác minh tài khoản trước khi lấy lịch sử", err)
	}

	arg := db.ListTransactionHistoryByAccountIDParams{
		AccountID: account.ID,
		Limit:     int32(req.PageSize),
		Offset:    int32((req.PageID - 1) * req.PageSize),
	}
	history, err := s.repo.ListTransactionHistoryByAccountID(ctx, arg)
	if err != nil {
		return nil, utils.NewInternalServerError("không thể lấy lịch sử giao dịch", err)
	}
	return history, nil
}

// DetermineStatusCode giúp ánh xạ lỗi nghiệp vụ sang HTTP status code
// (đã chuyển vào utils/error.go)
func (s *accountService) publishTransactionNotification(ctx context.Context, account db.Account, notiType, title, message string) {
	go func() {
		bgCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		notificationTopic := "notifications_topic" // Nên lấy từ config

		userID := account.OwnerName // Lấy OwnerName làm UserID

		event := kafkaclient.NotificationEvent{
			UserID:  &userID,
			Type:    notiType,
			Title:   title,
			Message: message,
		}

		if err := s.publisher.Publish(bgCtx, notificationTopic, []byte(userID), event); err != nil {
			log.Printf("CRITICAL: Failed to publish transaction notification for account %d: %v", account.ID, err)
		}
	}()
}
