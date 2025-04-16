package utils

func VaildatorCurrency(currency string) bool {
	switch currency {
	case "USD", "EUR", "CAD", "JPY", "GBP":
		return true
	}
	return false
}
