package usecase

import (
	"context"
	"fmt"
	"time"
	"webscraper-v2/internal/domain/entity"
	"webscraper-v2/internal/domain/repository"
	"webscraper-v2/internal/infrastructure/config"
	"webscraper-v2/pkg/crypto"
	pkgerrors "webscraper-v2/pkg/errors"
	"webscraper-v2/pkg/validator"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCase struct {
	userRepo       repository.UserRepository
	tokenRepo      repository.TokenRepository
	config         *config.Config
	ctx            context.Context
	cancel         context.CancelFunc
	validator      *validator.Validator
	passwordHasher *crypto.PasswordHasher
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthUseCase(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, cfg *config.Config) *AuthUseCase {
	ctx, cancel := context.WithCancel(context.Background())

	uc := &AuthUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		config:         cfg,
		ctx:            ctx,
		cancel:         cancel,
		validator:      validator.NewValidator(),
		passwordHasher: crypto.NewPasswordHasher(cfg.Auth.BCryptCost),
	}
	go uc.cleanupExpiredTokens()
	return uc
}

func (uc *AuthUseCase) Shutdown() {
	uc.cancel()
}

func (uc *AuthUseCase) Register(req *entity.RegisterRequest) (*entity.AuthResponse, error) {

	if err := uc.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	existsUsername, err := uc.userRepo.ExistsUsername(req.Username)
	if err != nil {
		return nil, pkgerrors.DatabaseError("check username", err)
	}
	if existsUsername {
		return nil, pkgerrors.ConflictError("username already exists")
	}

	existsEmail, err := uc.userRepo.ExistsEmail(req.Email)
	if err != nil {
		return nil, pkgerrors.DatabaseError("check email", err)
	}
	if existsEmail {
		return nil, pkgerrors.ConflictError("email already exists")
	}

	hashedPassword, err := uc.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to hash password", err)
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     uc.config.Auth.DefaultRole,
		Active:   true,
	}
	if err := uc.userRepo.Create(user); err != nil {
		return nil, pkgerrors.DatabaseError("create user", err)
	}

	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to generate token", err)
	}

	user.Password = ""

	return &entity.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) Login(req *entity.LoginRequest) (*entity.AuthResponse, error) {

	if err := uc.validateLoginRequest(req); err != nil {
		return nil, err
	}

	user, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, pkgerrors.DatabaseError("find user", err)
	}
	if user == nil {
		return nil, pkgerrors.AuthenticationError("invalid credentials")
	}
	if !user.Active {
		return nil, pkgerrors.New(pkgerrors.CodeAuthentication, "user account is deactivated", pkgerrors.ErrUserInactive)
	}
	if err := uc.passwordHasher.Compare(user.Password, req.Password); err != nil {
		return nil, pkgerrors.AuthenticationError("invalid credentials")
	}

	token, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to generate token", err)
	}

	user.Password = ""

	return &entity.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*Claims, error) {
	isRevoked, err := uc.tokenRepo.IsTokenRevoked(tokenString)
	if err != nil {
		return nil, pkgerrors.DatabaseError("check token revocation", err)
	}
	if isRevoked {
		return nil, pkgerrors.New(pkgerrors.CodeAuthentication, "token has been revoked", pkgerrors.ErrTokenRevoked)
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, pkgerrors.New(pkgerrors.CodeAuthentication, "error parsing token", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, pkgerrors.New(pkgerrors.CodeAuthentication, "invalid token", pkgerrors.ErrInvalidToken)
}

func (uc *AuthUseCase) GetUserByID(id int64) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(id)
	if err != nil {
		return nil, pkgerrors.DatabaseError("find user", err)
	}
	if user != nil {
		user.Password = ""
	}
	return user, nil
}

func (uc *AuthUseCase) RefreshToken(tokenString string) (*entity.AuthResponse, error) {

	if err := uc.validator.ValidateRequired(tokenString, "token"); err != nil {
		return nil, pkgerrors.ValidationError(err.Error())
	}

	claims, err := uc.ValidateToken(tokenString)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "invalid token")
	}

	user, err := uc.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, pkgerrors.DatabaseError("find user", err)
	}
	if user == nil {
		return nil, pkgerrors.NotFoundError("user")
	}

	newToken, expiresAt, err := uc.generateToken(user)
	if err != nil {
		return nil, pkgerrors.InternalError("failed to generate token", err)
	}

	user.Password = ""

	return &entity.AuthResponse{
		Token:     newToken,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) RevokeToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return pkgerrors.Wrap(err, "failed to parse token")
	}

	if claims, ok := token.Claims.(*Claims); ok {
		if err := uc.tokenRepo.RevokeToken(tokenString, claims.ExpiresAt.Time); err != nil {
			return pkgerrors.DatabaseError("revoke token", err)
		}
		return nil
	}

	return pkgerrors.New(pkgerrors.CodeAuthentication, "invalid token claims", pkgerrors.ErrInvalidToken)
}

func (uc *AuthUseCase) generateToken(user *entity.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(uc.config.Auth.TokenDuration) * time.Hour)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "webscraper",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(uc.config.Auth.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (uc *AuthUseCase) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := uc.tokenRepo.CleanupExpiredTokens(); err != nil {
				fmt.Printf("Error cleaning up expired tokens: %v\n", err)
			}
		case <-uc.ctx.Done():
			return
		}
	}
}

func (uc *AuthUseCase) validateRegisterRequest(req *entity.RegisterRequest) error {

	if err := uc.validator.ValidateStruct(req, "registration request"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateUsername(req.Username); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateEmail(req.Email); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidatePassword(req.Password); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}

	return nil
}

func (uc *AuthUseCase) validateLoginRequest(req *entity.LoginRequest) error {

	if err := uc.validator.ValidateStruct(req, "login request"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateRequired(req.Username, "username"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}
	if err := uc.validator.ValidateRequired(req.Password, "password"); err != nil {
		return pkgerrors.ValidationError(err.Error())
	}

	return nil
}
