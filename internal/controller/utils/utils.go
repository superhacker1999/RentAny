package utils

import (
	"RentAny/internal/types"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phone = strings.ReplaceAll(phone, " ", "")

	return re.Find([]byte(phone)) != nil
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	rules := [4]string{"([a-z])+", "([A-Z])+", "([0-9])+", "([!@#$%^&*.?-])+"}

	for _, rule := range rules {
		if !regexp.MustCompile(rule).MatchString(password) {
			return false
		}
	}

	return true
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ValidateLoginCredentials(sl validator.StructLevel) {
	creds := sl.Current().Interface().(types.LoginCredentials)

	// Проверяем, что заполнено хотя бы одно из полей: email или телефон
	if creds.Email == "" && creds.Phone == "" {
		sl.ReportError(creds.Email, "email", "Email", "either-email-or-phone", "")
		sl.ReportError(creds.Phone, "phone", "Phone", "either-email-or-phone", "")
	}
}
