package usecase

import (
	"errors"
	"fmt"
	"time"
	"webscraper-v2/config"
	"webscraper-v2/domain/entity"
	"webscraper-v2/domain/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
	config   *config.Config
}

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthUseCase(userRepo repository.UserRepository, cfg *config.Config) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		config:   cfg,
	}
}

func (uc *AuthUseCase) Register(req *entity.RegisterRequest) (*entity.AuthResponse, error) {
	existsUsername, err := uc.userRepo.ExistsUsername(req.Username)

	if err != nil {
		return nil, fmt.Errorf("error checking username: %w", err)
	}
	if existsUsername {
		return nil, errors.New("username already exists")
	}
	existsEmail, err := uc.userRepo.ExistsEmail(req.Email)

	if err != nil {
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if existsEmail {
		return nil, errors.New("email already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		uc.config.Auth.BCryptCost,
	)

	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     uc.config.Auth.DefaultRole,
		Active:   true,
	}

	if err := uc.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	token, expiresAt, err := uc.generateToken(user)

	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}
	user.Password = ""

	return &entity.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) Login(req *entity.LoginRequest) (*entity.AuthResponse, error) {
	user, err := uc.userRepo.FindByUsername(req.Username)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	token, expiresAt, err := uc.generateToken(user)

	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}
	user.Password = ""

	return &entity.AuthResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *AuthUseCase) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(uc.config.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (uc *AuthUseCase) GetUserByID(id int64) (*entity.User, error) {
	user, err := uc.userRepo.FindByID(id)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user != nil {
		user.Password = "" // Limpiar password
	}
	return user, nil
}

func (uc *AuthUseCase) RefreshToken(tokenString string) (*entity.AuthResponse, error) {
	claims, err := uc.ValidateToken(tokenString)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	user, err := uc.userRepo.FindByID(claims.UserID)

	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	newToken, expiresAt, err := uc.generateToken(user)

	if err != nil {
		return nil, fmt.Errorf("error generating token: %w", err)
	}
	user.Password = ""

	return &entity.AuthResponse{
		Token:     newToken,
		User:      *user,
		ExpiresAt: expiresAt,
	}, nil
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
