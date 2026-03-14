package validator

import (
	"regexp"
)

//
// type (
// 	Validator struct {
// 		v *validator.Validate
// 	}
//
// 	FieldError struct {
// 		Field   string `json:"field"`
// 		Rule    string `json:"rule"`
// 		Message string `json:"message"`
// 	}
// )
//
// func getErrorMessage(e validator.FieldError) string {
// 	switch e.Tag() {
// 	case "required":
// 		return fmt.Sprintf("%s is required", strings.ToLower(e.Field()))
// 	case "min":
// 		return fmt.Sprintf("%s must be atleast %s characters", strings.ToLower(e.Field()), e.Param())
// 	case "max":
// 		return fmt.Sprintf("%s must not exceed %s characters", strings.ToLower(e.Field()), e.Param())
// 	case "username_regex":
// 		return fmt.Sprintf("%s can only contain letters, number and underscores (3-20 characters)", strings.ToLower(e.Field()))
// 	case "password":
// 		return fmt.Sprintf("%s must be 8-32 chars, include uppercase, lowercase, number, and special char", strings.ToLower(e.Field()))
// 	default:
// 		return fmt.Sprintf("%s is invalid (%s)", strings.ToLower(e.Field()), e.Tag())
// 	}
// }
//
// func New() *Validator {
// 	v := validator.New()
//
// 	if err := v.RegisterValidation("username_regex", func(fl validator.FieldLevel) bool {
// 		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, fl.Field().String())
// 		return matched
// 	}); err != nil {
// 		log.Fatal("Failed to register validation \"username_regex\"")
// 	}
//
// 	if err := v.RegisterValidation("email", func(fl validator.FieldLevel) bool {
// 		return validatePassword(fl.Field().String())
// 	}); err != nil {
// 		log.Fatal("Failed to register validation \"password\"")
// 	}
// 	return &Validator{v}
// }

func ValidateUsername(u string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`).MatchString(u)
}

func ValidateEmail(e string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(e)
}

func ValidatePassword(p string) bool {
	if len(p) < 8 || len(p) > 32 {
		return false
	}

	hasLower := regexp.MustCompile(`[a-z]`).MatchString(p)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(p)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(p)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(p)

	return hasLower && hasUpper && hasNumber && hasSpecial
}

func ValidateStatus(s string) bool {
	if s == "pending" || s == "active" || s == "completed" {
		return true
	}
	return false
}
