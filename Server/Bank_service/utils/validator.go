package utils

import "github.com/go-playground/validator/v10"

// ValidCurrency là một custom validator function cho Gin.
// Nó kiểm tra xem giá trị của trường có phải là một loại tiền tệ được hỗ trợ không.
func ValidCurrency(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedCurrency(currency)
	}
	return false
}
