package service

import (
	"context"
	"errors"
	"time"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService struct {
	cfg             entities.Config
	tokenRepository tokenRepository
}

func NewJwtService(cfg entities.Config, toktokenRepository tokenRepository) *JwtService {
	return &JwtService{
		cfg:             cfg,
		tokenRepository: toktokenRepository,
	}
}

func (j *JwtService) GetToken(ctx context.Context) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iat":      time.Now().Unix(),
		"token_id": uuid.NewString(),
	})

	return token.SignedString([]byte(j.cfg.JwtKey))
}

func (j *JwtService) Validate(ctx context.Context, token string) error {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Error token parsing")
		}
		return j.cfg.JwtKey, nil
	})
	if err != nil {
		return err
	}

	if jwtToken == nil {
		return errors.New("Invalid token")
	}

	issuedAt, err := jwtToken.Claims.GetIssuedAt()
	if err != nil {
		return err
	}

	if issuedAt.Add(time.Millisecond * time.Duration(j.cfg.JwtExpiryMs)).After(time.Now()) {
		return errors.New("Token expired")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Token error")
	}

	tokenId := claims["token_id"].(string)

	return j.tokenRepository.Save(ctx, tokenId)
}
