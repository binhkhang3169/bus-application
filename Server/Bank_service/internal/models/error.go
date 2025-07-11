package models

import "errors"

var (
	ErrAccountNotFound      = errors.New("tài khoản không tồn tại")
	ErrInsufficientFunds    = errors.New("số dư không đủ")
	ErrInvalidAccountStatus = errors.New("trạng thái tài khoản không hợp lệ cho hành động này")
	ErrSameAccountTransfer  = errors.New("không thể chuyển tiền vào cùng một tài khoản") // Mặc dù không có chức năng chuyển tiền, nhưng để làm ví dụ
	ErrCurrencyMismatch     = errors.New("loại tiền tệ không khớp")
)
