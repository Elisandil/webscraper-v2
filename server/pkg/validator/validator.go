package validator

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/robfig/cron/v3"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Validator struct {
	cronParser cron.Parser
}

func NewValidator() *Validator {
	return &Validator{
		cronParser: cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}
}

func (v *Validator) IsNotNil(ptr interface{}, fieldName string) error {

	if ptr == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}
	return nil
}

// --- String Validation ---

func (v *Validator) ValidateRequired(value, fieldName string) error {

	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

func (v *Validator) ValidateLength(value, fieldName string, min, max int) error {
	length := len(strings.TrimSpace(value))
	if length < min || length > max {
		return fmt.Errorf("%s must be between %d and %d characters", fieldName, min, max)
	}
	return nil
}

func (v *Validator) ValidateMinLength(value, fieldName string, min int) error {

	if len(strings.TrimSpace(value)) < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	return nil
}

func (v *Validator) ValidateMaxLength(value, fieldName string, max int) error {

	if len(strings.TrimSpace(value)) > max {
		return fmt.Errorf("%s must not exceed %d characters", fieldName, max)
	}
	return nil
}

// --- Email Validation ---

func (v *Validator) ValidateEmail(email string) error {

	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// --- URL Validation ---

func (v *Validator) ValidateURL(targetURL string) error {

	if strings.TrimSpace(targetURL) == "" {
		return errors.New("URL is required")
	}

	parsedURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("URL must include protocol (http:// or https://) and valid domain")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("only HTTP and HTTPS protocols are supported")
	}

	return nil
}

// --- Cron Expression Validation ---

func (v *Validator) ValidateCronExpression(cronExpr string) error {

	if strings.TrimSpace(cronExpr) == "" {
		return errors.New("cron expression is required")
	}

	_, err := v.cronParser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}
	return nil
}

// --- Struct Validation ---

func (v *Validator) ValidateStruct(value interface{}, fieldName string) error {

	if value == nil {
		return fmt.Errorf("%s cannot be nil", fieldName)
	}
	return nil
}

// --- Password Validation ---

func (v *Validator) ValidatePassword(password string) error {

	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// --- Username Validation ---

func (v *Validator) ValidateUsername(username string) error {

	if username == "" {
		return errors.New("username is required")
	}
	if err := v.ValidateLength(username, "username", 3, 50); err != nil {
		return err
	}
	return nil
}
