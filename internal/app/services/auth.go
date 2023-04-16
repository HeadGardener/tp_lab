package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/HeadHardener/tp_lab/internal/app/models"
	"github.com/HeadHardener/tp_lab/internal/app/repositories"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt      = "qetuoadgjlzcbmwryipsfhkxvn"
	tokenTTL  = 2 * time.Hour
	secretKey = "qazwsxedcrfvtgbyhnujm"
)

type AuthService struct {
	repos *repositories.Repository
}

func NewAuthService(repos *repositories.Repository) *AuthService {
	return &AuthService{repos: repos}
}

type tokenClaims struct {
	jwt.StandardClaims
	WorkerID int    `json:"worker_id"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

func (s *AuthService) GenerateToken(workerInput models.LogWorkerInput) (string, error) {
	worker := models.Worker{
		Name:         workerInput.Name,
		Surname:      workerInput.Surname,
		Phone:        workerInput.Phone,
		PasswordHash: getPasswordHash(workerInput.Password),
	}

	err := s.repos.WorkerInterface.GetWorker(&worker)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		worker.ID,
		worker.Role,
		fmt.Sprintf("%s %s", worker.Name, worker.Surname),
	})

	return token.SignedString([]byte(secretKey))
}

func (s *AuthService) ParseToken(accessToken string) (models.WorkerAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return models.WorkerAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return models.WorkerAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	return models.WorkerAttributes{
		ID:   claims.WorkerID,
		Role: claims.Role,
		Name: claims.Name,
	}, nil
}

func getPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
