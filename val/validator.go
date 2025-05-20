package val

import (
	"fmt"
	"regexp"
)

const (
	usernameRegex = `^[a-zA-Z0-9_]{3,20}$`
	passwordRegex = `^[a-zA-Z0-9!@#\$%\^&\*]{6,20}$`
	fullNameRegex = `^[a-zA-Z\s]{1,50}$`
	emailRegex    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func ValidateUsername(username string) error {
	if err := VaildateStringLen(username, 3, 20); err != nil {
		return err
	}
	if !matchRegex(usernameRegex, username) {
		return fmt.Errorf("用户名不符合要求: %s", username)
	}
	return nil
}

func ValidatePassword(password string) error {
	if err := VaildateStringLen(password, 6, 20); err != nil {
		return err
	}
	if !matchRegex(passwordRegex, password) {
		return fmt.Errorf("密码不符合要求: %s", password)
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if err := VaildateStringLen(fullName, 1, 20); err != nil {
		return err
	}
	if !matchRegex(fullNameRegex, fullName) {
		return fmt.Errorf("全名不符合要求: %s", fullName)
	}
	return nil
}

func ValidateEmail(email string) error {
	if !matchRegex(emailRegex, email) {
		return fmt.Errorf("电子邮件不符合要求: %s", email)
	}
	return nil
}

func matchRegex(pattern, value string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

func VaildateStringLen(value string, min int, max int) error {
	length := len(value)
	if length < min || length > max {
		return fmt.Errorf("the length of %s should be between %d and %d", value, min, max)
	}
	return nil
}
