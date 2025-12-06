package crypto

import (
	"errors"
	"fmt"
)

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user account is deactivated")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenRevoked       = errors.New("token has been revoked")
	ErrUnauthorized       = errors.New("unauthorized access")

	// Validation errors
	ErrInvalidInput  = errors.New("invalid input")
	ErrRequiredField = errors.New("required field missing")
	ErrInvalidFormat = errors.New("invalid format")
	ErrInvalidURL    = errors.New("invalid URL")
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrInvalidCron   = errors.New("invalid cron expression")

	// Resource errors
	ErrResourceNotFound = errors.New("resource not found")
	ErrAlreadyExists    = errors.New("resource already exists")
	ErrDuplicate        = errors.New("duplicate entry")

	// Database errors
	ErrDatabase     = errors.New("database error")
	ErrCreateFailed = errors.New("failed to create resource")
	ErrUpdateFailed = errors.New("failed to update resource")
	ErrDeleteFailed = errors.New("failed to delete resource")
	ErrQueryFailed  = errors.New("failed to query database")

	// General errors
	ErrInternal       = errors.New("internal server error")
	ErrNotImplemented = errors.New("not implemented")
)

const (
	CodeValidation     = "VALIDATION_ERROR"
	CodeAuthentication = "AUTH_ERROR"
	CodeAuthorization  = "AUTHZ_ERROR"
	CodeNotFound       = "NOT_FOUND"
	CodeConflict       = "CONFLICT"
	CodeDatabase       = "DATABASE_ERROR"
	CodeInternal       = "INTERNAL_ERROR"
	CodeBadRequest     = "BAD_REQUEST"
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func New(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// --- Specific Error Constructors ---

func ValidationError(message string) *AppError {
	return &AppError{
		Code:    CodeValidation,
		Message: message,
		Err:     ErrInvalidInput,
	}
}

func AuthenticationError(message string) *AppError {
	return &AppError{
		Code:    CodeAuthentication,
		Message: message,
		Err:     ErrInvalidCredentials,
	}
}

func NotFoundError(resource string) *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Err:     ErrResourceNotFound,
	}
}

func ConflictError(message string) *AppError {
	return &AppError{
		Code:    CodeConflict,
		Message: message,
		Err:     ErrAlreadyExists,
	}
}

func DatabaseError(operation string, err error) *AppError {
	return &AppError{
		Code:    CodeDatabase,
		Message: fmt.Sprintf("database %s failed", operation),
		Err:     err,
	}
}

func InternalError(message string, err error) *AppError {
	return &AppError{
		Code:    CodeInternal,
		Message: message,
		Err:     err,
	}
}

// -----------------------------------

func (e *AppError) Error() string {

	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// --- Wrappers ---

func Wrap(err error, message string) error {

	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func WrapWithCode(code string, err error, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
