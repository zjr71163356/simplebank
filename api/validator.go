package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/zjr71163356/simplebank/utils"
)

var vaildatorCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if v, ok := fl.Field().Interface().(string); ok {
		return utils.VaildatorCurrency(v)
	}
	return false
}
