package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"bank/internal/db"
	"bank/internal/models"
	"bank/internal/service"
	"bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// AccountController xử lý các request liên quan đến tài khoản.
type AccountController struct {
	accountService service.AccountService
}

// NewAccountController tạo một instance mới của AccountController.
func NewAccountController(accountService service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

// *** HÀM TRỢ GIÚP MỚI ***
// getAccountIDFromHeader lấy và xác thực ID tài khoản từ header "X-User-ID".
func getAccountIDFromHeader(ctx *gin.Context) (int64, error) {
	userIDStr := ctx.GetHeader("X-User-ID")
	if userIDStr == "" {
		// Trả về lỗi cụ thể để có thể tạo AppError phù hợp ở ngoài
		return 0, errors.New("missing X-User-ID header")
	}

	accountID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid X-User-ID format, must be an integer")
	}
	return accountID, nil
}

// CreateAccount godoc
// @Summary Tạo tài khoản mới
// @Description Tạo một tài khoản ngân hàng mới với thông tin chủ sở hữu và tiền tệ.
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   account body models.CreateAccountRequest true "Thông tin tài khoản cần tạo"
// @Success 201 {object} models.AccountResponse "Tài khoản đã được tạo thành công"
// @Failure 400 {object} models.ErrorResponse "Dữ liệu đầu vào không hợp lệ"
// @Failure 409 {object} models.ErrorResponse "Tài khoản đã tồn tại"
// @Failure 500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts [post]
func (ctrl *AccountController) CreateAccount(ctx *gin.Context) {
	var req models.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		appErr := utils.NewBadRequestError("dữ liệu không hợp lệ", err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	// Các validation khác giữ nguyên...
	if !utils.IsSupportedCurrency(req.Currency) {
		appErr := utils.NewBadRequestError("loại tiền tệ không được hỗ trợ", nil)
		ctx.JSON(appErr.Code, appErr)
		return
	}
	if req.Balance < 0 {
		appErr := utils.NewBadRequestError("số dư ban đầu không thể âm", nil)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	account, err := ctrl.accountService.CreateAccount(ctx.Request.Context(), req)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			appErr := utils.NewAppError("tài khoản với các thuộc tính này đã tồn tại", http.StatusConflict)
			ctx.JSON(appErr.Code, appErr)
			return
		}
		appErr := utils.HandleServiceError(err, "không thể tạo tài khoản")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := utils.ToAccountResponse(account)
	ctx.JSON(http.StatusCreated, rsp)
}

// GetMyAccount godoc
// @Summary Lấy thông tin tài khoản của tôi
// @Description Lấy thông tin chi tiết của tài khoản đang đăng nhập dựa trên X-User-ID header.
// @Tags accounts
// @Produce  json
// @Param   X-User-ID header int true "ID Tài khoản của người dùng"
// @Success 200 {object} models.AccountResponse "Thông tin chi tiết tài khoản"
// @Failure 400 {object} models.ErrorResponse "Header X-User-ID không hợp lệ hoặc bị thiếu"
// @Failure 404 {object} models.ErrorResponse "Tài khoản không tồn tại"
// @Failure 500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts/me [get]
func (ctrl *AccountController) GetMyAccount(ctx *gin.Context) {
	// Sửa đổi: Lấy ID từ header
	accountID, err := getAccountIDFromHeader(ctx)
	if err != nil {
		appErr := utils.NewBadRequestError(err.Error(), err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	account, err := ctrl.accountService.GetAccount(ctx.Request.Context(), accountID)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			appErr := utils.NewNotFoundError("tài khoản không tồn tại", err)
			ctx.JSON(appErr.Code, appErr)
			return
		}
		appErr := utils.HandleServiceError(err, "không thể lấy thông tin tài khoản")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := utils.ToAccountResponse(account)
	ctx.JSON(http.StatusOK, rsp)
}

// ListAccounts godoc
// @Summary [Admin] Liệt kê các tài khoản
// @Description Lấy danh sách tất cả các tài khoản với phân trang (chỉ dành cho admin).
// @Tags accounts
// @Produce  json
// @Param page_id query int true "Page ID (bắt đầu từ 1)" default(1)
// @Param page_size query int true "Page Size (tối thiểu 5, tối đa 20)" default(10)
// @Success 200 {array} models.AccountResponse "Danh sách tài khoản"
// @Failure 400 {object} models.ErrorResponse "Tham số không hợp lệ"
// @Failure 500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts [get]
func (ctrl *AccountController) ListAccounts(ctx *gin.Context) {
	var req models.ListAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		appErr := utils.NewBadRequestError("tham số phân trang không hợp lệ", err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	accounts, err := ctrl.accountService.ListAccounts(ctx.Request.Context(), req)
	if err != nil {
		appErr := utils.HandleServiceError(err, "không thể liệt kê tài khoản")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := make([]models.AccountResponse, len(accounts))
	for i, acc := range accounts {
		rsp[i] = utils.ToAccountResponse(acc)
	}
	ctx.JSON(http.StatusOK, rsp)
}

// DepositToMyAccount godoc
// @Summary Nạp tiền vào tài khoản của tôi
// @Description Nạp một số tiền vào tài khoản được chỉ định bởi X-User-ID header.
// @Tags accounts
// @Accept   json
// @Produce  json
// @Param    X-User-ID header int true "ID Tài khoản của người dùng"
// @Param    deposit_request body models.DepositRequest true "Thông tin nạp tiền"
// @Success  200 {object} models.AccountResponse "Tài khoản sau khi nạp tiền"
// @Failure  400 {object} models.ErrorResponse "Dữ liệu không hợp lệ hoặc header bị thiếu/sai"
// @Failure  404 {object} models.ErrorResponse "Tài khoản không tồn tại"
// @Failure  422 {object} models.ErrorResponse "Không thể xử lý yêu cầu (ví dụ: tiền tệ không khớp, tài khoản không hoạt động)"
// @Failure  500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts/deposit [post]
func (ctrl *AccountController) DepositToMyAccount(ctx *gin.Context) {
	// Sửa đổi: Lấy ID từ header
	accountID, err := getAccountIDFromHeader(ctx)
	if err != nil {
		appErr := utils.NewBadRequestError(err.Error(), err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	var bodyReq models.DepositRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		appErr := utils.NewBadRequestError("dữ liệu nạp tiền không hợp lệ", err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	// Các validation khác giữ nguyên...
	if !utils.IsSupportedCurrency(bodyReq.Currency) {
		appErr := utils.NewBadRequestError("loại tiền tệ không được hỗ trợ khi nạp tiền", nil)
		ctx.JSON(appErr.Code, appErr)
		return
	}
	if bodyReq.Amount <= 0 {
		appErr := utils.NewBadRequestError("số tiền nạp phải lớn hơn 0", nil)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	account, err := ctrl.accountService.DepositToAccount(ctx.Request.Context(), accountID, bodyReq)
	if err != nil {
		appErr := utils.HandleServiceError(err, "lỗi khi nạp tiền")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := utils.ToAccountResponse(account)
	ctx.JSON(http.StatusOK, rsp)
}

// MakePaymentOnMyAccount godoc
// @Summary Thanh toán từ tài khoản của tôi
// @Description Thực hiện thanh toán (trừ tiền) từ tài khoản được chỉ định bởi X-User-ID header.
// @Tags accounts
// @Accept   json
// @Produce  json
// @Param    X-User-ID header int true "ID Tài khoản của người dùng"
// @Param    payment_request body models.PaymentRequest true "Thông tin thanh toán"
// @Success  200 {object} models.AccountResponse "Tài khoản sau khi thanh toán"
// @Failure  400 {object} models.ErrorResponse "Dữ liệu không hợp lệ hoặc header bị thiếu/sai"
// @Failure  402 {object} models.ErrorResponse "Số dư không đủ"
// @Failure  404 {object} models.ErrorResponse "Tài khoản không tồn tại"
// @Failure  422 {object} models.ErrorResponse "Không thể xử lý yêu cầu (ví dụ: tiền tệ không khớp, tài khoản không hoạt động)"
// @Failure  500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts/payment [post]
func (ctrl *AccountController) MakePaymentOnMyAccount(ctx *gin.Context) {
	// Sửa đổi: Lấy ID từ header
	accountID, err := getAccountIDFromHeader(ctx)
	if err != nil {
		appErr := utils.NewBadRequestError(err.Error(), err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	var bodyReq models.PaymentRequest
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		appErr := utils.NewBadRequestError("dữ liệu thanh toán không hợp lệ", err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	// Các validation khác giữ nguyên...
	if !utils.IsSupportedCurrency(bodyReq.Currency) {
		appErr := utils.NewBadRequestError("loại tiền tệ không được hỗ trợ khi thanh toán", nil)
		ctx.JSON(appErr.Code, appErr)
		return
	}
	if bodyReq.Amount <= 0 {
		appErr := utils.NewBadRequestError("số tiền thanh toán phải lớn hơn 0", nil)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	account, err := ctrl.accountService.MakePayment(ctx.Request.Context(), accountID, bodyReq)
	if err != nil {
		appErr := utils.HandleServiceError(err, "lỗi khi thực hiện thanh toán")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := utils.ToAccountResponse(account)
	ctx.JSON(http.StatusOK, rsp)
}

// CloseMyAccount godoc
// @Summary Đóng tài khoản của tôi
// @Description Chuyển trạng thái của tài khoản (chỉ định bởi X-User-ID) sang 'closed'.
// @Tags accounts
// @Produce  json
// @Param    X-User-ID header int true "ID Tài khoản của người dùng"
// @Success  200 {object} models.AccountResponse "Tài khoản sau khi được đóng"
// @Failure  400 {object} models.ErrorResponse "Header X-User-ID không hợp lệ hoặc bị thiếu"
// @Failure  404 {object} models.ErrorResponse "Tài khoản không tồn tại"
// @Failure  422 {object} models.ErrorResponse "Không thể xử lý yêu cầu (ví dụ: tài khoản đã đóng)"
// @Failure  500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts/close [patch]
func (ctrl *AccountController) CloseMyAccount(ctx *gin.Context) {
	// Sửa đổi: Lấy ID từ header
	accountID, err := getAccountIDFromHeader(ctx)
	if err != nil {
		appErr := utils.NewBadRequestError(err.Error(), err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	account, err := ctrl.accountService.CloseAccount(ctx.Request.Context(), accountID)
	if err != nil {
		appErr := utils.HandleServiceError(err, "lỗi khi đóng tài khoản")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := utils.ToAccountResponse(account)
	ctx.JSON(http.StatusOK, rsp)
}

// GetMyTransactionHistory godoc
// @Summary Lấy lịch sử giao dịch của tôi
// @Description Lấy lịch sử giao dịch của tài khoản (chỉ định bởi X-User-ID) với phân trang.
// @Tags accounts
// @Produce  json
// @Param    X-User-ID header int true "ID Tài khoản của người dùng"
// @Param 	 page_id query int false "Page ID (bắt đầu từ 1)" default(1)
// @Param 	 page_size query int false "Page Size (tối thiểu 1, tối đa 50)" default(10)
// @Success  200 {array} models.TransactionHistoryResponse "Lịch sử giao dịch"
// @Failure  400 {object} models.ErrorResponse "Header hoặc tham số phân trang không hợp lệ"
// @Failure  404 {object} models.ErrorResponse "Tài khoản không tồn tại"
// @Failure  500 {object} models.ErrorResponse "Lỗi máy chủ nội bộ"
// @Router /accounts/history [get]
func (ctrl *AccountController) GetMyTransactionHistory(ctx *gin.Context) {
	// Sửa đổi: Lấy ID từ header
	accountID, err := getAccountIDFromHeader(ctx)
	if err != nil {
		appErr := utils.NewBadRequestError(err.Error(), err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	var pageReq models.ListTransactionHistoryRequest
	if err := ctx.ShouldBindQuery(&pageReq); err != nil {
		appErr := utils.NewBadRequestError("tham số phân trang không hợp lệ", err)
		ctx.JSON(appErr.Code, appErr)
		return
	}

	history, err := ctrl.accountService.GetTransactionHistory(ctx.Request.Context(), accountID, pageReq)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			appErr := utils.NewNotFoundError("tài khoản không tồn tại khi lấy lịch sử", err)
			ctx.JSON(appErr.Code, appErr)
			return
		}
		appErr := utils.HandleServiceError(err, "không thể lấy lịch sử giao dịch")
		ctx.JSON(appErr.Code, appErr)
		return
	}

	rsp := make([]models.TransactionHistoryResponse, len(history))
	for i, tx := range history {
		rsp[i] = toTransactionHistoryResponse(tx)
	}
	ctx.JSON(http.StatusOK, rsp)
}

// toTransactionHistoryResponse chuyển đổi db.TransactionHistory sang models.TransactionHistoryResponse.
func toTransactionHistoryResponse(tx db.TransactionHistory) models.TransactionHistoryResponse {
	resp := models.TransactionHistoryResponse{
		ID:                   tx.ID,
		AccountID:            tx.AccountID,
		TransactionType:      models.TransactionType(tx.TransactionType),
		TransactionTimestamp: tx.CreatedAt.Format(time.RFC3339Nano),
	}
	resp.Description = tx.Description
	if tx.Amount.Valid {
		amount := tx.Amount.Int64
		resp.Amount = &amount
	}
	if tx.Currency.Valid {
		currency := tx.Currency.String
		resp.Currency = &currency
	}
	return resp
}

// Struct 'GetAccountRequest' đã được loại bỏ vì không còn được sử dụng.
